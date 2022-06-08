package model

import (
	"reflect"

	"github.com/shopspring/decimal"
)

type SymbolPrice struct {
	Symbol    string          `bson:"symbol"`
	Price     decimal.Decimal `bson:"price"`
	Timestamp int64           `bson:"timestamp"`
}

func (a SymbolPrice) IsEmpty() bool {
	return reflect.DeepEqual(a, SymbolPrice{})
}

type SymbolPriceByTimestamp struct {
	SymbolPrices []struct {
		Symbol string
		Price  decimal.Decimal
	}
	Timestamp int64
}
