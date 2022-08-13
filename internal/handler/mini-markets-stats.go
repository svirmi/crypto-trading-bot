package handler

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool/v2"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/exchange"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/operations"
	"github.com/valerioferretti92/crypto-trading-bot/internal/prices"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

var (
	mmsChannel      chan []model.MiniMarketStats
	callbackChannel chan model.MiniMarketStatsAck
)

func Initialize(mmsCh chan []model.MiniMarketStats, callbackCh chan model.MiniMarketStatsAck) {
	mmsChannel = mmsCh
	callbackChannel = callbackCh
}

func HandleMiniMarketsStats() {
	go handle_mini_markets_stats()
}

var handle_mini_markets_stats = func() {
	sentinel := abool.New()

	for mmss := range mmsChannel {
		// Store prices in db
		store_prices_deferred(mmss)

		// Getting execution
		exe, err := get_latest_exe()
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_SKIP_MMSS_UPDATE, err.Error())
			ack_mmss(len(mmss))
			continue
		}
		// If the execution is not ACTIVE, no action should be applied
		if exe.IsEmpty() || exe.Status != model.EXE_ACTIVE {
			logrus.Debug(logger.HANDL_NO_ACTIVE_EXECUTION)
			ack_mmss(len(mmss))
			continue
		}

		// Getting local account
		lacc, err := get_latest_lacc(exe.ExeId)
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_SKIP_MMSS_UPDATE, err.Error())
			ack_mmss(len(mmss))
			continue
		}
		if lacc == nil {
			err := errors.NotFound(logger.HANDL_ERR_LACC_NOT_FOUND, exe.ExeId)
			logrus.Warnf(logger.HANDL_ERR_SKIP_MMSS_UPDATE, err.Error())
			ack_mmss(len(mmss))
			continue
		}

		for _, mms := range mmss {
			// Check if asset is in wallet
			assetStatuses := get_asset_amounts(lacc)
			if _, found := assetStatuses[mms.Asset]; !found {
				ack_mmss(1)
				continue
			}
			logrus.Tracef(logger.HANDL_MMS_HANDLING, mms.Asset)

			// Check that trading is enabled for given asset
			symbol, err := utils.GetSymbolFromAsset(mms.Asset)
			if err != nil {
				logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, mms.Asset, err.Error())
				ack_mmss(1)
				continue
			}
			if !can_spot_trade(symbol) {
				logrus.Warnf(logger.HANDL_TRADING_DISABLED, symbol)
				ack_mmss(1)
				continue
			}

			// Getting spot market limits
			slimits, err := get_spot_market_limits(symbol)
			if err != nil {
				logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, mms.Asset, err.Error())
				ack_mmss(1)
				continue
			}

			// Trading ongoing, skip market stats update
			ok := sentinel.SetToIf(false, true)
			if !ok {
				skip_mini_market_stats(mmss)
				ack_mmss(1)
				continue
			}

			// Getting operation
			op, err := get_operation(exe, lacc, mms, slimits)
			if err != nil {
				logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, mms.Asset, err.Error())
				ack_mmss(1)
				sentinel.UnSet()
				continue
			}
			if op.IsEmpty() {
				ack_mmss(1)
				sentinel.UnSet()
				continue
			}

			// Set sentinel, handle operation and defer sentinel reset
			go func(op model.Operation) {
				defer func() {
					ack_mmss(1)
					sentinel.UnSet()
				}()
				lacc = handle_operation(lacc, op)
			}(op)
		}
	}
}

var ack_mmss = func(size int) {
	if callbackChannel == nil {
		return
	}

	select {
	case callbackChannel <- model.MiniMarketStatsAck{Count: size}:
	default:
		logrus.Errorf(logger.HANDL_ERR_FAILED_TO_ACK_MMSS,
			len(callbackChannel), cap(callbackChannel))
	}
}

var get_latest_exe = func() (model.Execution, errors.CtbError) {
	return executions.GetLatest()
}

var get_latest_lacc = func(exeId string) (model.ILocalAccount, errors.CtbError) {
	return laccount.GetLatestByExeId(exeId)
}

var store_prices_deferred = func(mmss []model.MiniMarketStats) {
	timestamp := time.Now().UnixMicro()
	symbolPrices := make([]model.SymbolPrice, 0, len(mmss))
	for _, mms := range mmss {
		symbol, err := utils.GetSymbolFromAsset(mms.Asset)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}

		symbolPrices = append(symbolPrices, model.SymbolPrice{
			Symbol:    symbol,
			Price:     mms.LastPrice,
			Timestamp: timestamp})
	}
	prices.InsertManyDeferred(symbolPrices)
}

var can_spot_trade = func(symbol string) bool {
	return exchange.CanSpotTrade(symbol)
}

var get_spot_market_limits = func(symbol string) (model.SpotMarketLimits, errors.CtbError) {
	return exchange.GetSpotMarketLimits(symbol)
}

var get_operation = func(exe model.Execution, lacc model.ILocalAccount, mms model.MiniMarketStats, slimits model.SpotMarketLimits) (model.Operation, errors.CtbError) {
	return lacc.GetOperation(exe.Props, mms, slimits)
}

var skip_mini_market_stats = func([]model.MiniMarketStats) {
	logrus.Info(logger.HANDL_SKIP_MMS_UPDATE)
}

