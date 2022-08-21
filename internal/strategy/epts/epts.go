package epts

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

const (
	_NP_OP          = "NO_OP"
	_NP_OP_FUNDS    = "NO_OP(funds)"
	_EPTS_BUY_DESC  = "epts_buy"
	_EPTS_SELL_DESC = "epts_sell"
)

const (
	_BUY_PERCENTAGE              = "buyPercentage"
	_SELL_PERCENTAGE             = "sellPercentage"
	_INIT_BUY_AMOUNT_PERCENTAGE  = "initBuyAmountPercentage"
	_INIT_SELL_AMOUNT_PERCENTAGE = "initSellAmountPercentage"
	_EXPONENTIAL_BASE            = "exponentialBase"
)

type OperationTypeEPTS string

const (
	_OP_INIT_EPTS OperationTypeEPTS = "OP_INIT_EPTS"
	_OP_BUY_EPTS  OperationTypeEPTS = "OP_BUY_EPTS"
	_OP_SELL_EPTS OperationTypeEPTS = "OP_SELL_EPTS"
)

type AssetStatusEPTS struct {
	Asset               string            `bson:"asset"`
	Amount              decimal.Decimal   `bson:"amount"`
	LastOperationPrice  decimal.Decimal   `bson:"lastOperationPrice"`
	LastOperationAmount decimal.Decimal   `bson:"lastOperationAmount"`
	LastOperationType   OperationTypeEPTS `bson:"lastOperationType"`
}

func (a AssetStatusEPTS) IsEmpty() bool {
	return reflect.DeepEqual(a, AssetStatusEPTS{})
}

type LocalAccountEPTS struct {
	model.LocalAccountMetadata `bson:"metadata"`
	Ignored                    map[string]decimal.Decimal `bson:"ignored"`
	Assets                     map[string]AssetStatusEPTS `bson:"assets"`
	Usdt                       decimal.Decimal            `bson:"usdt"`
}

func (a LocalAccountEPTS) IsEmpty() bool {
	return reflect.DeepEqual(a, LocalAccountEPTS{})
}

func (a LocalAccountEPTS) Initialize(req model.LocalAccountInit) (model.ILocalAccount, errors.CtbError) {
	var ignored = make(map[string]decimal.Decimal)
	var assets = make(map[string]AssetStatusEPTS)
	var usdt = decimal.Zero

	for _, rbalance := range req.RAccount.Balances {
		if rbalance.Asset == "USDT" {
			usdt = rbalance.Amount
			continue
		}

		if decimal.Zero.Equals(rbalance.Amount) {
			continue
		}
		price, found := req.TradableAssetsPrice[rbalance.Asset]
		if !found {
			logrus.WithField("comp", "epts").
				Warnf(logger.XXX_IGNORED_ASSET, rbalance.Asset)
			ignored[rbalance.Asset] = rbalance.Amount
			continue
		}
		assetStatus := AssetStatusEPTS{
			Asset:               rbalance.Asset,
			Amount:              rbalance.Amount,
			LastOperationType:   _OP_INIT_EPTS,
			LastOperationPrice:  price.Price,
			LastOperationAmount: rbalance.Amount}
		assets[rbalance.Asset] = assetStatus
	}

	return LocalAccountEPTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        req.ExeId,
			StrategyType: req.StrategyType,
			Timestamp:    time.Now().UnixMicro()},
		Ignored: ignored,
		Assets:  assets,
		Usdt:    usdt}, nil
}

func (a LocalAccountEPTS) GetAssetAmounts() map[string]model.AssetAmount {
	assets := make(map[string]model.AssetAmount)
	for asset, amount := range a.Ignored {
		assets[asset] = model.AssetAmount{asset, amount}
	}
	for asset, assetStatusEPTS := range a.Assets {
		assets[asset] = model.AssetAmount{asset, assetStatusEPTS.Amount}
	}
	assets["USDT"] = model.AssetAmount{"USDT", a.Usdt}
	return assets
}

