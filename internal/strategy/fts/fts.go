package fts

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

type OperationTypeFts string

const (
	OP_BUY_FTS  OperationTypeFts = "OP_BUY_FTS"
	OP_SELL_FTS OperationTypeFts = "OP_SELL_FTS"
)

type AssetStatusFTS struct {
	Asset              string           `bson:"asset"`              // Asset being tracked
	Amount             decimal.Decimal  `bson:"amount"`             // Amount of that asset currently owned
	Usdt               decimal.Decimal  `bson:"usdt"`               // Usdt got by selling the asset
	LastOperationType  OperationTypeFts `bson:"lastOperationType"`  // Last FTS operation type
	LastOperationPrice decimal.Decimal  `bson:"lastOperationPrice"` // Asset value at the time last op was executed
}

func (a AssetStatusFTS) IsEmpty() bool {
	return reflect.DeepEqual(a, AssetStatusFTS{})
}

type LocalAccountFTS struct {
	model.LocalAccountMetadata `bson:"metadata"`
	Ignored                    map[string]decimal.Decimal `bson:"ignored"` // Usdt not to be invested
	Assets                     map[string]AssetStatusFTS  `bson:"assets"`  // Value allocation across assets
}

func (a LocalAccountFTS) IsEmpty() bool {
	return reflect.DeepEqual(a, LocalAccountFTS{})
}

type operation_init struct {
	asset       string
	amount      decimal.Decimal
	targetPrice decimal.Decimal
}

func (a LocalAccountFTS) Initialize(creationRequest model.LocalAccountInit) (model.ILocalAccount, error) {
	var ignored = make(map[string]decimal.Decimal)
	var assets = make(map[string]AssetStatusFTS)

	for _, rbalance := range creationRequest.RAccount.Balances {
		price, found := creationRequest.TradableAssetsPrice[rbalance.Asset]
		if !found {
			logrus.WithField("comp", "fts").
				Warnf(logger.FTS_IGNORED_ASSET, rbalance.Asset)
			ignored[rbalance.Asset] = rbalance.Amount
			continue
		}
		symbol := utils.GetSymbolFromAsset(rbalance.Asset)
		if !binance.CanSpotTrade(symbol) {
			ignored[rbalance.Asset] = rbalance.Amount
			continue
		}
		assetStatus := init_asset_status_fts(rbalance, price)
		assets[rbalance.Asset] = assetStatus
	}

	return LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:        uuid.NewString(),
			ExeId:            creationRequest.ExeId,
			StrategyType:     model.FIXED_THRESHOLD_STRATEGY,
			SpotMarketLimits: creationRequest.SpotMarketLimits,
			Timestamp:        time.Now().UnixMicro()},
		Ignored: ignored,
		Assets:  assets}, nil
}

