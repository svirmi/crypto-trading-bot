package model

import (
	"reflect"
)

const (
	OP_INIT  = "OP_INIT" // Synch at exeuction start, manual intervention
	OP_SELL  = "OP_SELL" // Crypto sell operation
	OP_BUY   = "OP_BUY"  // Crypto buy operation
	OP_CLOSE = "OP_DONE" // Sell off at execution end, manual intervention
)

type Operation struct {
	OpId         string  `bson:"opId"`         // Operation id
	ExeId        string  `bson:"exeId"`        // Execution id
	Type         string  `bson:"type"`         // Operation type
	Symbol       string  `bson:"symbol"`       // Crypto traded
	Price        float32 `bson:"price"`        // Current crypto price
	ForcastedQty float32 `bson:"forcastedQty"` // Qty we expect to get, BUY -> Symbol, SELL -> USDT
	Qty          float32 `bbson:"qty"`         // Quantity actually gotten, BUY -> Symbol, SELL -> USDT
	Timestamp    int64   `bson:"timestamp"`    // Operation timestamp
}

func (o Operation) IsEmpty() bool {
	return reflect.DeepEqual(o, Operation{})
}
