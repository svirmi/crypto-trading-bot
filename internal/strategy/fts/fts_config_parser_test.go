package fts

import (
	"testing"

	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

func TestStrategyConfig(t *testing.T) {
	expected :=
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

	gotten := get_fts_config(expected)
	testutils.AssertStructEq(t, expected, gotten)
}
