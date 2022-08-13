package utils

import (
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

const (
	_MAX_NUM       = "1000000000000"  // 10^12
	_MIN_NUM       = "-1000000000000" // -10^12
	_INTEGER_REGEX = `^[1-9][0-9]*$`
	_FLOAT_REGEX   = `^(([1-9][0-9]+)|[0-9])(\.[0-9]+)?$`
)

/*** Symbols ***/

func IsSymbolTradable(symbol string) bool {
	symbolre, _ := regexp.Compile(`^.+USDT$`)
	return symbolre.Match([]byte(symbol))
}

func GetSymbolFromAsset(base string) (string, errors.CtbError) {
	// base is already a symbol
	symbolre, _ := regexp.Compile(`^.*USDT$`)
	if symbolre.Match([]byte(base)) {
		return "", errors.Internal(logger.UTILS_ERR_ASSET_TO_SYMBOL, base)
	}

	return base + "USDT", nil
}

func GetSymbolsFromAssets(bases []string) ([]string, []errors.CtbError) {
	symbols := make([]string, 0, len(bases))
	errors := make([]errors.CtbError, 0)
	for _, base := range bases {
		symbol, err := GetSymbolFromAsset(base)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		symbols = append(symbols, symbol)
	}
	return symbols, errors
}

func GetAssetFromSymbol(symbol string) (string, errors.CtbError) {
	if !IsSymbolTradable(symbol) {
		return "", errors.Internal(logger.UTILS_ERR_SYMBOL_TO_ASSET, symbol)
	}

	return strings.TrimSuffix(symbol, "USDT"), nil
}

func GetAssetsFromSymbols(symbols []string) ([]string, []errors.CtbError) {
	assets := make([]string, 0, len(symbols))
	errors := make([]errors.CtbError, 0)
	for _, symbol := range symbols {
		asset, err := GetAssetFromSymbol(symbol)
		if err != nil {
			errors = append(errors, err)
		}
		assets = append(assets, asset)
	}
	return assets, errors
}

/*** Decimals ***/

func GetFloatRegexp() string {
	return _FLOAT_REGEX
}

func GetIntegerRegexp() string {
	return _INTEGER_REGEX
}

func DecimalFromString(str string) decimal.Decimal {
	decimal, err := decimal.NewFromString(str)
	if err != nil {
		logrus.Panicf(logger.UTILS_ERR_FAILED_TO_DECODE_DECIMAL, str)
	}
	return decimal
}

func DecimalFromFloat64(num float64) decimal.Decimal {
	return decimal.NewFromFloat(num)
}

func SignChangeDecimal(d decimal.Decimal) decimal.Decimal {
	return decimal.NewFromInt(-1).Mul(d).Round(8)
}

func MaxDecimal() decimal.Decimal {
	return DecimalFromString(_MAX_NUM)
}

func MinDecimal() decimal.Decimal {
	return DecimalFromString(_MIN_NUM)
}

func IncrementByPercentage(amt decimal.Decimal, perc decimal.Decimal) decimal.Decimal {
	if perc == decimal.Zero {
		return amt
	}

	percabs := perc.Abs()
	sign := perc.Div(percabs).Round(8)
	amtabs := amt.Abs()
	delta := amtabs.Div(decimal.NewFromInt(100)).Mul(percabs).Round(8)
	return amt.Add(delta.Mul(sign)).Round(8)
}

func PercentageOf(amt decimal.Decimal, perc decimal.Decimal) decimal.Decimal {
	if perc == decimal.Zero {
		return decimal.Zero
	}

	amt = amt.Abs()
	perc = perc.Abs()
	return amt.Div(decimal.NewFromInt(100)).Mul(perc).Round(8)
}

// mmss vs assets
func ToAssets(mmss []model.MiniMarketStats) []string {
	assets := make([]string, 0, len(mmss))
	for _, mms := range mmss {
		assets = append(assets, mms.Asset)
	}
	return assets
}
