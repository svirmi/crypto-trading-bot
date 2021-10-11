package main

import (
	"fmt"
	"log"
	"time"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/handler"
	"github.com/valerioferretti92/trading-bot-demo/internal/repository"
)

func main() {
	defer shutdown()

	account, _ := binance.GetAccout()
	fmt.Printf("%s\n", account.Balances)

	binance.MiniMarketsStatServe(handler.HandlePriceUpdate)
	err := binance.SendMarketOrder("USDT", "ETH", 1000)
	if err != nil {
		log.Printf("%s\n", err)
	}

	account, _ = binance.GetAccout()
	log.Printf("%s\n", account.Balances)

	time.Sleep(10 * time.Second)
}

func shutdown() {
	binance.Close()
	repository.Disconnect()
}
