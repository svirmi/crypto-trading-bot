package dts

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func TestStrategyConfig(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)

	props := map[string]string{
		_BUY_THRESHOLD:         "12.34",
		_SELL_THRESHOLD:        "23.45",
		_MISS_PROFIT_THRESHOLD: "45.67",
		_STOP_LOSS_THRESHOLD:   "34.561"}

	exp := strategy_config_dts{
		BuyThreshold:        utils.DecimalFromString("12.34"),
		SellThreshold:       utils.DecimalFromString("23.45"),
		MissProfitThreshold: utils.DecimalFromString("45.67"),
		StopLossThreshold:   utils.DecimalFromString("34.56"),
	}

	testutils.AssertEq(t, exp, parse_config(props), "dts_config")
}

func TestStrategyConfig_FailedToParseConfig(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)

	props := map[string]string{
		"nonExisting":          "12.34",
		_SELL_THRESHOLD:        "23.45",
		_MISS_PROFIT_THRESHOLD: "45.67",
		_STOP_LOSS_THRESHOLD:   "34.56"}

	testutils.AssertPanic(t, func() {
		parse_config(props)
	})
}

func TestStrategyConfig_BelowZeroThresholds(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)

	props := map[string]string{
		_BUY_THRESHOLD:         "12.34",
		_SELL_THRESHOLD:        "23.45",
		_MISS_PROFIT_THRESHOLD: "0",
		_STOP_LOSS_THRESHOLD:   "34.56"}

	testutils.AssertPanic(t, func() {
		parse_config(props)
	})
}
