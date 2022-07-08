package model

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
)

type Env string

const (
	SIMULATION Env = "simulation" // Local exchange
	TESTNET    Env = "testnet"    // Binance testnet
	MAINNET    Env = "mainnet"    // Binance mainnet
)

func ParseEnv(s string) Env {
	if s != string(SIMULATION) && s != string(TESTNET) && s != string(MAINNET) {
		envs := fmt.Sprintf("[%s,%s,%s]", SIMULATION, TESTNET, MAINNET)
		err := fmt.Errorf(logger.MODEL_ERR_UNKNOWN_ENV, s, envs)
		logrus.Panic(err.Error())
	}
	return Env(s)
}

type IExchange interface {
	Initialize(chan []MiniMarketStats, chan MiniMarketStatsAck) error
	CanSpotTrade(string) bool
	GetSpotMarketLimits(string) (SpotMarketLimits, error)
	FilterTradableAssets([]string) []string
	GetAssetsValue([]string) (map[string]AssetPrice, error)
	GetAccout() (RemoteAccount, error)
	SendSpotMarketOrder(Operation) (Operation, error)
	MiniMarketsStatsServe() error
	MiniMarketsStatsStop()
}
