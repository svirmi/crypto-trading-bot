package fts

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"github.com/valerioferretti92/trading-bot-demo/internal/utils"
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

func InitLocalAccountFTS(creationRequest model.LocalAccountInit) (LocalAccountFTS, error) {
	var ignored = make(map[string]float32)
	var assets = make(map[string]AssetStatusFTS)

	for _, rbalance := range creationRequest.RAccount.Balances {
		symbol := utils.GetSymbolFromAsset(rbalance.Asset)
		price, found := creationRequest.TradableAssetsPrice[symbol]
		if !found {
			ignored[rbalance.Asset] = rbalance.Amount
			continue
		}
		assetStatus, err := init_asset_status_FTS(rbalance, price)
		if err != nil {
			return LocalAccountFTS{}, err
		}
		assets[rbalance.Asset] = assetStatus
	}

	return LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        creationRequest.ExeId,
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMilli()},
		Ignored: ignored,
		Assets:  assets}, nil
}

func init_asset_status_FTS(rbalance model.RemoteBalance, price model.SymbolPrice) (AssetStatusFTS, error) {
	return AssetStatusFTS{
		Asset:             rbalance.Asset,
		Amount:            rbalance.Amount,
		Usdt:              0,
		LastOperationType: OP_BUY_FTS,
		LastOperationRate: price.Price}, nil
}