func (a LocalAccountFTS) RegisterTrading(op model.Operation) (model.ILocalAccount, error) {
	// Check execution ids
	if op.ExeId != a.ExeId {
		msg := fmt.Sprintf(logger.FTS_ERR_MISMATCHING_EXE_IDS, a.ExeId, op.ExeId)
		logrus.Error(msg)
		return a, model.NewCtbError(msg, false)
	}

	// If the result status is failed, NOP
	if op.Status == model.FAILED {
		msg := fmt.Sprintf(logger.FTS_ERR_FAILED_OP, op.OpId)
		logrus.Error(msg)
		return a, model.NewCtbError(msg, true)
	}

	// FTS only handle operation back and forth USDT
	if op.Quote != "USDT" {
		msg := fmt.Sprintf(logger.FTS_ERR_BAD_QUOTE_CURRENCY, op.Quote)
		logrus.Error(msg)
		return a, model.NewCtbError(msg, true)
	}

	// Getting asset status
	assetStatus, found := a.Assets[op.Base]
	if !found {
		msg := fmt.Sprintf(logger.FTS_ERR_ASSET_NOT_FOUND, op.Base)
		logrus.Error(msg)
		return a, model.NewCtbError(msg, true)
	}

	// Updating asset status
	baseAmount := op.Results.BaseDiff
	quoteAmount := op.Results.QuoteDiff
	if op.Side == model.BUY {
		assetStatus.Amount = assetStatus.Amount.Add(baseAmount).Round(8)
		assetStatus.Usdt = assetStatus.Usdt.Sub(quoteAmount).Round(8)
		assetStatus.LastOperationType = OP_BUY_FTS
	} else if op.Side == model.SELL {
		assetStatus.Amount = assetStatus.Amount.Sub(baseAmount).Round(8)
		assetStatus.Usdt = assetStatus.Usdt.Add(quoteAmount).Round(8)
		assetStatus.LastOperationType = OP_SELL_FTS
	} else {
		msg := fmt.Sprintf(logger.FTS_ERR_UNKNWON_OP_TYPE, op.Type)
		logrus.Error(msg)
		return a, model.NewCtbError(msg, true)
	}
	if assetStatus.Amount.LessThan(decimal.Zero) {
		msg := fmt.Sprintf(logger.FTS_ERR_NEGATIVE_BALANCE, assetStatus.Asset, assetStatus.Amount)
		logrus.Error(msg)
		return a, model.NewCtbError(msg, false)
	}
	if assetStatus.Usdt.LessThan(decimal.Zero) {
		msg := fmt.Sprintf(logger.FTS_ERR_NEGATIVE_BALANCE, "USDT", assetStatus.Usdt)
		logrus.Error(msg)
		return a, model.NewCtbError(msg, false)
	}
	assetStatus.LastOperationPrice = op.Results.ActualPrice

	// Returning results
	a.Assets[op.Base] = assetStatus
	a.Timestamp = time.Now().UnixMicro()
	a.AccountId = uuid.NewString()
	return a, nil
}

func (a LocalAccountFTS) GetOperation(mms model.MiniMarketStats) (model.Operation, error) {
	asset := mms.Asset
	assetStatus, found := a.Assets[asset]
	if !found {
		msg := fmt.Sprintf(logger.FTS_ERR_ASSET_NOT_FOUND, asset)
		logrus.Error(msg)
		return model.Operation{}, model.NewCtbError(msg, true)
	}

	lastOpType := assetStatus.LastOperationType
	lastOpPrice := assetStatus.LastOperationPrice
	currentAmnt := assetStatus.Amount
	currentAmntUsdt := assetStatus.Usdt
	currentPrice := mms.LastPrice

	ftsConfig, err := get_fts_config(config.GetStrategyConfig())
	if err != nil {
		return model.Operation{}, err
	}

	sellPrice := get_threshold_rate(lastOpPrice, ftsConfig.SellThreshold)
	stopLossPrice := get_threshold_rate(lastOpPrice, utils.SignChangeDecimal(ftsConfig.StopLossThreshold))
	buyPrice := get_threshold_rate(lastOpPrice, utils.SignChangeDecimal(ftsConfig.BuyThreshold))
	missProfitPrice := get_threshold_rate(lastOpPrice, ftsConfig.MissProfitThreshold)

	if lastOpType == OP_BUY_FTS && currentPrice.GreaterThanOrEqual(sellPrice) {
		// sell command
		operationInit := build_operation_init(asset, currentAmnt, currentPrice)
		logrus.WithField("comp", "fts").
			Infof(logger.FTS_TRADE, "SELL", asset, lastOpType, lastOpPrice, currentPrice)
		return build_sell_op(a, operationInit)

	} else if lastOpType == OP_BUY_FTS && currentPrice.LessThanOrEqual(stopLossPrice) {
		// stop loss command
		operationInit := build_operation_init(asset, currentAmnt, currentPrice)
		logrus.WithField("comp", "fts").
			Infof(logger.FTS_TRADE, "STOP_LOSS", asset, lastOpType, lastOpPrice, currentPrice)
		return build_sell_op(a, operationInit)

	} else if lastOpType == OP_SELL_FTS && currentPrice.LessThanOrEqual(buyPrice) {
		// buy command
		operationInit := build_operation_init(asset, currentAmntUsdt, currentPrice)
		logrus.WithField("comp", "fts").
			Infof(logger.FTS_TRADE, "BUY", asset, lastOpType, lastOpPrice, currentPrice)
		return build_buy_op(a, operationInit)

	} else if lastOpType == OP_SELL_FTS && currentPrice.GreaterThanOrEqual(missProfitPrice) {
		// miss profit command
		operationInit := build_operation_init(asset, currentAmntUsdt, currentPrice)
		logrus.WithField("comp", "fts").
			Infof(logger.FTS_TRADE, "MISS_PROFIT", asset, lastOpType, lastOpPrice, currentPrice)
		return build_buy_op(a, operationInit)
	}

	logrus.WithField("comp", "fts").Debugf(logger.FTS_TRADE, "NO_OP", asset, lastOpType, lastOpPrice, currentPrice)
	return model.Operation{}, nil
}

