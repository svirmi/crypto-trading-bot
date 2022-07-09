package binance

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func TestCanSpotTrade(t *testing.T) {
	logger.Initialize(false, true, true)
	old := symbols

	defer func() {
		symbols = old
	}()

	exchange := binance_exchange{}
	symbols = get_symbols()

	testutils.AssertTrue(t, exchange.CanSpotTrade("BTCUSDT"), "can_spot_trade")
	testutils.AssertFalse(t, exchange.CanSpotTrade("SHIBAUSDT"), "can_spot_trade")
	testutils.AssertFalse(t, exchange.CanSpotTrade("SHITUSDT"), "can_spot_trade")
}

func TestGetSpotMarketLimits(t *testing.T) {
	logger.Initialize(false, true, true)
	old := symbols

	defer func() {
		symbols = old
	}()

	exchange := binance_exchange{}
	symbols = get_symbols()

	btcusdt := symbols["BTCUSDT"]
	btcusdt.Filters = make([]map[string]interface{}, 0, 3)
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters[0]["filterType"] = "LOT_SIZE"
	btcusdt.Filters[0]["minQty"] = "0.001"
	btcusdt.Filters[0]["maxQty"] = "999.999"
	btcusdt.Filters[0]["stepSize"] = "0.1"
	btcusdt.Filters[1]["filterType"] = "MARKET_LOT_SIZE"
	btcusdt.Filters[1]["minQty"] = "0.002"
	btcusdt.Filters[1]["maxQty"] = "999.998"
	btcusdt.Filters[1]["stepSize"] = "0.05"
	btcusdt.Filters[2]["filterType"] = "MIN_NOTIONAL"
	btcusdt.Filters[2]["minNotional"] = "10.00"
	symbols["BTCUSDT"] = btcusdt

	exp := model.SpotMarketLimits{
		MinBase:  utils.DecimalFromString("0.002"),
		MaxBase:  utils.DecimalFromString("999.998"),
		StepBase: utils.DecimalFromString("0.1"),
		MinQuote: utils.DecimalFromString("10.00")}
	got, err := exchange.GetSpotMarketLimits("BTCUSDT")

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "spot_market_limits")
}

func TestGetSpotMarketLimits_InvalidValues(t *testing.T) {
	logger.Initialize(false, true, true)
	old := symbols

	defer func() {
		symbols = old
	}()

	exchange := binance_exchange{}
	symbols = get_symbols()

	btcusdt := symbols["BTCUSDT"]
	btcusdt.Filters = make([]map[string]interface{}, 0, 3)
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters[0]["filterType"] = "LOT_SIZE"
	btcusdt.Filters[0]["minQty"] = ""
	btcusdt.Filters[0]["maxQty"] = "999.999"
	btcusdt.Filters[0]["stepSize"] = "0.1"
	btcusdt.Filters[1]["filterType"] = "MARKET_LOT_SIZE"
	btcusdt.Filters[1]["minQty"] = "0.002"
	btcusdt.Filters[1]["stepSize"] = "0.05"
	btcusdt.Filters[2]["filterType"] = "MIN_NOTIONAL"
	btcusdt.Filters[2]["minNotional"] = "10.00"
	symbols["BTCUSDT"] = btcusdt

	exp := model.SpotMarketLimits{
		MinBase:  utils.DecimalFromString("0.002"),
		MaxBase:  utils.DecimalFromString("999.999"),
		StepBase: utils.DecimalFromString("0.1"),
		MinQuote: utils.DecimalFromString("10.00")}
	got, err := exchange.GetSpotMarketLimits("BTCUSDT")

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "spot_market_limits")

	btcusdt = symbols["BTCUSDT"]
	btcusdt.Filters = make([]map[string]interface{}, 0, 3)
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters[0]["filterType"] = "LOT_SIZE"
	btcusdt.Filters[1]["filterType"] = "MARKET_LOT_SIZE"
	btcusdt.Filters[2]["filterType"] = "MIN_NOTIONAL"
	symbols["BTCUSDT"] = btcusdt

	exp = model.SpotMarketLimits{
		MinBase:  decimal.Zero,
		MaxBase:  utils.MaxDecimal(),
		StepBase: decimal.Zero,
		MinQuote: decimal.Zero}
	got, err = exchange.GetSpotMarketLimits("BTCUSDT")

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "spot_market_limits")
}

func TestGetSpotMarketLimits_FilterNotFound(t *testing.T) {
	logger.Initialize(false, true, true)
	old := symbols

	defer func() {
		symbols = old
	}()

	exchange := binance_exchange{}
	symbols = get_symbols()

	btcusdt := symbols["BTCUSDT"]
	btcusdt.Filters = make([]map[string]interface{}, 0, 3)
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters = append(btcusdt.Filters, make(map[string]interface{}))
	btcusdt.Filters[0]["filterType"] = "LOT_SIZE"
	btcusdt.Filters[0]["minQty"] = "0.001"
	btcusdt.Filters[0]["maxQty"] = "999.999"
	btcusdt.Filters[0]["stepSize"] = "0.1"
	btcusdt.Filters[1]["filterType"] = "MIN_NOTIONAL"
	btcusdt.Filters[1]["minNotional"] = "10.00"
	symbols["BTCUSDT"] = btcusdt

	got, err := exchange.GetSpotMarketLimits("BTCUSDT")

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "spot_market_limits")
}
