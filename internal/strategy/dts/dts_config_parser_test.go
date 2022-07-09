package dts

import (
	"testing"

	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func TestStrategyConfig(t *testing.T) {
	logger.Initialize(false, true, true)

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

	config, err := parse_config(props)
	testutils.AssertNil(t, err, "error")
	testutils.AssertEq(t, exp, config, "dts_config")
}

func TestStrategyConfig_FailedToParseConfig(t *testing.T) {
	logger.Initialize(false, true, true)

	props := map[string]string{
		"nonExisting":          "12.34",
		_SELL_THRESHOLD:        "23.45",
		_MISS_PROFIT_THRESHOLD: "45.67",
		_STOP_LOSS_THRESHOLD:   "34.56"}

	config, err := parse_config(props)
	testutils.AssertNotNil(t, err, "error")
	testutils.AssertEq(t, strategy_config_dts{}, config, "pts_config")
}

func TestStrategyConfig_BelowZeroThresholds(t *testing.T) {
	logger.Initialize(false, true, true)

	props := map[string]string{
		_BUY_THRESHOLD:         "12.34",
		_SELL_THRESHOLD:        "23.45",
		_MISS_PROFIT_THRESHOLD: "0",
		_STOP_LOSS_THRESHOLD:   "34.56"}

	config, err := parse_config(props)
	testutils.AssertNotNil(t, err, "error")
	testutils.AssertEq(t, strategy_config_dts{}, config, "pts_config")
}
