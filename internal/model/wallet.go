package model

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
)

// Strategy types
type StrategyType string

// Strategies
const (
	DTS_STRATEGY       StrategyType = "DEMO_TRADING_STRATEGY"
	DTS_STRATEGY_SHORT StrategyType = "dts"
	PTS_STRATEGY       StrategyType = "PERCENTAGE_TRADING_STRATEGY"
	PTS_STRATEGY_SHORT StrategyType = "pts"
)

func ParseStr(s string) (StrategyType, error) {
	if s == string(DTS_STRATEGY) || s == string(DTS_STRATEGY_SHORT) {
		return DTS_STRATEGY, nil
	}
	if s == string(PTS_STRATEGY) || s == string(PTS_STRATEGY_SHORT) {
		return PTS_STRATEGY, nil
	}

	envs := fmt.Sprintf("[%s|%s,%s|%s]", DTS_STRATEGY, DTS_STRATEGY_SHORT, PTS_STRATEGY, PTS_STRATEGY_SHORT)
	err := fmt.Errorf(logger.MODEL_ERR_UNKNOWN_ENV, s, envs)
	return StrategyType(s), err
}

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

type AssetAmount struct {
	Asset  string
	Amount decimal.Decimal
}

func (p AssetAmount) IsEmpty() bool {
	return reflect.DeepEqual(p, AssetAmount{})
}

type ILocalAccount interface {
	GetAccountId() string
	GetExeId() string
	GetStrategyType() StrategyType
	GetTimestamp() int64
	Initialize(LocalAccountInit) (ILocalAccount, error)
	RegisterTrading(Operation) (ILocalAccount, error)
	GetOperation(map[string]string, MiniMarketStats, SpotMarketLimits) (Operation, error)
	GetAssetAmounts() map[string]AssetAmount
	ValidateConfig(map[string]string) error
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