func build_operation_init(asset string, amount, price decimal.Decimal) operation_init {
	return operation_init{
		asset:       asset,
		amount:      amount,
		targetPrice: price}
}

func build_buy_op(laccount LocalAccountFTS, operationInit operation_init) (model.Operation, error) {
	symbol := utils.GetSymbolFromAsset(operationInit.asset)
	limit, found := laccount.SpotMarketLimits[symbol]
	if !found {
		msg := fmt.Sprintf(logger.FTS_ERR_SPOT_MARKET_SIZE_NOT_FOUND, symbol)
		logrus.Error(msg)
		return model.Operation{}, model.NewCtbError(msg, true)
	}

	if operationInit.amount.LessThan(limit.MinQuote) {
		logrus.WithField("comp", "fts").Infof(logger.FTS_BELOW_QUOTE_LIMIT,
			symbol, model.BUY, operationInit.amount, model.QUOTE_AMOUNT, limit.MinQuote)
		return model.Operation{}, nil
	}

	op := model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      laccount.ExeId,
		Type:       model.AUTO,
		Base:       operationInit.asset,
		Quote:      "USDT",
		Side:       model.BUY,
		Amount:     operationInit.amount,
		AmountSide: model.QUOTE_AMOUNT,
		Price:      operationInit.targetPrice,
		Status:     model.PENDING}

	logrus.WithField("comp", "fts").
		Infof(logger.FTS_OPERATION, op.Base, op.Quote, op.Amount, op.AmountSide, op.Side)
	return op, nil
}

func build_sell_op(laccount LocalAccountFTS, operationInit operation_init) (model.Operation, error) {
	symbol := utils.GetSymbolFromAsset(operationInit.asset)
	limit, found := laccount.SpotMarketLimits[symbol]
	if !found {
		msg := fmt.Sprintf(logger.FTS_ERR_SPOT_MARKET_SIZE_NOT_FOUND, symbol)
		logrus.Error(msg)
		return model.Operation{}, model.NewCtbError(msg, true)
	}

	if operationInit.amount.LessThan(limit.MinBase) {
		logrus.WithField("comp", "fts").Infof(logger.FTS_BELOW_BASE_LIMIT,
			symbol, model.SELL, operationInit.amount, model.BASE_AMOUNT, limit.MinBase)
		return model.Operation{}, nil
	}

	op := model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      laccount.ExeId,
		Type:       model.AUTO,
		Base:       operationInit.asset,
		Quote:      "USDT",
		Side:       model.SELL,
		Amount:     operationInit.amount,
		AmountSide: model.BASE_AMOUNT,
		Price:      operationInit.targetPrice,
		Status:     model.PENDING}

	logrus.WithField("comp", "fts").
		Infof(logger.FTS_OPERATION, op.Base, op.Quote, op.Amount, op.AmountSide, op.Side)
	return op, nil
}

func get_threshold_rate(price decimal.Decimal, percentage decimal.Decimal) decimal.Decimal {
	abs := percentage.Abs()
	sign := percentage.Div(abs).Round(8)
	delta := price.Div(decimal.NewFromInt(100)).Mul(abs).Round(8)
	return price.Add(delta.Mul(sign)).Round(8)
}

func init_asset_status_fts(rbalance model.RemoteBalance, price model.AssetPrice) AssetStatusFTS {
	return AssetStatusFTS{
		Asset:              rbalance.Asset,
		Amount:             rbalance.Amount,
		Usdt:               decimal.Zero,
		LastOperationType:  OP_BUY_FTS,
		LastOperationPrice: price.Price}
}
