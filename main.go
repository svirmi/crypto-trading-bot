package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/executions"
	"github.com/valerioferretti92/trading-bot-demo/internal/laccount"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"github.com/valerioferretti92/trading-bot-demo/internal/operations"
)

func main() {
	defer shutdown()
	sigc := interrupt_handler()

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

	// Terminate when the application is stopped
	<-sigc
}

func interrupt_handler() chan os.Signal {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	return sigc
}

func shutdown() {
	binance.Close()
	executions.Close()
	operations.Close()
	laccount.Close()
	log.Printf("bye, bye")
}
