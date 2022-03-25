package binance

import (
	"os"
	"testing"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
)

func TestMain(m *testing.M) {
	logger.Initialize(true, logrus.TraceLevel)
	code := m.Run()
	os.Exit(code)
}

func get_symbols() map[string]binanceapi.Symbol {
	symbols = make(map[string]binanceapi.Symbol)
	symbols["BTCUSDT"] = binanceapi.Symbol{
		Status:               string(binanceapi.SymbolStatusTypeTrading),
		IsSpotTradingAllowed: true}
	symbols["ETHUSDT"] = binanceapi.Symbol{
		Status:               string(binanceapi.SymbolStatusTypeTrading),
		IsSpotTradingAllowed: true}
	symbols["DOTUSDT"] = binanceapi.Symbol{
		Status:               string(binanceapi.SymbolStatusTypeTrading),
		IsSpotTradingAllowed: true}
	symbols["SHIBAUSDT"] = binanceapi.Symbol{
		Status:               string(binanceapi.SymbolStatusTypeTrading),
		IsSpotTradingAllowed: false}
	symbols["SHITUSDT"] = binanceapi.Symbol{
		Status:               string(binanceapi.SymbolStatusTypeHalt),
		IsSpotTradingAllowed: true}
	return symbols
}
