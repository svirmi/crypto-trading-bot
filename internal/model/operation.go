package model

import (
	"reflect"

	"github.com/google/uuid"
)

const (
	OP_INIT  = "OP_INIT"  // Crypto balance at the start of an execution
	OP_SELL  = "OP_SELL"  // Crypto sell operation
	OP_BUY   = "OP_BUY"   // Crypto buy operation
	OP_CLOSE = "OP_CLOSE" // Crypto sell off at the end of an execution
)

type Operation struct {
	OpId         uuid.UUID `bson:"opId"`         // Operation id
	PreviousId   uuid.UUID `bson:"previousId"`   // Previous operation id
	ExeId        uuid.UUID `bson:"exeId"`        // Execution id
	Type         string    `bson:"type"`         // Operation type
	Symbol       string    `bson:"symbol"`       // Crypto traded
	Price        float32   `bson:"price"`        // Current crypto price
	ForcastedQty float32   `bson:"forcastedQty"` // Qty we expect to get, BUY -> Symbol, SELL -> USDT
	Qty          float32   `bbson:"qty"`         // Quantity actually gotten, BUY -> Symbol, SELL -> USDT
	Timestamp    int64     `bson:"timestamp"`    // Operation timestamp
}

func (o Operation) IsEmpty() bool {
	return reflect.DeepEqual(o, Operation{})
}
