package utils

import (
	"strings"

	"github.com/shopspring/decimal"
)

func GetSymbolFromAsset(base string) string {
	return base + "USDT"
}

func GetAssetFromSymbol(symbol string) string {
	return strings.TrimSuffix(symbol, "USDT")
}

func SignChangeDecimal(d decimal.Decimal) decimal.Decimal {
	return decimal.NewFromInt(-1).Mul(d)
}
