package fts

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
		Type:   string(model.FIXED_THRESHOLD_STRATEGY),
		Config: exp,
	}
	got, _ := get_fts_config(strategyConfig)

	testutils.AssertEq(t, exp, got, "fts_config")
}
