package model

import (
	"reflect"

	"github.com/shopspring/decimal"
)

type ExchangeType string

const (
	LOCALEX  ExchangeType = "LOCAL"
	BINANCEX ExchangeType = "BINANCEX"
)

type MiniMarketStats struct {
	Event       string
	Time        int64
	Asset       string
	LastPrice   decimal.Decimal
	OpenPrice   decimal.Decimal
	HighPrice   decimal.Decimal
	LowPrice    decimal.Decimal
	BaseVolume  decimal.Decimal
	QuoteVolume decimal.Decimal
}

func (m MiniMarketStats) IsEmpty() bool {
	return reflect.DeepEqual(m, MiniMarketStats{})
}

type MiniMarketStatsAck struct {
	Count int
}

func (m MiniMarketStatsAck) IsEmpty() bool {
	return reflect.DeepEqual(m, MiniMarketStatsAck{})
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
