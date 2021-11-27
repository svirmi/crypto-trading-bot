package model

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

const (
	NO_OP_CMD = "NO_OP_CMD"
	BUY_CMD   = "BUY_CMD"
	SELL_CMD  = "SELL_CMD"
)

type TradingCommand struct {
	Base       string
	Quote      string
	Amount     string
	AmountSide string
}
