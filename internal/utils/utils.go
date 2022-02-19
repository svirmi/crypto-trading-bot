package utils

import (
	"strconv"
	"strings"
)

func ParseFloat32(payload string) (float32, error) {
	value, err := strconv.ParseFloat(payload, 32)
	if err != nil {
		return 0, err
	}
	return float32(value), nil
}

func GetSymbolFromAsset(base string) string {
	return base + "USDT"
}

func GetAssetFromSymbol(symbol string) string {
	return strings.TrimSuffix(symbol, "USDT")
}

func Xor(a, b bool) bool {
	return (a || b) && !(a && b)
}
