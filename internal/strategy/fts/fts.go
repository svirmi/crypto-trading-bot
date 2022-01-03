package fts

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

type OperationTypeFts string

const (
	OP_BUY_FTS  OperationTypeFts = "OP_BUY_FTS"
	OP_SELL_FTS OperationTypeFts = "OP_SELL_FTS"
)

type AssetStatusFTS struct {
	Asset              string           `bson:"asset"`              // Asset being tracked
	Amount             float32          `bson:"amount"`             // Amount of that asset currently owned
	Usdt               float32          `bson:"usdt"`               // Usdt gotten by selling the asset
	LastOperationType  OperationTypeFts `bson:"lastOperationType"`  // Last FTS operation type
	LastOperationPrice float32          `bson:"lastOperationPrice"` // Asset value at the time last op was executed
}

func (a AssetStatusFTS) IsEmpty() bool {
	return reflect.DeepEqual(a, AssetStatusFTS{})
}

type LocalAccountFTS struct {
	model.LocalAccountMetadata `bson:"metadata"`
	Ignored                    map[string]float32        `bson:"ignored"` // Usdt not to be invested
	Assets                     map[string]AssetStatusFTS `bson:"assets"`  // Value allocation across assets
}

func (a LocalAccountFTS) IsEmpty() bool {
	return reflect.DeepEqual(a, LocalAccountFTS{})
}

type operation_init struct {
	asset       string
	amount      float32
	targetPrice float32
}

func (a LocalAccountFTS) Initialize(creationRequest model.LocalAccountInit) (model.ILocalAccount, error) {
	var ignored = make(map[string]float32)
	var assets = make(map[string]AssetStatusFTS)

	for _, rbalance := range creationRequest.RAccount.Balances {
		price, found := creationRequest.TradableAssetsPrice[rbalance.Asset]
		if !found {
			ignored[rbalance.Asset] = rbalance.Amount
			continue
		}
		assetStatus, err := init_asset_status_fts(rbalance, price)
		if err != nil {
			return nil, err
		}
		assets[rbalance.Asset] = assetStatus
	}

	a = LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        creationRequest.ExeId,
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMilli()},
		Ignored: ignored,
		Assets:  assets}
	return a, nil
}

func (a LocalAccountFTS) RegisterTrading(op model.Operation) (model.ILocalAccount, error) {
	// Check execution ids
	if op.ExeId != a.ExeId {
		err := fmt.Errorf("mismatching execution ids")
		return a, err
	}

	// If the result status is failed, NOP
	if op.Status == model.FAILED {
		return a, nil
	}

	// FTS only handle operation back and forth USDT
	if op.Quote != "USDT" {
		err := fmt.Errorf("FTS can only hande trading to USDT")
		return a, err
	}

	// Getting asset status
	assetStatus, found := a.Assets[op.Base]
	if !found {
		err := fmt.Errorf("asset %s not found in local wallet", op.Base)
		return a, err
	}

	// Updating asset status
	baseAmount := op.Results.BaseAmount
	quoteAmount := op.Results.QuoteAmount
	if op.Side == model.BUY {
		assetStatus.Amount = assetStatus.Amount + baseAmount
		assetStatus.Usdt = assetStatus.Usdt - quoteAmount
		assetStatus.LastOperationType = OP_BUY_FTS
	} else if op.Side == model.SELL {
		assetStatus.Amount = assetStatus.Amount - baseAmount
		assetStatus.Usdt = assetStatus.Usdt + quoteAmount
		assetStatus.LastOperationType = OP_SELL_FTS
	} else {
		err := fmt.Errorf("unsupported operation type %s", op.Type)
		return a, err
	}
	if assetStatus.Amount < 0 || assetStatus.Usdt < 0 {
		err := fmt.Errorf("negative balance detected")
		return a, err
	}
	assetStatus.LastOperationPrice = op.Results.ActualPrice

	// Returning results
	a.Assets[op.Base] = assetStatus
	a.Timestamp = time.Now().UnixMilli()
	return a, nil
}

