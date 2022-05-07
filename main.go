package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/exchange/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/handler"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
)

type cmdline_args struct {
	env      string
	colors   bool
	logLevel logrus.Level
}

// TODO: hanlde big orders in the exchange package
// TODO: check assets with 0 balance
func main() {
	// Parsing command line
	args := parse_cmdline()

	// Register interrupt handler
	exchange := binance.GetExchange()
	register_interrupt_handler(exchange)

	// Initializing logger
	logger.Initialize(args.colors, args.logLevel)

	// Parsing config
	config.Initialize(model.ParseEnv(args.env))

	// Initializing mongodb
	err := mongodb.Initialize()
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Initializing exchange
	mms := make(chan []model.MiniMarketStats)
	err = exchange.Initialize(mms)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Retrieving remote account
	raccount, err := exchange.GetAccout()
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Creating or restoring execution
	exe, err := executions.CreateOrRestore(raccount)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Getting tradable assets
	tradableAssets := exchange.FilterTradableAssets(exe.Assets)
	prices, err := exchange.GetAssetsValue(tradableAssets)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Creating or restoring local account
	strategyConfig := config.GetStrategyConfig()
	strategyType := model.StrategyType(strategyConfig.Type)
	req := model.LocalAccountInit{
		ExeId:               exe.ExeId,
		RAccount:            raccount,
		StrategyType:        strategyType,
		TradableAssetsPrice: prices}
	lacc, err := laccount.CreateOrRestore(req)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Initializing handler
	handler.Initialize(lacc, exe, mms, exchange)

	// Handling price updates
	handler.HandleMiniMarketsStats()
	exchange.MiniMarketsStatsServe(tradableAssets)

	// Wait until the application is stopped
	select {}
}

func register_interrupt_handler(exchange model.IExchange) chan os.Signal {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-sigc
		exchange.MiniMarketsStatsStop()
		mongodb.Disconnect()
		logrus.Info("bye, bye")
		os.Exit(0)
	}()
	return sigc
}

func parse_cmdline() cmdline_args {
	env := flag.String("env", string(model.MAINNET), "if present, application runs on testnet")
	v := flag.Bool("v", false, "if present, debug logs are shown")
	vv := flag.Bool("vv", false, "if present, trace logs are shown")
	colors := flag.Bool("colors", false, "if present, logs are colored")
	flag.Parse()

	var level logrus.Level = logrus.InfoLevel
	if *vv {
		level = logrus.TraceLevel
	} else if *v {
		level = logrus.DebugLevel
	}

	return cmdline_args{
		env:      *env,
		colors:   *colors,
		logLevel: level,
	}
}