var get_asset_amounts = func(lacc model.ILocalAccount) map[string]model.AssetAmount {
	return lacc.GetAssetAmounts()
}

var handle_operation = func(lacc model.ILocalAccount, op model.Operation) model.ILocalAccount {
	// Price equals to zero
	if op.Amount.Equals(decimal.Zero) {
		logrus.Errorf(logger.HANDL_ERR_ZERO_REQUESTED_AMOUNT)
		return lacc
	}
	if op.Price.Equals(decimal.Zero) {
		logrus.Errorf(logger.HANDL_ERR_ZERO_EXP_PRICE)
		return lacc
	}

	// Getting remote account before operation
	raccount1, err := exchange.GetAccount()
	if err != nil {
		logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
		return lacc
	}

	// Sending market order
	op, err = exchange.SendSpotMarketOrder(op)
	if err != nil {
		logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
		return lacc
	}

	// Getting remote account after operation
	raccount2, err := exchange.GetAccount()
	if err != nil {
		logrus.Panicf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
	}

	// Computing operation results
	op.FromId = lacc.GetAccountId()
	op, err = compute_op_results(raccount1, raccount2, op)
	if err != nil {
		logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
		return lacc
	}

	// Updating local account
	lacc, err = lacc.RegisterTrading(op)
	if err != nil {
		logrus.Panicf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
	}

	// Inserting operation and updating laccount in DB
	op.ToId = lacc.GetAccountId()
	err = operations.Create(op)
	if err != nil {
		logrus.Panicf(err.Error())
	}

	lacc, err = laccount.Update(lacc)
	if err != nil {
		logrus.Panicf(err.Error())
	}
	return lacc
}

func compute_op_results(old, new model.RemoteAccount, op model.Operation) (model.Operation, errors.CtbError) {
	var oldBaseBalance, newBaseBalance decimal.Decimal
	var oldQuoteBalance, newQuoteBalance decimal.Decimal

	// Getting old abase and quote balances
	for _, balance := range old.Balances {
		if balance.Asset == op.Base {
			oldBaseBalance = balance.Amount
		}
		if balance.Asset == op.Quote {
			oldQuoteBalance = balance.Amount
		}
	}

	// Getting new base and quote balances
	for _, balance := range new.Balances {
		if balance.Asset == op.Base {
			newBaseBalance = balance.Amount
		}
		if balance.Asset == op.Quote {
			newQuoteBalance = balance.Amount
		}
	}

	// Checking base diff and quote diff
	baseDiff := (newBaseBalance.Sub(oldBaseBalance)).Abs().Round(8)
	quoteDiff := (newQuoteBalance.Sub(oldQuoteBalance)).Abs().Round(8)
	if baseDiff.Equals(decimal.Zero) && quoteDiff.Equals(decimal.Zero) {
		err := errors.Internal(logger.HANDL_ERR_ZERO_BASE_QUOTE_DIFF)
		logrus.Error(err.Error())
		op.Status = model.FAILED
		return op, err
	} else if baseDiff.Equals(decimal.Zero) {
		logrus.Warnf(logger.HANDL_ZERO_BASE_DIFF, op.OpId)
	} else if quoteDiff.Equals(decimal.Zero) {
		logrus.Warnf(logger.HANDL_ZERO_QUOTE_DIFF, op.OpId)
	}

	// Computing status
	if baseDiff.Equals(decimal.Zero) || quoteDiff.Equals(decimal.Zero) {
		op.Status = model.PARTIALLY_FILLED
	} else if op.AmountSide == model.BASE_AMOUNT && baseDiff.Equals(op.Amount) {
		op.Status = model.FILLED
	} else if op.AmountSide == model.QUOTE_AMOUNT && quoteDiff.Equals(op.Amount) {
		op.Status = model.FILLED
	} else {
		op.Status = model.PARTIALLY_FILLED
	}

	// Computing actual price and spread
	var actualPrice decimal.Decimal
	var spread decimal.Decimal
	if baseDiff.Equals(decimal.Zero) && op.Side == model.BUY {
		actualPrice = utils.MaxDecimal()
		spread = utils.MaxDecimal()
	} else if baseDiff.Equals(decimal.Zero) && op.Side == model.SELL {
		actualPrice = decimal.Zero
		spread = utils.DecimalFromString("-100")
	} else if quoteDiff.Equals(decimal.Zero) && op.Side == model.BUY {
		actualPrice = decimal.Zero
		spread = utils.DecimalFromString("-100")
	} else if quoteDiff.Equals(decimal.Zero) && op.Side == model.SELL {
		actualPrice = utils.MaxDecimal()
		spread = utils.MaxDecimal()
	} else {
		actualPrice = quoteDiff.Div(baseDiff).Round(8)
		spread = ((actualPrice.Sub(op.Price)).
			Div(op.Price)).
			Mul(decimal.NewFromInt(100)).
			Round(8)
	}

	// Setting results
	results := model.OpResults{
		ActualPrice: actualPrice,
		BaseDiff:    baseDiff,
		QuoteDiff:   quoteDiff,
		Spread:      spread}
	op.Results = results

	logrus.Infof(logger.HANDL_OPERATION_RESULTS,
		op.Results.BaseDiff, op.Results.QuoteDiff, op.Results.ActualPrice, op.Results.Spread, op.Status)
	return op, nil
}
