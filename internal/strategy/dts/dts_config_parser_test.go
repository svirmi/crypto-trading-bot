package dts

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

func TestStrategyConfig(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	exp :=
		struct {
			BuyThreshold        string
			SellThreshold       string
			StopLossThreshold   string
			MissProfitThreshold string
		}{
			BuyThreshold:        "12.34",
			SellThreshold:       "23.45",
			StopLossThreshold:   "34.56",
			MissProfitThreshold: "45.67",
		}

	strategyConfig := config.StrategyConfig{
		Type:   string(model.DTS_STRATEGY),
		Config: exp,
	}
	got := get_dts_config(strategyConfig)

	testutils.AssertEq(t, exp, got, "dts_config")
}

func TestStrategyConfig_HighPrecision(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	exp :=
		struct {
			BuyThreshold        string
			SellThreshold       string
			StopLossThreshold   string
			MissProfitThreshold string
		}{
			BuyThreshold:        "12.339",
			SellThreshold:       "23.451",
			StopLossThreshold:   "34.56",
			MissProfitThreshold: "45.665",
		}

	strategyConfig := config.StrategyConfig{
		Type:   string(model.DTS_STRATEGY),
		Config: exp}
	got := get_dts_config(strategyConfig)

	exp.BuyThreshold = "12.34"
	exp.SellThreshold = "23.45"
	exp.StopLossThreshold = "34.56"
	exp.MissProfitThreshold = "45.67"
	testutils.AssertEq(t, exp, got, "pts_config")
}

func TestStrategyConfig_MismatchingStrategyType(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	strategyConfig := config.StrategyConfig{
		Type:   "FAKE_STRATEGY",
		Config: struct{}{},
	}

	testutils.AssertPanic(t, func() {
		get_dts_config(strategyConfig)
	})
}

func TestStrategyConfig_FailedToParseConfig(t *testing.T) {
	exp :=
		struct {
			WrongBuyThreshold        string
			WrongSellThreshold       string
			WrongStopLossThreshold   string
			WrongMissProfitThreshold string
		}{
			WrongBuyThreshold:        "12.34",
			WrongSellThreshold:       "23.45",
			WrongStopLossThreshold:   "34.56",
			WrongMissProfitThreshold: "45.67",
		}

	strategyConfig := config.StrategyConfig{
		Type:   string(model.DTS_STRATEGY),
		Config: exp,
	}

	testutils.AssertPanic(t, func() {
		get_dts_config(strategyConfig)
	})
}

func TestStrategyConfig_ZeroThresholds(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	exp :=
		struct {
			BuyThreshold        string
			SellThreshold       string
			StopLossThreshold   string
			MissProfitThreshold string
		}{
			BuyThreshold:        "0",
			SellThreshold:       "23.45",
			StopLossThreshold:   "34.56",
			MissProfitThreshold: "0",
		}

	strategyConfig := config.StrategyConfig{
		Type:   string(model.DTS_STRATEGY),
		Config: exp,
	}

	testutils.AssertPanic(t, func() {
		get_dts_config(strategyConfig)
	})
}
