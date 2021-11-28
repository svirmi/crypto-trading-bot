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

const (
	OP_BUY_FTS  = "OP_BUY_FTS"
	OP_SELL_FTS = "OP_SELL_FTS"
)

type AssetStatusFTS struct {
	Asset             string  `bson:"asset"`             // Asset being tracked
	Amount            float32 `bson:"amount"`            // Amount of that asset currently owned
	Usdt              float32 `bson:"usdt"`              // Usdt gotten by selling the asset
	LastOperationType string  `bson:"lastOperationType"` // Last FTS operation type
	LastOperationRate float32 `bson:"lastOperationRate"` // Asset value at the time last op was executed
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

func (a LocalAccountFTS) Initialize(creationRequest model.LocalAccountInit) (model.ILocalAccount, error) {
	var ignored = make(map[string]float32)
	var assets = make(map[string]AssetStatusFTS)

	for _, rbalance := range creationRequest.RAccount.Balances {
		price, found := creationRequest.TradableAssetsPrice[rbalance.Asset]
		if !found {
			ignored[rbalance.Asset] = rbalance.Amount
			continue
		}
		assetStatus, err := init_asset_status_FTS(rbalance, price)
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
	if op.OrderResults.Status == model.FAILED {
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
	baseAmount := op.OrderResults.BaseAmount
	quoteAmount := op.OrderResults.QuoteAmount
	if op.Type == model.OP_BUY_AUTO || op.Type == model.OP_BUY_MANUAL {
		assetStatus.Amount = assetStatus.Amount + baseAmount
		assetStatus.Usdt = assetStatus.Usdt - quoteAmount
		assetStatus.LastOperationType = OP_BUY_FTS
	} else if op.Type == model.OP_SELL_AUTO || op.Type == model.OP_SELL_MANUAL {
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
	assetStatus.LastOperationRate = op.OrderResults.ActualRate

	// Returning results
	a.Assets[op.Base] = assetStatus
	a.Timestamp = time.Now().UnixMilli()
	return a, nil
}

func (a LocalAccountFTS) GetCommand(mms model.MiniMarketStats) (model.TradingCommand, error) {
	asset := mms.Asset
	assetStatus, found := a.Assets[asset]
	if !found {
		err := fmt.Errorf("asset %s not in local wallet", asset)
		return model.TradingCommand{}, err
	}

	lastOpType := assetStatus.LastOperationType
	lastOpRate := assetStatus.LastOperationRate
	currentAmnt := assetStatus.Amount
	currentAmntUsdt := assetStatus.Usdt
	currentRate := mms.LastPrice

	sellPrice := get_threshold_rate(lastOpRate, strategy_config.SellThreshold)
	stopLossPrice := get_threshold_rate(lastOpRate, -strategy_config.StopLossThreshold)
	buyPrice := get_threshold_rate(lastOpRate, -strategy_config.BuyThreshold)
	missProfitPrice := get_threshold_rate(lastOpRate, strategy_config.MissProfitThreshold)

	if lastOpType == OP_BUY_FTS && currentRate >= sellPrice {
		// sell command
		log_trading_intent("SELL", asset, lastOpRate, currentRate)
		return build_sell_command(asset, currentAmnt), nil

	} else if lastOpType == OP_BUY_FTS && currentRate <= stopLossPrice {
		// stop loss command
		log_trading_intent("STOP_LOSS", asset, lastOpRate, currentRate)
		return build_sell_command(asset, currentAmnt), nil

	} else if lastOpType == OP_SELL_FTS && currentRate <= buyPrice {
		// buy command
		log_trading_intent("BUY", asset, lastOpRate, currentRate)
		return build_buy_command(asset, currentAmntUsdt), nil

	} else if lastOpType == OP_SELL_FTS && currentRate >= missProfitPrice {
		// miss profit command
		log_trading_intent("MISS_PROFIT", asset, lastOpRate, currentRate)
		return build_buy_command(asset, currentAmntUsdt), nil

	} else {
		// noop command
		return build_no_op_command(), nil
	}
}

func build_no_op_command() model.TradingCommand {
	return model.TradingCommand{
		CommandType: model.NO_OP_CMD}
}

func build_buy_command(asset string, amount float32) model.TradingCommand {
	return model.TradingCommand{
		Base:        asset,
		Quote:       "USDT",
		Amount:      amount,
		AmountSide:  model.QUOTE_AMOUNT,
		CommandType: model.BUY_CMD}
}

func build_sell_command(asset string, amount float32) model.TradingCommand {
	return model.TradingCommand{
		Base:        asset,
		Quote:       "USDT",
		Amount:      amount,
		AmountSide:  model.BASE_AMOUNT,
		CommandType: model.SELL_CMD}
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

func init_asset_status_FTS(rbalance model.RemoteBalance, price model.AssetPrice) (AssetStatusFTS, error) {
	return AssetStatusFTS{
		Asset:             rbalance.Asset,
		Amount:            rbalance.Amount,
		Usdt:              0,
		LastOperationType: OP_BUY_FTS,
		LastOperationRate: price.Price}, nil
}
