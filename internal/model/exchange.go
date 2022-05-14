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
	Initialize(mmsch chan []MiniMarketStats) error
	CanSpotTrade(symbol string) bool
	GetSpotMarketLimits(symbol string) (SpotMarketLimits, error)
	FilterTradableAssets(bases []string) []string
	GetAssetsValue(bases []string) (map[string]AssetPrice, error)
	GetAccout() (RemoteAccount, error)
	SendSpotMarketOrder(op Operation) (Operation, error)
	MiniMarketsStatsServe(assets []string) error
	MiniMarketsStatsStop()
}
