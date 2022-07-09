package model

import (
	"reflect"

	"github.com/shopspring/decimal"
)

type IExchange interface {
	Initialize(chan []MiniMarketStats, chan MiniMarketStatsAck) error
	CanSpotTrade(string) bool
	GetSpotMarketLimits(string) (SpotMarketLimits, error)
	FilterTradableAssets([]string) []string
	GetAssetsValue([]string) (map[string]AssetPrice, error)
	GetAccout() (RemoteAccount, error)
	SendSpotMarketOrder(Operation) (Operation, error)
	MiniMarketsStatsServe() error
	MiniMarketsStatsStop()
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
