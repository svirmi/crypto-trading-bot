package utils

import (
	"log"
	"strings"

	"github.com/shopspring/decimal"
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
		log.Fatalf("failed to convert string \"%s\" to decimal", str)
	}
	return decimal
}

func SignChangeDecimal(d decimal.Decimal) decimal.Decimal {
	return decimal.NewFromInt(-1).Mul(d)
}
