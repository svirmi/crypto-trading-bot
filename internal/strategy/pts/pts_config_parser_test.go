package pts

import (
	"testing"

	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func TestStrategyConfig(t *testing.T) {
	logger.Initialize(false, true, true)

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

	config, err := parse_config(props)
	testutils.AssertNil(t, err, "error")
	testutils.AssertEq(t, exp, config, "pts_config")
}

func TestStrategyConfig_FailedToParseConfig(t *testing.T) {
	logger.Initialize(false, true, true)

	props := map[string]string{
		_BUY_PERCENTAGE:        "10.05",
		_SELL_PERCENTAGE:       "5.454",
		_BUY_AMOUNT_PERCENTAGE: "15",
		"nonExisting":          "10"}

	config, err := parse_config(props)
	testutils.AssertNotNil(t, err, "error")
	testutils.AssertEq(t, strategy_config_pts{}, config, "pts_config")
}

func TestStrategyConfig_BelowZeroThresholds(t *testing.T) {
	logger.Initialize(false, true, true)

	props := map[string]string{
		_BUY_PERCENTAGE:         "10.05",
		_SELL_PERCENTAGE:        "-1",
		_BUY_AMOUNT_PERCENTAGE:  "15",
		_SELL_AMOUNT_PERCENTAGE: "10"}

	config, err := parse_config(props)
	testutils.AssertNotNil(t, err, "error")
	testutils.AssertEq(t, strategy_config_pts{}, config, "pts_config")
}
