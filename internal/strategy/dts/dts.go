package dts

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

// DTS config props
const (
	_BUY_THRESHOLD         = "buyThreshold"
	_SELL_THRESHOLD        = "sellThreshold"
	_STOP_LOSS_THRESHOLD   = "stopLossThreshold"
	_MISS_PROFIT_THRESHOLD = "missProfitThreshold"
)

type OperationTypeDTS string

const (
	OP_BUY_DTS  OperationTypeDTS = "OP_BUY_DTS"
	OP_SELL_DTS OperationTypeDTS = "OP_SELL_DTS"
)

type AssetStatusDTS struct {
	Asset              string           `bson:"asset"`              // Asset being tracked
	Amount             decimal.Decimal  `bson:"amount"`             // Amount of that asset currently owned
	Usdt               decimal.Decimal  `bson:"usdt"`               // Usdt got by selling the asset
	LastOperationType  OperationTypeDTS `bson:"lastOperationType"`  // Last DTS operation type
	LastOperationPrice decimal.Decimal  `bson:"lastOperationPrice"` // Asset value at the time last op was executed
}

func (a AssetStatusDTS) IsEmpty() bool {
	return reflect.DeepEqual(a, AssetStatusDTS{})
}

type LocalAccountDTS struct {
	model.LocalAccountMetadata `bson:"metadata"`
	Ignored                    map[string]decimal.Decimal `bson:"ignored"` // Ignored assets
	Assets                     map[string]AssetStatusDTS  `bson:"assets"`  // Value allocation across assets
}

func (a LocalAccountDTS) IsEmpty() bool {
	return reflect.DeepEqual(a, LocalAccountDTS{})
}

type operation_init struct {
	asset       string
	amount      decimal.Decimal
	targetPrice decimal.Decimal
	cause       string
}

const (
	_NP_OP                = "NO_OP"
	_DTS_BUY_DESC         = "dts_buy"
	_DTS_SELL_DESC        = "dts_sell"
	_DTS_STOP_LOSS_DESC   = "dts_stop_loss"
	_DTS_MISS_PROFIT_DESC = "dts_miss_profit"
)

func (a LocalAccountDTS) Initialize(req model.LocalAccountInit) (model.ILocalAccount, errors.CtbError) {
	var ignored = make(map[string]decimal.Decimal)
	var assets = make(map[string]AssetStatusDTS)

	for _, rbalance := range req.RAccount.Balances {
		if decimal.Zero.Equals(rbalance.Amount) {
			continue
		}
		price, found := req.TradableAssetsPrice[rbalance.Asset]
		if !found {
			logrus.WithField("comp", "dts").
				Warnf(logger.XXX_IGNORED_ASSET, rbalance.Asset)
			ignored[rbalance.Asset] = rbalance.Amount
			continue
		}
		assetStatus := init_asset_status_dts(rbalance, price)
		assets[rbalance.Asset] = assetStatus
	}

	return LocalAccountDTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        req.ExeId,
			StrategyType: req.StrategyType,
			Timestamp:    time.Now().UnixMicro()},
		Ignored: ignored,
		Assets:  assets}, nil
}

func (a LocalAccountDTS) GetAssetAmounts() map[string]model.AssetAmount {
	assets := make(map[string]model.AssetAmount)
	assets["USDT"] = model.AssetAmount{"USDT", decimal.Zero}
	for asset, amount := range a.Ignored {
		assets[asset] = model.AssetAmount{asset, amount}
	}
	for asset, assetStatusDts := range a.Assets {
		assets[asset] = model.AssetAmount{asset, assetStatusDts.Amount}
		usdtStatus := assets["USDT"]
		usdtStatus.Amount = usdtStatus.Amount.Add(assetStatusDts.Usdt)
		assets["USDT"] = usdtStatus
	}
	return assets
}

func (a LocalAccountDTS) RegisterTrading(op model.Operation) (model.ILocalAccount, errors.CtbError) {
	// Check execution ids
	if op.ExeId != a.ExeId {
		logrus.WithField("comp", "dts").
			Panicf(logger.XXX_ERR_MISMATCHING_EXE_IDS, a.ExeId, op.ExeId)
	}

	// If the result status is failed, NOP
	if op.Status == model.FAILED {
		err := errors.Internal(logger.XXX_ERR_FAILED_OP, op.OpId)
		logrus.WithField("comp", "dts").Error(err.Error())
		return a, err
	}

	// DTS only handle operation back and forth USDT
	if op.Quote != "USDT" {
		err := errors.Internal(logger.DTS_ERR_BAD_QUOTE_CURRENCY, op.Quote)
		logrus.WithField("comp", "dts").Error(err.Error())
		return a, err
	}

	// Getting asset status
	assetStatus, found := a.Assets[op.Base]
	if !found {
		err := errors.Internal(logger.XXX_ERR_ASSET_NOT_FOUND, op.Base)
		logrus.WithField("comp", "dts").Error(err.Error())
		return a, err
	}

	// Updating asset status
	baseAmount := op.Results.BaseDiff
	quoteAmount := op.Results.QuoteDiff
	if op.Side == model.BUY {
		assetStatus.Amount = assetStatus.Amount.Add(baseAmount).Round(8)
		assetStatus.Usdt = assetStatus.Usdt.Sub(quoteAmount).Round(8)
		assetStatus.LastOperationType = OP_BUY_DTS
	} else if op.Side == model.SELL {
		assetStatus.Amount = assetStatus.Amount.Sub(baseAmount).Round(8)
		assetStatus.Usdt = assetStatus.Usdt.Add(quoteAmount).Round(8)
		assetStatus.LastOperationType = OP_SELL_DTS
	} else {
		err := errors.Internal(logger.XXX_ERR_UNKNWON_OP_TYPE, op.Type)
		logrus.WithField("comp", "dts").Error(err.Error())
		return a, err
	}
	if assetStatus.Amount.LessThan(decimal.Zero) {
		logrus.WithField("comp", "dts").
			Panicf(logger.XXX_ERR_NEGATIVE_BALANCE, assetStatus.Asset, assetStatus.Amount)
	}
	if assetStatus.Usdt.LessThan(decimal.Zero) {
		logrus.WithField("comp", "dts").
			Panicf(logger.XXX_ERR_NEGATIVE_BALANCE, "USDT", assetStatus.Usdt)
	}
	assetStatus.LastOperationPrice = op.Results.ActualPrice

	// Returning results
	a.Assets[op.Base] = assetStatus
	a.Timestamp = time.Now().UnixMicro()
	a.AccountId = uuid.NewString()
	return a, nil
}

