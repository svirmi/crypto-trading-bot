package model

import (
	"reflect"
)

// Operation types
const (
	OP_BUY_AUTO    = "OP_BUY_AUTO"
	OP_SELL_AUTO   = "OP_SELL_AUTO"
	OP_BUY_MANUAL  = "OP_BUY_MANUAL"
	OP_SELL_MANUAL = "OP_SELL_MANUAL"
)

// Order details, amount sides
const (
	BASE_AMOUNT  = "BASE_AMOUNT"
	QUOTE_AMOUNT = "QUOTE_AMOUNT"
)

// Order results, status
const (
	FULLY_FILLED    = "FULLY_FILLED"
	PARTIALLY_FIELD = "PARTIALLY_FIELD"
	FAILED          = "FAILED"
)

type OrderDetails struct {
	Rate       float32 `bson:"rate"`     // How much of "quote" to get one unit of "base"
	Amount     float32 `bson:"amount"`   // Amount to be bought or sold
	AmountSide string  `bson:"quoteQty"` // What amount refers to, base or quote
}

func (o OrderDetails) IsEmpty() bool {
	return reflect.DeepEqual(o, OrderDetails{})
}

type OrderResults struct {
	ActualRate  float32 `bson:"actualRate"`  // Actual rate
	BaseAmount  float32 `bson:"baseAmount"`  // Base amount actually traded
	QuoteAmount float32 `bson:"quoteAmount"` // Quote amount actually traded
	Status      string  `bson:"status"`      // Status
}

func (o OrderResults) IsEmpty() bool {
	return reflect.DeepEqual(o, OrderResults{})
}

type Operation struct {
	OpId         string       `bson:"opId"`         // Operation id
	ExeId        string       `bson:"exeId"`        // Execution id
	Type         string       `bson:"type"`         // Operation type
	Base         string       `bson:"base"`         // Base crypto
	Quote        string       `bson:"quote"`        // Quote crypto
	OrderDetails OrderDetails `bson:"orderDetails"` // Expected order details
	OrderResults OrderResults `bson:"orderResults"` // Actual order details
	Spread       float32      `bson:"spread"`       // Spread percentage expected - actual
	Timestamp    int64        `bson:"timestamp"`    // Operation timestamp
}

func (o Operation) IsEmpty() bool {
	return reflect.DeepEqual(o, Operation{})
}
