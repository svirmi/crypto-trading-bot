package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/valerioferretti92/crypto-trading-bot/internal/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/handler"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
)

// TODO: min / max LOT_SIZE
func main() {
	defer shutdown()
	sigc := interrupt_handler()

	config.ParseConfig()
	mongodb.Initialize()
	binance.Initialize()

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
	strategyType := model.StrategyType(strategyConfig.Type)
	laCreationRequest := model.LocalAccountInit{
		ExeId:               exe.ExeId,
		RAccount:            raccount,
		StrategyType:        strategyType,
		TradableAssetsPrice: prices}
	laccount, err := laccount.CreateOrRestore(laCreationRequest)
	if err != nil {
		log.Fatalf(err.Error())
	}

	mms := make(chan []model.MiniMarketStats)
	handler.InitTradingContext(laccount, exe)
	handler.InitMmsChannel(mms)
	binance.InitMmsChannel(mms)

	handler.HandleMiniMarketsStats()
	binance.MiniMarketsStatsServe(tradableAssets)

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
	binance.MiniMarketsStatsStop()
	mongodb.Disconnect()
	log.Printf("bye, bye")
}