func (a LocalAccountEPTS) RegisterTrading(op model.Operation) (model.ILocalAccount, errors.CtbError) {
	// Check execution ids
	if op.ExeId != a.ExeId {
		logrus.WithField("comp", "epts").
			Panicf(logger.XXX_ERR_MISMATCHING_EXE_IDS, a.ExeId, op.ExeId)
	}

	// If the result status is failed, NOP
	if op.Status == model.FAILED {
		err := errors.Internal(logger.XXX_ERR_FAILED_OP, op.OpId)
		logrus.WithField("comp", "epts").Error(err.Error())
		return a, err
	}

	// DTS only handle operation back and forth USDT
	if op.Quote != "USDT" {
		err := errors.Internal(logger.PTS_ERR_BAD_QUOTE_CURRENCY, op.Quote)
		logrus.WithField("comp", "epts").Error(err.Error())
		return a, err
	}

	// Getting asset status
	assetStatus, found := a.Assets[op.Base]
	if !found {
		err := errors.Internal(logger.XXX_ERR_ASSET_NOT_FOUND, op.Base)
		logrus.WithField("comp", "epts").Error(err.Error())
		return a, err
	}

	// Updating asset status
	currentAmntUsdt := a.Usdt
	baseAmount := op.Results.BaseDiff
	quoteAmount := op.Results.QuoteDiff
	if op.Side == model.BUY {
		assetStatus.Amount = assetStatus.Amount.Add(baseAmount).Round(8)
		assetStatus.LastOperationType = _OP_BUY_EPTS
		assetStatus.LastOperationAmount = quoteAmount
		currentAmntUsdt = currentAmntUsdt.Sub(quoteAmount).Round(8)
	} else if op.Side == model.SELL {
		assetStatus.Amount = assetStatus.Amount.Sub(baseAmount).Round(8)
		assetStatus.LastOperationType = _OP_SELL_EPTS
		assetStatus.LastOperationAmount = baseAmount
		currentAmntUsdt = currentAmntUsdt.Add(quoteAmount).Round(8)
	} else {
		err := errors.Internal(logger.XXX_ERR_UNKNWON_OP_TYPE, op.Type)
		logrus.WithField("comp", "epts").Error(err.Error())
		return a, err
	}
	if assetStatus.Amount.LessThan(decimal.Zero) {
		logrus.WithField("comp", "epts").
			Panicf(logger.XXX_ERR_NEGATIVE_BALANCE, assetStatus.Asset, assetStatus.Amount)
	}
	if currentAmntUsdt.LessThan(decimal.Zero) {
		logrus.WithField("comp", "epts").
			Panicf(logger.XXX_ERR_NEGATIVE_BALANCE, "USDT", currentAmntUsdt)
	}
	assetStatus.LastOperationPrice = op.Results.ActualPrice

	// Returning results
	a.Assets[op.Base] = assetStatus
	a.Usdt = currentAmntUsdt
	a.Timestamp = time.Now().UnixMicro()
	a.AccountId = uuid.NewString()
	return a, nil
}

