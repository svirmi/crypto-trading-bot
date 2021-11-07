package model

import "reflect"

// Representation of remote Binance wallet
type RemoteAccount struct {
	MakerCommission  int64
	TakerCommission  int64
	BuyerCommission  int64
	SellerCommission int64
	Balances         []RemoteBalance
}

func (a RemoteAccount) IsEmpty() bool {
	return reflect.DeepEqual(a, RemoteAccount{})
}

type RemoteBalance struct {
	Asset  string
	Amount float32
}

func (b RemoteBalance) IsEmpty() bool {
	return reflect.DeepEqual(b, RemoteBalance{})
}

type SymbolPrice struct {
	Symbol string
	Price  float32
}

func (p SymbolPrice) IsEmpty() bool {
	return reflect.DeepEqual(p, SymbolPrice{})
}

const (
	FIXED_THRESHOLD_STRATEGY = "FIXED_THRESHOLD_STRATEGY"
)

type ILocalAccount interface {
	GetAccountId() string
	GetExeId() string
	GetStrategyType() string
	GetTimestamp() int64
}

// Abstract local account representation
// It contains fields that are common to all strategy dependant
// local wallet representations. To be composed with those
// strategy dependant types
type LocalAccountMetadata struct {
	AccountId    string `bson:"accountId"`    // Local account object id
	ExeId        string `bson:"exeId"`        // Execution id this local wallet is bound to
	StrategyType string `bson:"strategyType"` // Strategy type
	Timestamp    int64  `bson:"timestamp"`    // Timestamp
}

func (a LocalAccountMetadata) GetAccountId() string {
	return a.AccountId
}

func (a LocalAccountMetadata) GetExeId() string {
	return a.ExeId
}

func (a LocalAccountMetadata) GetStrategyType() string {
	return a.StrategyType
}

func (a LocalAccountFTS) GetTimestamp() int64 {
	return a.Timestamp
}

func (a LocalAccountMetadata) IsEmpty() bool {
	return reflect.DeepEqual(a, LocalAccountMetadata{})
}

type LocalAccountInit struct {
	ExeId               string
	RAccount            RemoteAccount
	TradableAssetsPrice map[string]SymbolPrice
	StrategyType        string
}

func (acr LocalAccountInit) IsEmpty() bool {
	return reflect.DeepEqual(acr, LocalAccountInit{})
}
