package main

import (
	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/handler"
	"github.com/valerioferretti92/trading-bot-demo/internal/repository"
)

func main() {
	defer shutdown()

	binance.MiniMarketsStatServe(handler.HandlePriceUpdate)
	binance.GetAccout()
}

func shutdown() {
	binance.Close()
	repository.Disconnect()
}
