package main

import (
	"log"
	"time"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/executions"
	"github.com/valerioferretti92/trading-bot-demo/internal/laccount"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

func main() {
	defer shutdown()

	raccount, err := binance.GetAccout()
	if err != nil {
		log.Fatalf(err.Error())
	}

	exe, err := executions.CreateOrRestore(raccount)
	if err != nil {
		log.Fatalf(err.Error())
	}

	tradableSymbols := binance.FilterTradableSymbols(exe.Symbols)
	prices, err := binance.GetAssetsValueUsdt(tradableSymbols)
	if err != nil {
		log.Fatalf(err.Error())
	}

	laCreationRequest := model.LocalAccountInit{
		ExeId:               exe.ExeId,
		RAccount:            raccount,
		StrategyType:        model.FIXED_THRESHOLD_STRATEGY,
		TradableAssetsPrice: prices}
	_, err = laccount.CreateOrRestore(laCreationRequest)
	if err != nil {
		log.Fatalf(err.Error())
	}

	symbols := make(map[string]bool)
	symbols["ETHUSDT"] = true
	symbols["BTCUSDT"] = true
	binance.MiniMarketsStatsServe(symbols)
	time.Sleep(30 * time.Second)
}

func shutdown() {
	binance.Close()
	log.Printf("bye, bye")
}
