package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/handler"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func main() {
	// Register interrupt handler
	register_interrupt_handler()

	// Parsing command line
	testnet := flag.Bool("testnet", false, "if present, application runs on testnet")
	colors := flag.Bool("colors", false, "if present, logs are colored")
	v := flag.Bool("v", false, "if present, debug logs are shown")
	vv := flag.Bool("vv", false, "if present, trace logs are shown")
	flag.Parse()

	logger.Initialize(*colors, get_log_level(*v, *vv))
	logrus.Infof(logger.MAIN_LOGICAL_CORES, runtime.NumCPU())

	config.Initialize(*testnet)

	err := mongodb.Initialize()
	if err != nil {
		logrus.Panic(err.Error())
	}

	err = binance.Initialize()
	if err != nil {
		logrus.Panic(err.Error())
	}

	raccount, err := binance.GetAccout()
	if err != nil {
		logrus.Panic(err.Error())
	}

	exe, err := executions.CreateOrRestore(raccount)
	if err != nil {
		logrus.Panic(err.Error())
	}

	tradableAssets := binance.FilterTradableAssets(exe.Assets)
	prices, err := binance.GetAssetsValue(tradableAssets)
	if err != nil {
		logrus.Panic(err.Error())
	}

	spotMarketLimits := make(map[string]model.SpotMarketLimits)
	for _, asset := range tradableAssets {
		symbol := utils.GetSymbolFromAsset(asset)
		spotLimit, err := binance.GetSpotMarketLimits(symbol)
		if err != nil {
			logrus.Panic(err.Error())
		}
		spotMarketLimits[symbol] = spotLimit
	}

	strategyConfig := config.GetStrategyConfig()
	strategyType := model.StrategyType(strategyConfig.Type)
	laCreationRequest := model.LocalAccountInit{
		ExeId:               exe.ExeId,
		RAccount:            raccount,
		StrategyType:        strategyType,
		TradableAssetsPrice: prices,
		SpotMarketLimits:    spotMarketLimits}
	laccount, err := laccount.CreateOrRestore(laCreationRequest)
	if err != nil {
		logrus.Panic(err.Error())
	}

	mms := make(chan []model.MiniMarketStats)
	handler.InitTradingContext(laccount, exe)
	handler.InitMmsChannel(mms)
	binance.InitMmsChannel(mms)

	handler.HandleMiniMarketsStats()
	binance.MiniMarketsStatsServe(tradableAssets)

	// Terminate when the application is stopped
	select {}
}

func register_interrupt_handler() chan os.Signal {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-sigc
		binance.MiniMarketsStatsStop()
		mongodb.Disconnect()
		logrus.Info("bye, bye")
		os.Exit(0)
	}()
	return sigc
}

func get_log_level(v, vv bool) logrus.Level {
	if vv {
		return logrus.TraceLevel
	} else if v {
		return logrus.DebugLevel
	} else {
		return logrus.InfoLevel
	}
}
