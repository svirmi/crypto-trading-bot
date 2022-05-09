package handler

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tevino/abool/v2"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/operations"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

type trading_context struct {
	laccount  model.ILocalAccount
	execution model.Execution
	mms       chan []model.MiniMarketStats
	exchange  model.IExchange
}

var tcontext trading_context

func Initialize(laccount model.ILocalAccount, execution model.Execution,
	mmsChannel chan []model.MiniMarketStats, exchange model.IExchange) {

	tcontext.laccount = laccount
	tcontext.execution = execution
	tcontext.mms = mmsChannel
	tcontext.exchange = exchange
}

func InvalidateTradingContext() {
	tcontext.execution = model.Execution{}
	tcontext.laccount = nil
}

func HandleMiniMarketsStats() {
	trading_context_init()

	go handle_mini_markets_stats()
}

var handle_mini_markets_stats = func() {
	sentinel := abool.New()

	for miniMarketsStats := range tcontext.mms {
		// If the execution is not ACTIVE, no action should be applied
		if tcontext.execution.Status != model.EXE_ACTIVE {
			continue
		}

		for _, miniMarketStats := range miniMarketsStats {
			// Check that trading is enabled for given asset
			symbol := utils.GetSymbolFromAsset(miniMarketStats.Asset)
			if !can_spot_trade(symbol) {
				logrus.Warnf(logger.HANDL_TRADING_DISABLED, symbol)
				continue
			}

			// Getting spot market limits
			slimits, err := get_spot_market_limits(symbol)
			if err != nil {
				logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, miniMarketStats.Asset, err.Error())
				continue
			}

			// Trading ongoing, skip market stats update
			ok := sentinel.SetToIf(false, true)
			if !ok {
				skip_mini_market_stats(miniMarketsStats)
				continue
			}

			// Getting operation
			op, err := get_operation(miniMarketStats, slimits)
			if err != nil {
				logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, miniMarketStats.Asset, err.Error())
				sentinel.UnSet()
				continue
			}
			if op.IsEmpty() {
				sentinel.UnSet()
				continue
			}

			// Set sentinel, handle operation and defer sentinel reset
			go func(op model.Operation) {
				defer sentinel.UnSet()
				handle_operation(op)
			}(op)
		}
	}
}

var can_spot_trade = func(symbol string) bool {
	return tcontext.exchange.CanSpotTrade(symbol)
}

var get_spot_market_limits = func(symbol string) (model.SpotMarketLimits, error) {
	return tcontext.exchange.GetSpotMarketLimits(symbol)
}

var get_operation = func(miniMarketStats model.MiniMarketStats, slimits model.SpotMarketLimits) (model.Operation, error) {
	return tcontext.laccount.GetOperation(miniMarketStats, slimits)
}

var skip_mini_market_stats = func([]model.MiniMarketStats) {
	logrus.Info(logger.HANDL_SKIP_MMS_UPDATE)
}

var handle_operation = func(op model.Operation) {
	// Price equals to zero
	if op.Amount.Equals(decimal.Zero) {
		logrus.Errorf(logger.HANDL_ERR_ZERO_REQUESTED_AMOUNT)
		return
	}
	if op.Price.Equals(decimal.Zero) {
		logrus.Errorf(logger.HANDL_ERR_ZERO_EXP_PRICE)
		return
	}

	// Getting remote account before operation
	raccount1, err := tcontext.exchange.GetAccout()
	if err != nil {
		logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
		return
	}

	// Sending market order
	op, err = tcontext.exchange.SendSpotMarketOrder(op)
	if err != nil {
		logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
		return
	}

	// Getting remote account after operation
	raccount2, err := tcontext.exchange.GetAccout()
	if err != nil {
		logrus.Panicf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
	}

	// Computing operation results
	op.FromId = tcontext.laccount.GetAccountId()
	op, err = compute_op_results(raccount1, raccount2, op)
	if err != nil {
		logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
		return
	}

	// Updating local account
	tcontext.laccount, err = tcontext.laccount.RegisterTrading(op)
	if err != nil {
		logrus.Panicf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
	}

	// Inserting operation and updating laccount in DB
	op.ToId = tcontext.laccount.GetAccountId()
	err = operations.Create(op)
	if err != nil {
		logrus.Panicf(err.Error())
	}

	err = laccount.Create(tcontext.laccount)
	if err != nil {
		logrus.Panicf(err.Error())
	}
}

func compute_op_results(old, new model.RemoteAccount, op model.Operation) (model.Operation, error) {
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
		err := fmt.Errorf(logger.HANDL_ERR_ZERO_BASE_QUOTE_DIFF)
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

func trading_context_init() {
	if tcontext.execution.IsEmpty() {
		execution, err := executions.GetCurrentlyActive()
		if err != nil {
			logrus.Panicf(err.Error())
		}
		tcontext.execution = execution
	}

	if tcontext.laccount == nil {
		laccount, err := laccount.GetLatestByExeId(tcontext.execution.ExeId)
		if err != nil {
			logrus.Panicf(err.Error())
		}
		tcontext.laccount = laccount
	}
}
