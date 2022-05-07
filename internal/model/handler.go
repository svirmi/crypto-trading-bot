package model

import (
	"reflect"

	"github.com/shopspring/decimal"
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
