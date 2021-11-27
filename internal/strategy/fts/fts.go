package fts

import (
	"fmt"
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

func (a LocalAccountFTS) GetCommand(model.MiniMarketStats) (model.TradingCommand, error) {
	return model.TradingCommand{}, nil
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

func (a LocalAccountFTS) Initialise(creationRequest model.LocalAccountInit) (model.ILocalAccount, error) {
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

func init_asset_status_FTS(rbalance model.RemoteBalance, price model.AssetPrice) (AssetStatusFTS, error) {
	return AssetStatusFTS{
		Asset:             rbalance.Asset,
		Amount:            rbalance.Amount,
		Usdt:              0,
		LastOperationType: OP_BUY_FTS,
		LastOperationRate: price.Price}, nil
}
