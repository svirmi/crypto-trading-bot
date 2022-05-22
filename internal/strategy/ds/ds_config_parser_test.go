package ds

import (
	"testing"

	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

func TestStrategyConfig(t *testing.T) {
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
		Type:   string(model.DEMO_STRATEGY),
		Config: exp,
	}
	got := get_ds_config(strategyConfig)

	testutils.AssertEq(t, exp, got, "ds_config")
}

func TestStrategyConfig_MismatchingStrategyType(t *testing.T) {
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
		Type:   "FAKE_STRATEGY",
		Config: exp,
	}

	testutils.AssertPanic(t, func() {
		get_ds_config(strategyConfig)
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
		Type:   string(model.DEMO_STRATEGY),
		Config: exp,
	}

	testutils.AssertPanic(t, func() {
		get_ds_config(strategyConfig)
	})
}

func TestStrategyConfig_ZeroThresholds(t *testing.T) {
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
		Type:   string(model.DEMO_STRATEGY),
		Config: exp,
	}

	testutils.AssertPanic(t, func() {
		get_ds_config(strategyConfig)
	})
}
