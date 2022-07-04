package pts

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
		_BUY_PERCENTAGE:         "10.05",
		_SELL_PERCENTAGE:        "5.454",
		_BUY_AMOUNT_PERCENTAGE:  "15",
		_SELL_AMOUNT_PERCENTAGE: "10"}

	exp := strategy_config_pts{
		BuyPercentage:        utils.DecimalFromString("10.05"),
		SellPercentage:       utils.DecimalFromString("5.45"),
		BuyAmountPercentage:  utils.DecimalFromString("15"),
		SellAmountPercentage: utils.DecimalFromString("10")}

	testutils.AssertEq(t, exp, parse_config(props), "pts_config")
}

func TestStrategyConfig_FailedToParseConfig(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)

	props := map[string]string{
		_BUY_PERCENTAGE:        "10.05",
		_SELL_PERCENTAGE:       "5.454",
		_BUY_AMOUNT_PERCENTAGE: "15",
		"nonExisting":          "10"}

	testutils.AssertPanic(t, func() {
		parse_config(props)
	})
}

func TestStrategyConfig_BelowZeroThresholds(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)

	props := map[string]string{
		_BUY_PERCENTAGE:         "10.05",
		_SELL_PERCENTAGE:        "-1",
		_BUY_AMOUNT_PERCENTAGE:  "15",
		_SELL_AMOUNT_PERCENTAGE: "10"}

	testutils.AssertPanic(t, func() {
		parse_config(props)
	})
}
