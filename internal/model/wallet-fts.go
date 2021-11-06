package model

import "reflect"

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
	LocalAccountMetadata `bson:"metadata"`
	FrozenUsdt           float32                   `bson:"frozenUsdt"` // Usdt not to be invested
	Assets               map[string]AssetStatusFTS `bson:"assets"`     // Value allocation across assets
}

func (a LocalAccountFTS) IsEmpty() bool {
	return reflect.DeepEqual(a, LocalAccountFTS{})
}
