package binance

import binanceapi "github.com/adshao/go-binance/v2"

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
