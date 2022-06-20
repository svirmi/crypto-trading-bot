package model

import (
	"reflect"

	"github.com/shopspring/decimal"
)

type AnalyticsType string

const (
	EXE_ANALYTICS    AnalyticsType = "execution-analytics"
	OP_ANALYTICS     AnalyticsType = "operation-analytics"
	WALLET_ANALYTICS AnalyticsType = "wallet-analytics"
)

type IAnalytics interface {
	GetExeId() string
	GetAnalyticsType() AnalyticsType
	GetTimestamp() int64
}

type ExeAnalytics struct {
	ExeId         string        `bson:"exeId"`
	AnalyticsType AnalyticsType `bson:"analyticsType"`
	Timestamp     int64         `bson:"timestamp"`
	Assets        []string      `bson:"assets"`
	Status        ExeStatus     `bson:"status"`
}

func (a ExeAnalytics) GetExeId() string {
	return a.ExeId
}

func (a ExeAnalytics) GetAnalyticsType() AnalyticsType {
	return a.AnalyticsType
}

func (a ExeAnalytics) GetTimestamp() int64 {
	return a.Timestamp
}

func (a ExeAnalytics) IsEmpty() bool {
	return reflect.DeepEqual(a, ExeAnalytics{})
}

type OpAnalytics struct {
	ExeId         string          `bson:"exeId"`
	AnalyticsType AnalyticsType   `bson:"analyticsType"`
	Timestamp     int64           `bson:"timestamp"`
	Base          string          `bson:"base"`
	Quote         string          `bson:"quote"`
	Amount        decimal.Decimal `bson:"amount"`
	Side          OpSide          `bson:"side"`
	AmountSide    AmountSide      `bson:"amountSide"`
	Price         decimal.Decimal `bson:"price"`
}

func (a OpAnalytics) GetExeId() string {
	return a.ExeId
}

func (a OpAnalytics) GetAnalyticsType() AnalyticsType {
	return a.AnalyticsType
}

func (a OpAnalytics) GetTimestamp() int64 {
	return a.Timestamp
}

func (a OpAnalytics) IsEmpty() bool {
	return reflect.DeepEqual(a, OpAnalytics{})
}

type WalletAnalytics struct {
	ExeId         string                 `bson:"exeId"`
	AnalyticsType AnalyticsType          `bson:"analyticsType"`
	Timestamp     int64                  `bson:"timestamp"`
	AssetStatuses map[string]AssetStatus `bson:"assetStatuses"`
	WalletValue   decimal.Decimal        `bson:"walletValue"`
}

func (a WalletAnalytics) GetExeId() string {
	return a.ExeId
}

func (a WalletAnalytics) GetAnalyticsType() AnalyticsType {
	return a.AnalyticsType
}

func (a WalletAnalytics) GetTimestamp() int64 {
	return a.Timestamp
}

func (a WalletAnalytics) IsEmpty() bool {
	return reflect.DeepEqual(a, WalletAnalytics{})
}

type AssetStatus struct {
	Asset  string          `bson:"asset"`
	Price  decimal.Decimal `bson:"price"`
	Amount decimal.Decimal `bson:"amount"`
}

func (a AssetStatus) IsEmpty() bool {
	return reflect.DeepEqual(a, AssetStatus{})
}