func (a LocalAccountEPTS) GetOperation(props map[string]string, mms model.MiniMarketStats, slimts model.SpotMarketLimits) (model.Operation, errors.CtbError) {
	asset := mms.Asset
	assetStatus, found := a.Assets[asset]
	if !found {
		err := errors.Internal(logger.XXX_ERR_ASSET_NOT_FOUND, asset)
		logrus.WithField("comp", "epts").Error(err.Error())
		return model.Operation{}, err
	}

	lastOpAmnt := assetStatus.LastOperationAmount
	lastOpType := assetStatus.LastOperationType
	lastOpPrice := assetStatus.LastOperationPrice
	currentAmnt := assetStatus.Amount
	currentAmntUsdt := a.Usdt
	currentPrice := mms.LastPrice
	if currentPrice.Equals(decimal.Zero) {
		err := errors.Internal(logger.XXX_ERR_ZERO_EXP_PRICE, asset)
		logrus.WithField("comp", "epts").Errorf(err.Error())
		return model.Operation{}, err
	}

	config, err := parse_config(props)
	if err != nil {
		return model.Operation{}, err
	}

	sellPrice := utils.IncrementByPercentage(lastOpPrice, config.SellPercentage)
	buyPrice := utils.IncrementByPercentage(lastOpPrice, utils.SignChangeDecimal(config.BuyPercentage))
	sellAmnt := utils.PercentageOf(currentAmnt, config.InitSellAmountPercentage)
	buyAmnt := utils.PercentageOf(currentAmntUsdt, config.InitBuyAmountPercentage)
	expAmnt := lastOpAmnt.Mul(config.ExponentialBase).Round(8)

	var op model.Operation
	if currentPrice.GreaterThanOrEqual(sellPrice) && lastOpType == _OP_INIT_EPTS {
		op = build_sell_op(a.ExeId, asset, _EPTS_SELL_DESC, sellAmnt, currentPrice)
	} else if currentPrice.GreaterThanOrEqual(sellPrice) && lastOpType == _OP_BUY_EPTS {
		op = build_sell_op(a.ExeId, asset, _EPTS_SELL_DESC, sellAmnt, currentPrice)
	} else if currentPrice.GreaterThanOrEqual(sellPrice) && lastOpType == _OP_SELL_EPTS {
		op = build_sell_op(a.ExeId, asset, _EPTS_SELL_DESC, expAmnt, currentPrice)
	} else if currentPrice.LessThanOrEqual(buyPrice) && lastOpType == _OP_INIT_EPTS {
		op = build_buy_op(a.ExeId, asset, _EPTS_BUY_DESC, buyAmnt, currentPrice)
	} else if currentPrice.LessThanOrEqual(buyPrice) && lastOpType == _OP_SELL_EPTS {
		op = build_buy_op(a.ExeId, asset, _EPTS_BUY_DESC, buyAmnt, currentPrice)
	} else if currentPrice.LessThanOrEqual(buyPrice) && lastOpType == _OP_BUY_EPTS {
		op = build_buy_op(a.ExeId, asset, _EPTS_BUY_DESC, expAmnt, currentPrice)
	} else {
		logrus.WithField("comp", "epts").
			Debugf(logger.EPTS_TRADE, _NP_OP, asset, lastOpType, lastOpPrice, currentPrice)
		return model.Operation{}, nil
	}

	if op.Side == model.SELL && op.Amount.GreaterThan(currentAmnt) {
		op.Amount = currentAmnt
	}
	if op.Side == model.BUY && op.Amount.GreaterThan(currentAmntUsdt) {
		op.Amount = currentAmntUsdt
	}
	if op.Amount.Equals(decimal.Zero) {
		logrus.WithField("comp", "epts").
			Debugf(logger.EPTS_TRADE, _NP_OP_FUNDS, asset, lastOpType, lastOpPrice, currentPrice)
		return model.Operation{}, nil
	}

	err = check_spot_market_limits(op, slimts)
	if err != nil {
		return model.Operation{}, err
	}

	logrus.WithField("comp", "epts").
		Infof(logger.EPTS_TRADE, op.Cause, asset, lastOpType, lastOpPrice, currentPrice)
	return op, nil
}

func check_spot_market_limits(op model.Operation, slimits model.SpotMarketLimits) errors.CtbError {
	if op.AmountSide == model.QUOTE_AMOUNT && op.Amount.LessThan(slimits.MinQuote) {
		err := errors.Internal(logger.XXX_BELOW_QUOTE_LIMIT,
			op.Base+op.Quote, op.Side, op.Amount, op.AmountSide, slimits.MinQuote)
		logrus.WithField("comp", "dts").Error(err.Error())
		return err
	}
	if op.AmountSide == model.BASE_AMOUNT && op.Amount.LessThan(slimits.MinBase) {
		err := errors.Internal(logger.XXX_BELOW_BASE_LIMIT,
			op.Base+op.Quote, op.Side, op.Amount, op.AmountSide, slimits.MinBase)
		logrus.WithField("comp", "dts").Error(err.Error())
		return err
	}
	// No checks on MaxBase as big orders should be broken down into
	// smaller orders by the exchange package
	return nil
}

func build_buy_op(exeId, asset, cause string, amnt, targetPrice decimal.Decimal) model.Operation {
	return model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      exeId,
		Type:       model.AUTO,
		Base:       asset,
		Quote:      "USDT",
		Side:       model.BUY,
		Amount:     amnt,
		AmountSide: model.QUOTE_AMOUNT,
		Price:      targetPrice,
		Cause:      cause,
		Status:     model.PENDING}
}

func build_sell_op(exeId, asset, cause string, amnt, targetPrice decimal.Decimal) model.Operation {
	return model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      exeId,
		Type:       model.AUTO,
		Base:       asset,
		Quote:      "USDT",
		Side:       model.SELL,
		Amount:     amnt,
		AmountSide: model.BASE_AMOUNT,
		Price:      targetPrice,
		Cause:      cause,
		Status:     model.PENDING}
}
