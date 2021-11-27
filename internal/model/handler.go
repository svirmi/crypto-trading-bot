package model

// WsMiniMarketsStatEvent define websocket market mini-ticker statistics event
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
