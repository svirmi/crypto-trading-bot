package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/alecthomas/kong"
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
	"github.com/valerioferretti92/crypto-trading-bot/internal/prices"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy"
)

const (
	config_folder              = "resources"
	simulation_config_filepath = "config-simulation.yaml"
	testnet_config_filepath    = "config-testnet.yaml"
	mainnet_config_filepath    = "config.yaml"
)

type Flags struct {
	v               bool
	vv              bool
	colors          bool
	config_filepath string
}

type Simulate struct {
	StrategyType   string            `arg:"" help:"Strategy type."`
	StrategyConfig map[string]string `arg:"" help:"Strategy config."`
}

func (r *Simulate) Run(flags Flags) error {
	if flags.config_filepath == "" {
		flags.config_filepath = filepath.Join(config_folder, simulation_config_filepath)
	}

	run_simulation(flags, r.StrategyType, r.StrategyConfig)
	return nil
}

type Testnet struct{}

func (r *Testnet) Run(flags Flags) error {
	if flags.config_filepath == "" {
		flags.config_filepath = filepath.Join(config_folder, testnet_config_filepath)
	}
	run(flags)
	return nil
}

type Mainnet struct{}

func (r *Mainnet) Run(flags Flags) error {
	if flags.config_filepath == "" {
		flags.config_filepath = filepath.Join(config_folder, mainnet_config_filepath)
	}
	run(flags)
	return nil
}

var cli struct {
	Debug          bool     `short:"d" help:"Debug level verbosity."`
	Trace          bool     `short:"t" help:"Trace level verbosity."`
	Colors         bool     `short:"c" help:"Enable log colors."`
	ConfigFilepath string   `short:"f" help:"config file path."`
	Simulate       Simulate `cmd:"" help:"Run strategy simulation."`
	Testnet        Testnet  `cmd:"" help:"Run against Binance testnet."`
	Mainnet        Mainnet  `cmd:"" help:"Run against Binance mainnet."`
}

func main() {
	// Parse args and call the Run() method of the selected command
	ctx := kong.Parse(&cli)
	ctx.Run(Flags{cli.Debug, cli.Trace, cli.Colors, cli.ConfigFilepath})
}

func run(flags Flags) {
	defer handle_panics()

	// Initialize logger
	logger.Initialize(flags.colors, flags.v, flags.vv)

	// Register interrupt handler
	register_interrupt_handler()

	// Parsing config files
	err := config.Initialize(flags.config_filepath)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Initializing mongodb
	err = mongodb.Initialize()
	if err != nil {
		logrus.Panic(err.Error())
	}
	terminate_mongodb = func() {
		mongodb.Disconnect()
	}

	// Instanciating channels
	mmsch := make(chan []model.MiniMarketStats)

	// Initializing exchange
	exchange := binance.GetExchange()
	err = exchange.Initialize(mmsch, nil)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Initializing prices
	prices.Initialize()
	terminate_prices = func() {
		prices.Terminate()
	}

	// Initializing handler
	handler.Initialize(mmsch, nil, exchange)
	handler.HandleMiniMarketsStats()

	// Start serving mini markets stats
	exchange.MiniMarketsStatsServe()
	terminate_exchange = func() {
		exchange.MiniMarketsStatsStop()
	}

	// Wait until the application is stopped
	select {}
}

func run_simulation(flags Flags, strategyName string, strategyConfig map[string]string) {
	defer handle_panics()

	// Initialize logger
	logger.Initialize(flags.colors, flags.v, flags.vv)

	// Register interrupt handler
	register_interrupt_handler()

	// Validating configuration
	strategyType, err := model.ParseStr(strategyName)
	if err != nil {
		logrus.Panic(err.Error())
	}
	err = strategy.ValidateConfig(strategyType, strategyConfig)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Parsing config files
	err = config.Initialize(flags.config_filepath)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Initializing mongodb
	err = mongodb.Initialize()
	if err != nil {
		logrus.Panic(err.Error())
	}
	terminate_mongodb = func() {
		mongodb.Disconnect()
	}

	// Instanciating channels
	mmsch := make(chan []model.MiniMarketStats)
	cllch := make(chan model.MiniMarketStatsAck, 10)

	// Initializing exchange
	exchange := local.GetExchange()
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
	exeReq := model.ExecutionInit{
		Raccount:     raccount,
		StrategyType: strategyType,
		Props:        strategyConfig}
	exe, err := executions.CreateOrRestore(exeReq)
	if err != nil {
		logrus.Panic(err.Error())
	}
	terminate_execution = func() {
		executions.Terminate(exe.ExeId)
		analytics.StoreAnalytics(exe.ExeId)
	}

	// Getting tradable assets
	tradableAssets := exchange.FilterTradableAssets(exe.Assets)
	assetPrices, err := exchange.GetAssetsValue(tradableAssets)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Creating or restoring local account
	laccReq := model.LocalAccountInit{
		ExeId:               exe.ExeId,
		RAccount:            raccount,
		StrategyType:        strategyType,
		TradableAssetsPrice: assetPrices}
	_, err = laccount.CreateOrRestore(laccReq)
	if err != nil {
		logrus.Panic(err.Error())
	}

	// Initializing prices
	prices.Initialize()
	terminate_prices = func() {
		prices.Terminate()
	}

	// Initializing handler
	handler.Initialize(mmsch, cllch, exchange)
	handler.HandleMiniMarketsStats()

	// Start serving mini markets stats
	exchange.MiniMarketsStatsServe()
	terminate_exchange = func() {
		exchange.MiniMarketsStatsStop()
	}

	// Wait until the application is stopped
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
		terminate()
	}()
	return sigc
}

func terminate() {
	terminate_exchange()
	terminate_prices()
	terminate_execution()
	terminate_mongodb()
	logrus.Info("bye, bye")
	os.Exit(0)
}

var terminate_exchange = func() {
	// Empty implementation
}

var terminate_prices = func() {
	// Empty implementation
}

var terminate_execution = func() {
	// Empty implementation
}

var terminate_mongodb = func() {
	// Empty implementation
}

func handle_panics() {
	if err := recover(); err != nil {
		fmt.Println(string(debug.Stack()))
		terminate()
	}
}
