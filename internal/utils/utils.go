package utils

import (
	"strings"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
)

func GetSymbolFromAsset(base string) string {
	return base + "USDT"
}

func GetAssetFromSymbol(symbol string) string {
	return strings.TrimSuffix(symbol, "USDT")
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
