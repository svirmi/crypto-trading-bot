package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/config"
)

func main() {
	// Parsing command line
	testnet := flag.Bool("testnet", false, "if present, application runs on testnet")
	flag.Parse()

	// Parsing config
	_, err := config.ParseConfig(*testnet)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	binance.New()
	//binance.BookTickerServe()
	//binance.SendMarketOrder("ETH", "USDT", 1)
	//time.Sleep(2 * time.Second)
	//binance.GetAccout()

	binance.SendMarketOrder("BTC", "BNB", 0.1)
	time.Sleep(2 * time.Second)
	binance.GetAccout()
	//binance.GetExchangeInfo()
}
