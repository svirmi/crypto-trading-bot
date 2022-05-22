package pts

import (
	"testing"

	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

func TestStrategyConfig(t *testing.T) {
	exp :=
		struct {
			BuyPercentage        string
			SellPercentage       string
			BuyAmountPercentage  string
			SellAmountPercentage string
		}{
			BuyPercentage:        "10.05",
			SellPercentage:       "5.45",
			BuyAmountPercentage:  "15",
			SellAmountPercentage: "10",
		}

	strategyConfig := config.StrategyConfig{
		Type:   string(model.PTS_STRATEGY),
		Config: exp}
	got := get_pts_config(strategyConfig)

	testutils.AssertEq(t, exp, got, "pts_config")
}

func TestStrategyConfig_HighPrecision(t *testing.T) {
	exp :=
		struct {
			BuyPercentage        string
			SellPercentage       string
			BuyAmountPercentage  string
			SellAmountPercentage string
		}{
			BuyPercentage:        "10.049445",
			SellPercentage:       "5.4459",
			BuyAmountPercentage:  "15.0001",
			SellAmountPercentage: "9.9999",
		}

	strategyConfig := config.StrategyConfig{
		Type:   string(model.PTS_STRATEGY),
		Config: exp}
	got := get_pts_config(strategyConfig)

	exp.BuyPercentage = "10.05"
	exp.SellPercentage = "5.45"
	exp.BuyAmountPercentage = "15"
	exp.SellAmountPercentage = "10"
	testutils.AssertEq(t, exp, got, "pts_config")
}

func TestStrategyConfig_MismatchingStrategyType(t *testing.T) {
	strategyConfig := config.StrategyConfig{
		Type:   "FAKE_STRATEGY",
		Config: struct{}{},
	}

	testutils.AssertPanic(t, func() {
		get_pts_config(strategyConfig)
	})
}

func TestStrategyConfig_FailedToParseConfig(t *testing.T) {
	exp :=
		struct {
			WrongBuyPercentage        string
			WrongSellPercentage       string
			WrongBuyAmountPercentage  string
			WrongSellAmountPercentage string
		}{
			WrongBuyPercentage:        "10.05",
			WrongSellPercentage:       "5.45",
			WrongBuyAmountPercentage:  "15",
			WrongSellAmountPercentage: "10",
		}

	strategyConfig := config.StrategyConfig{
		Type:   string(model.PTS_STRATEGY),
		Config: exp}

	testutils.AssertPanic(t, func() {
		get_pts_config(strategyConfig)
	})
}

func TestStrategyConfig_ZeroThresholds(t *testing.T) {
	exp :=
		struct {
			BuyPercentage        string
			SellPercentage       string
			BuyAmountPercentage  string
			SellAmountPercentage string
		}{
			BuyPercentage:        "0",
			SellPercentage:       "5.45",
			BuyAmountPercentage:  "15.56",
			SellAmountPercentage: "0",
		}

	strategyConfig := config.StrategyConfig{
		Type:   string(model.PTS_STRATEGY),
		Config: exp}

	testutils.AssertPanic(t, func() {
		get_pts_config(strategyConfig)
	})
}
