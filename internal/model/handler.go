package model

import "reflect"

type MiniMarketStats struct {
	Event       string
	Time        int64
	Asset       string
	LastPrice   float32
	OpenPrice   float32
	HighPrice   float32
	LowPrice    float32
	BaseVolume  float32
	QuoteVolume float32
}

func (m MiniMarketStats) IsEmpty() bool {
	return reflect.DeepEqual(m, MiniMarketStats{})
}

type TradingContext struct {
	Laccount  ILocalAccount
	Execution Execution
}

func (t TradingContext) IsEmpty() bool {
	return reflect.DeepEqual(t, TradingContext{})
}