func (a LocalAccountFTS) GetOperation(mms model.MiniMarketStats) (model.Operation, error) {
	asset := mms.Asset
	assetStatus, found := a.Assets[asset]
	if !found {
		err := fmt.Errorf("asset %s not in local wallet", asset)
		return model.Operation{}, err
	}

	lastOpType := assetStatus.LastOperationType
	lastOpPrice := assetStatus.LastOperationPrice
	currentAmnt := assetStatus.Amount
	currentAmntUsdt := assetStatus.Usdt
	currentPrice := mms.LastPrice

	ftsConfig := get_fts_config()
	sellPrice := get_threshold_rate(lastOpPrice, ftsConfig.SellThreshold)
	stopLossPrice := get_threshold_rate(lastOpPrice, -ftsConfig.StopLossThreshold)
	buyPrice := get_threshold_rate(lastOpPrice, -ftsConfig.BuyThreshold)
	missProfitPrice := get_threshold_rate(lastOpPrice, ftsConfig.MissProfitThreshold)

	if lastOpType == OP_BUY_FTS && currentPrice >= sellPrice {
		// sell command
		operationInit := build_operation_init(asset, currentAmnt/10, currentPrice)
		log_trading_intent("SELL", asset, lastOpPrice, currentPrice)
		return build_sell_op(a, operationInit), nil

	} else if lastOpType == OP_BUY_FTS && currentPrice <= stopLossPrice {
		// stop loss command
		operationInit := build_operation_init(asset, currentAmnt/10, currentPrice)
		log_trading_intent("STOP_LOSS", asset, lastOpPrice, currentPrice)
		return build_sell_op(a, operationInit), nil

	} else if lastOpType == OP_SELL_FTS && currentPrice <= buyPrice {
		// buy command
		operationInit := build_operation_init(asset, currentAmntUsdt/10, currentPrice)
		log_trading_intent("BUY", asset, lastOpPrice, currentPrice)
		return build_buy_op(a, operationInit), nil

	} else if lastOpType == OP_SELL_FTS && currentPrice >= missProfitPrice {
		// miss profit command
		operationInit := build_operation_init(asset, currentAmntUsdt/10, currentPrice)
		log_trading_intent("MISS_PROFIT", asset, lastOpPrice, currentPrice)
		return build_buy_op(a, operationInit), nil

	}

	log_noop(asset, lastOpType, lastOpPrice, currentPrice)
	return model.Operation{}, nil
}

func build_operation_init(asset string, amount float32, price float32) operation_init {
	return operation_init{
		asset:       asset,
		amount:      amount,
		targetPrice: price}
}

func build_buy_op(laccount LocalAccountFTS, operationInit operation_init) model.Operation {
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
		Status:     model.PENDING}
}

func build_sell_op(laccount LocalAccountFTS, operationInit operation_init) model.Operation {
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
		Status:     model.PENDING}
}

func log_noop(asset string, lastOpType OperationTypeFts, lastOpPrice, currentPrice float32) {
	log.Printf("FTS NOOP: asset=%s, lastOpType=%s, lastOpPrice=%f, currentPrice=%f",
		asset, lastOpType, lastOpPrice, currentPrice)
}

func log_trading_intent(cond, asset string, last, current float32) {
	message := fmt.Sprintf("FTS %s condition verified: asset=%s, last=%v, current=%v",
		cond, asset, last, current)
	log.Println(message)
}

func get_threshold_rate(price float32, percentage float32) float32 {
	abs := math.Abs(float64(percentage))
	sign := float64(percentage) / abs
	delta := (float64(price) / 100) * abs
	return price + float32(delta*sign)
}

func init_asset_status_fts(rbalance model.RemoteBalance, price model.AssetPrice) (AssetStatusFTS, error) {
	return AssetStatusFTS{
		Asset:              rbalance.Asset,
		Amount:             rbalance.Amount,
		Usdt:               0,
		LastOperationType:  OP_BUY_FTS,
		LastOperationPrice: price.Price}, nil
}
