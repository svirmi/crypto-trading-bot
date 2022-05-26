package utils

import (
	"strings"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
)

const (
	_MAX_NUM = "1000000000000"  // 10^12
	_MIN_NUM = "-1000000000000" // -10^12
)

/*** Symbols ***/

func GetSymbolFromAsset(base string) string {
	return base + "USDT"
}

func GetAssetFromSymbol(symbol string) string {
	return strings.TrimSuffix(symbol, "USDT")
}

/*** Decimals ***/

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
