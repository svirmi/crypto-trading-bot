package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/valerioferretti92/crypto-trading-bot/internal/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
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

	tradableAssets := binance.FilterTradableAssets(exe.Assets)
	prices, err := binance.GetAssetsValue(tradableAssets)
	if err != nil {
		log.Fatalf(err.Error())
	}

	strategyConfig := config.GetStrategyConfig()
	laCreationRequest := model.LocalAccountInit{
		ExeId:               exe.ExeId,
		RAccount:            raccount,
		StrategyType:        strategyConfig.Type,
		TradableAssetsPrice: prices}
	laccount, err := laccount.CreateOrRestore(laCreationRequest)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// NOOP
	mms1 := model.MiniMarketStats{
		Asset:     "BTC",
		LastPrice: 54240.11}
	cmd, _ := laccount.GetCommand(mms1)
	log.Printf("TradingCommand: %+v", cmd)
	// NO0P
	mms1 = model.MiniMarketStats{
		Asset:     "BTC",
		LastPrice: 54242.11}
	cmd, _ = laccount.GetCommand(mms1)
	log.Printf("TradingCommand: %+v", cmd)
	// SELL
	mms1 = model.MiniMarketStats{
		Asset:     "BTC",
		LastPrice: 61000.12}
	cmd, _ = laccount.GetCommand(mms1)
	log.Printf("TradingCommand: %+v", cmd)
	// STOP_LOSS
	mms1 = model.MiniMarketStats{
		Asset:     "BTC",
		LastPrice: 20000.12}
	cmd, _ = laccount.GetCommand(mms1)
	log.Printf("TradingCommand: %+v", cmd)

	binance.MiniMarketsStatsServe([]string{"BTC", "ETH"})

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
	mongodb.Disconnect()
	log.Printf("bye, bye")
}