func (a LocalAccountDTS) GetOperation(props map[string]string, mms model.MiniMarketStats, slimts model.SpotMarketLimits) (model.Operation, errors.CtbError) {
	asset := mms.Asset
	assetStatus, found := a.Assets[asset]
	if !found {
		err := errors.Internal(logger.XXX_ERR_ASSET_NOT_FOUND, asset)
		logrus.WithField("comp", "dts").Error(err.Error())
		return model.Operation{}, err
	}

	lastOpType := assetStatus.LastOperationType
	lastOpPrice := assetStatus.LastOperationPrice
	currentAmnt := assetStatus.Amount
	currentAmntUsdt := assetStatus.Usdt
	currentPrice := mms.LastPrice
	if currentPrice.Equals(decimal.Zero) {
		err := errors.Internal(logger.XXX_ERR_ZERO_EXP_PRICE, asset)
		logrus.WithField("comp", "dts").Errorf(err.Error())
		return model.Operation{}, err
	}

	dtsConfig, err := parse_config(props)
	if err != nil {
		return model.Operation{}, nil
	}

	sellPrice := utils.IncrementByPercentage(lastOpPrice, dtsConfig.SellThreshold)
	stopLossPrice := utils.IncrementByPercentage(lastOpPrice, utils.SignChangeDecimal(dtsConfig.StopLossThreshold))
	buyPrice := utils.IncrementByPercentage(lastOpPrice, utils.SignChangeDecimal(dtsConfig.BuyThreshold))
	missProfitPrice := utils.IncrementByPercentage(lastOpPrice, dtsConfig.MissProfitThreshold)

	var op model.Operation = model.Operation{}
	if lastOpType == OP_BUY_DTS && currentPrice.GreaterThanOrEqual(sellPrice) {
		// sell command
		operationInit := build_operation_init(asset, currentAmnt, currentPrice, _DTS_SELL_DESC)
		op = build_sell_op(a, operationInit)

	} else if lastOpType == OP_BUY_DTS && currentPrice.LessThanOrEqual(stopLossPrice) {
		// stop loss command
		operationInit := build_operation_init(asset, currentAmnt, currentPrice, _DTS_STOP_LOSS_DESC)
		op = build_sell_op(a, operationInit)

	} else if lastOpType == OP_SELL_DTS && currentPrice.LessThanOrEqual(buyPrice) {
		// buy command
		operationInit := build_operation_init(asset, currentAmntUsdt, currentPrice, _DTS_BUY_DESC)
		op = build_buy_op(a, operationInit)

	} else if lastOpType == OP_SELL_DTS && currentPrice.GreaterThanOrEqual(missProfitPrice) {
		// miss profit command
		operationInit := build_operation_init(asset, currentAmntUsdt, currentPrice, _DTS_MISS_PROFIT_DESC)
		op = build_buy_op(a, operationInit)
	} else {
		// no op
		logrus.WithField("comp", "dts").
			Debugf(logger.DTS_TRADE, _NP_OP, asset, lastOpType, lastOpPrice, currentPrice)
		return model.Operation{}, nil
	}

	err = check_spot_market_limits(op, slimts)
	if err != nil {
		return model.Operation{}, err
	}

	logrus.WithField("comp", "dts").
		Infof(logger.DTS_TRADE, op.Cause, asset, lastOpType, lastOpPrice, currentPrice)
	return op, nil
}

func build_operation_init(asset string, amount, price decimal.Decimal, cause string) operation_init {
	return operation_init{
		asset:       asset,
		amount:      amount,
		targetPrice: price,
		cause:       cause}
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

func build_buy_op(laccount LocalAccountDTS, operationInit operation_init) model.Operation {
	return model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      laccount.ExeId,
		Type:       model.AUTO,
		Base:       operationInit.asset,
		Quote:      "USDT",
		Side:       model.BUY,
		Amount:     operationInit.amount,
		AmountSide: model.QUOTE_AMOUNT,
		Price:      operationInit.targetPrice,
		Cause:      operationInit.cause,
		Status:     model.PENDING}
}

func build_sell_op(laccount LocalAccountDTS, operationInit operation_init) model.Operation {
	return model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      laccount.ExeId,
		Type:       model.AUTO,
		Base:       operationInit.asset,
		Quote:      "USDT",
		Side:       model.SELL,
		Amount:     operationInit.amount,
		AmountSide: model.BASE_AMOUNT,
		Price:      operationInit.targetPrice,
		Cause:      operationInit.cause,
		Status:     model.PENDING}
}

func init_asset_status_dts(rbalance model.RemoteBalance, price model.AssetPrice) AssetStatusDTS {
	return AssetStatusDTS{
		Asset:              rbalance.Asset,
		Amount:             rbalance.Amount,
		Usdt:               decimal.Zero,
		LastOperationType:  OP_BUY_DTS,
		LastOperationPrice: price.Price}
}
