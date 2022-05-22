package model

import (
	"reflect"

	"github.com/shopspring/decimal"
)

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
	Amount decimal.Decimal
}

func (b RemoteBalance) IsEmpty() bool {
	return reflect.DeepEqual(b, RemoteBalance{})
}

type AssetPrice struct {
	Asset string
	Price decimal.Decimal
}

func (p AssetPrice) IsEmpty() bool {
	return reflect.DeepEqual(p, AssetPrice{})
}

// Strategy types
type StrategyType string

const (
	DEMO_STRATEGY StrategyType = "DEMO_STRATEGY"
)

type ILocalAccount interface {
	GetAccountId() string
	GetExeId() string
	GetStrategyType() StrategyType
	GetTimestamp() int64
	Initialize(LocalAccountInit) (ILocalAccount, error)
	RegisterTrading(Operation) (ILocalAccount, error)
	GetOperation(MiniMarketStats, SpotMarketLimits) (Operation, error)
}

// Abstract local account representation
// It contains fields that are common to all strategy dependant
// local wallet representations. To be composed with those
// strategy dependant types
type LocalAccountMetadata struct {
	AccountId    string       `bson:"accountId"`    // Local account object id
	ExeId        string       `bson:"exeId"`        // Execution id this local wallet is bound to
	StrategyType StrategyType `bson:"strategyType"` // Strategy type
	Timestamp    int64        `bson:"timestamp"`    // Timestamp
}

func (a LocalAccountMetadata) GetAccountId() string {
	return a.AccountId
}

func (a LocalAccountMetadata) GetExeId() string {
	return a.ExeId
}

func (a LocalAccountMetadata) GetStrategyType() StrategyType {
	return a.StrategyType
}

func (a LocalAccountMetadata) GetTimestamp() int64 {
	return a.Timestamp
}

func (a LocalAccountMetadata) IsEmpty() bool {
	return reflect.DeepEqual(a, LocalAccountMetadata{})
}

type LocalAccountInit struct {
	ExeId               string
	RAccount            RemoteAccount
	TradableAssetsPrice map[string]AssetPrice
	StrategyType        StrategyType
}

func (acr LocalAccountInit) IsEmpty() bool {
	return reflect.DeepEqual(acr, LocalAccountInit{})
}

type SpotMarketLimits struct {
	MinBase  decimal.Decimal
	MaxBase  decimal.Decimal
	StepBase decimal.Decimal
	MinQuote decimal.Decimal
}

func (s SpotMarketLimits) IsEmpty() bool {
	return reflect.DeepEqual(s, SpotMarketLimits{})
}
