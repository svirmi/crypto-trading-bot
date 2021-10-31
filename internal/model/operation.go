package model

import (
	"reflect"
)

const (
	OP_INIT   = "OP_INIT"   // Synch at exeuction start, manual intervention
	OP_SELL   = "OP_SELL"   // Crypto sell operation, crypto asset to USDT
	OP_BUY    = "OP_BUY"    // Crypto buy operation, USDT to crypto asset
	OP_MANUAL = "OP_MANUAL" // Manual operation, arbitrary base, arbitrary quote
	OP_CLOSE  = "OP_DONE"   // Sell off at execution end, manual intervention
)

type OrderDetails struct {
	Rate     float32 `bson:"rate"`     // how much of "quote" to get one unit of "base"
	BaseQty  float32 `bson:"baseQty"`  // baseQty = nil if quoteQty != nil
	QuoteQty float32 `bson:"quoteQty"` // quoteQty = nil if baseQty != nil
}

func (o OrderDetails) IsEmpty() bool {
	return reflect.DeepEqual(o, OrderDetails{})
}

type Operation struct {
	OpId      string       `bson:"opId"`      // Operation id
	ExeId     string       `bson:"exeId"`     // Execution id
	Type      string       `bson:"type"`      // Operation type
	Base      string       `bson:"base"`      // Base crypto
	Quote     string       `bson:"quote"`     // Quote crypto
	Expected  OrderDetails `bson:"expected"`  // Expected order details
	Actual    OrderDetails `bson:"actual"`    // Actual order details
	Spread    float32      `bson:"spread"`    // Spread percentage expected - actual
	Timestamp int64        `bson:"timestamp"` // Operation timestamp
}

func (o Operation) IsEmpty() bool {
	return reflect.DeepEqual(o, Operation{})
}
