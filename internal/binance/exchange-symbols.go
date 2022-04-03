package binance

import (
	"fmt"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

var CanSpotTrade = func(symbol string) bool {
	status, found := symbols[symbol]

	if !found {
		return false
	}
	return status.Status == string(binanceapi.SymbolStatusTypeTrading) && status.IsSpotTradingAllowed
}

var GetSpotMarketLimits = func(symbol string) (model.SpotMarketLimits, error) {
	iLotSize, err := get_spot_limit_sizes(symbol)
	if err != nil {
		return model.SpotMarketLimits{}, err
	}
	iMarketLotSize, err := get_spot_market_sizes(symbol)
	if err != nil {
		return model.SpotMarketLimits{}, err
	}
	iNotional, err := get_min_notional(symbol)
	if err != nil {
		return model.SpotMarketLimits{}, err
	}

	minBase := decimal.Max(iLotSize.MinBase, iMarketLotSize.MinBase)
	maxBase := decimal.Min(iLotSize.MaxBase, iMarketLotSize.MaxBase)
	stepBase := decimal.Max(iLotSize.StepBase, iMarketLotSize.StepBase)

	return model.SpotMarketLimits{
		MinBase:  minBase,
		MaxBase:  maxBase,
		StepBase: stepBase,
		MinQuote: iNotional}, nil
}

func get_min_notional(symbol string) (decimal.Decimal, error) {
	status, found := symbols[symbol]

	if !found {
		err := fmt.Errorf("exchange symbol %s not found", symbol)
		logrus.WithField("comp", "binance").Error(err.Error())
		return decimal.Zero, err
	}

	iNotional := extract_filter(status.Filters, "MIN_NOTIONAL")
	if iNotional == nil {
		err := fmt.Errorf("MIN_NOTIONAL filter not found for %s", symbol)
		logrus.WithField("comp", "binance").Error(err.Error())
		return decimal.Zero, err
	}

	return parse_number(iNotional["minNotional"], decimal.Zero), nil
}

func get_spot_market_sizes(symbol string) (model.SpotMarketLimits, error) {
	status, found := symbols[symbol]

	if !found {
		err := fmt.Errorf(logger.BINANCE_ERR_SYMBOL_NOT_FOUND, symbol)
		logrus.WithField("comp", "binance").Error(err.Error())
		return model.SpotMarketLimits{}, err
	}

	iMarketLotSize := extract_filter(status.Filters, "MARKET_LOT_SIZE")
	if iMarketLotSize == nil {
		err := fmt.Errorf(logger.BINANCE_ERR_FILTER_NOT_FOUND, "MARKET_LOT_SIZE", symbol)
		logrus.WithField("comp", "binance").Error(err.Error())
		return model.SpotMarketLimits{}, err
	}

	return model.SpotMarketLimits{
		MinBase:  parse_number(iMarketLotSize["minQty"], decimal.Zero),
		MaxBase:  parse_number(iMarketLotSize["maxQty"], utils.MaxDecimal()),
		StepBase: parse_number(iMarketLotSize["stepSize"], decimal.Zero)}, nil
}

func get_spot_limit_sizes(symbol string) (model.SpotMarketLimits, error) {
	status, found := symbols[symbol]

	if !found {
		err := fmt.Errorf(logger.BINANCE_ERR_SYMBOL_NOT_FOUND, symbol)
		logrus.WithField("comp", "binance").Error(err.Error())
		return model.SpotMarketLimits{}, err
	}

	iLotSize := extract_filter(status.Filters, "LOT_SIZE")
	if iLotSize == nil {
		err := fmt.Errorf(logger.BINANCE_ERR_FILTER_NOT_FOUND, "LOT_SIZE", symbol)
		logrus.WithField("comp", "binance").Error(err.Error())
		return model.SpotMarketLimits{}, err
	}

	return model.SpotMarketLimits{
		MinBase:  parse_number(iLotSize["minQty"], decimal.Zero),
		MaxBase:  parse_number(iLotSize["maxQty"], utils.MaxDecimal()),
		StepBase: parse_number(iLotSize["stepSize"], decimal.Zero)}, nil
}

func extract_filter(filters []map[string]interface{}, filterType string) map[string]interface{} {
	for _, filter := range filters {
		if filterType == fmt.Sprintf("%v", filter["filterType"]) {
			return filter
		}
	}
	return nil
}

func parse_number(num interface{}, def decimal.Decimal) decimal.Decimal {
	if num == nil {
		return def
	}
	str := fmt.Sprintf("%s", num)
	if str == "" {
		return def
	}
	return utils.DecimalFromString(str)
}
