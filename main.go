package main

import (
	"github.com/valerioferretti92/trading-bot-demo/internal/repository"
)

func main() {
	defer repository.Disconnect()

	// Parsing command line

	//binance.BookTickerServe()
	//binance.SendMarketOrder("ETH", "USDT", 1)
	//time.Sleep(2 * time.Second)
	//binance.GetAccout()

	//binance.SendMarketOrder("BTC", "BNB", 0.1)
	//time.Sleep(2 * time.Second)
	//binance.GetAccout()
	//binance.GetExchangeInfo()
	repository.Ping()
}
