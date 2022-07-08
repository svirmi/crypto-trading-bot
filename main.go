package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/analytics"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/exchange/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/exchange/local"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/handler"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
)

var (
	exchange model.IExchange
	exe      model.Execution
)

func main() {
	// Parsing command line
	envstr := flag.String("env", string(model.MAINNET), "if present, application runs on testnet")
	v := flag.Bool("v", false, "if present, debug logs are shown")
	vv := flag.Bool("vv", false, "if present, trace logs are shown")
	colors := flag.Bool("colors", false, "if present, logs are colored")
	flag.Parse()

	// Initalizing logger
	var level logrus.Level = logrus.InfoLevel
	if *vv {
		level = logrus.TraceLevel
	} else if *v {
		level = logrus.DebugLevel
	}
	logger.Initialize(*colors, level)

	// Register interrupt handler
	env := model.ParseEnv(*envstr)
	register_interrupt_handler(env)

	// Getting exchange instance
	if model.SIMULATION == env {
		exchange = local.GetExchange()
	} else if model.TESTNET == env || model.MAINNET == env {
		exchange = binance.GetExchange()
	} else {
		logrus.Panicf(logger.MAIN_ERR_UNSUPPORTED_ENV, env)
	}

	// Parsing config
	err := config.Initialize(env)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Initializing mongodb
	err = mongodb.Initialize()
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Instanciating channels
	var mmsch chan []model.MiniMarketStats
	var cllch chan model.MiniMarketStatsAck
	mmsch = make(chan []model.MiniMarketStats)
	if env == model.SIMULATION {
		cllch = make(chan model.MiniMarketStatsAck, 10)
	}

	// Initializing exchange
	err = exchange.Initialize(mmsch, cllch)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Retrieving remote account
	raccount, err := exchange.GetAccout()
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Creating or restoring execution
	strategyConfig := config.GetStrategyConfig()
	exeReq := model.ExecutionInit{
		Raccount:     raccount,
		StrategyType: model.StrategyType(strategyConfig.Type),
		Props:        strategyConfig.Config}
	exe, err = executions.CreateOrRestore(exeReq)
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
	strategyType := model.StrategyType(strategyConfig.Type)
	laccReq := model.LocalAccountInit{
		ExeId:               exe.ExeId,
		RAccount:            raccount,
		StrategyType:        strategyType,
		TradableAssetsPrice: prices}
	_, err = laccount.CreateOrRestore(laccReq)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Initializing handler
	handler.Initialize(mmsch, cllch, exchange)

	// Handling price updates
	handler.HandleMiniMarketsStats()
	exchange.MiniMarketsStatsServe()

	// Wait until the application is stopped
	select {}
}

func register_interrupt_handler(env model.Env) chan os.Signal {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-sigc

		if env != model.SIMULATION {
			terminate()
		} else {
			terminate_simulation()
		}

		logrus.Info("bye, bye")
		os.Exit(0)
	}()
	return sigc
}

func terminate() {
	if exchange != nil {
		exchange.MiniMarketsStatsStop()
	}
	handler.Terminate()
	mongodb.Disconnect()
}

func terminate_simulation() {
	if exchange != nil {
		exchange.MiniMarketsStatsStop()
	}
	handler.Terminate()
	executions.Terminate(exe.ExeId)
	analytics.StoreAnalytics(exe.ExeId)
	mongodb.Disconnect()
}
