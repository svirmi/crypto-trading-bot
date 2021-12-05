package model

import (
	"log"
	"reflect"
)

// Order type
type OpType string

const (
	AUTO   OpType = "AUTO"
	MANUAL OpType = "MANUAL"
)

// Order side
type OpSide string

const (
	BUY  OpSide = "BUY"
	SELL OpSide = "SELL"
)

func (s OpSide) Invert() OpSide {
	if s == BUY {
		return SELL
	} else if s == SELL {
		return BUY
	} else {
		log.Fatalf("unknown side value %s", s)
		return OpSide("")
	}
}

// Order result status
type OpStatus string

const (
	FILLED           OpStatus = "FILLED"
	PARTIALLY_FILLED OpStatus = "PARTIALLY_FILLED"
	FAILED           OpStatus = "FAILED"
	PENDING          OpStatus = "PENDING"
)

// Order details, amount sides
type AmountSide string

const (
	BASE_AMOUNT  AmountSide = "BASE_AMOUNT"
	QUOTE_AMOUNT AmountSide = "QUOTE_AMOUNT"
)

func (s AmountSide) Invert() AmountSide {
	if s == BASE_AMOUNT {
		return QUOTE_AMOUNT
	} else if s == QUOTE_AMOUNT {
		return BASE_AMOUNT
	} else {
		log.Fatalf("unknown amount side value %s", s)
		return AmountSide("")
	}
}

type OpDetails struct {
	TargetPrice float32    `bson:"targetPrice"` // How much of "quote" to get one unit of "base"
	Amount      float32    `bson:"amount"`      // Amount to be bought or sold
	AmountSide  AmountSide `bson:"amountSide"`  // What amount refers to, base or quote
}

func (o OpDetails) IsEmpty() bool {
	return reflect.DeepEqual(o, OpDetails{})
}

type OpResults struct {
	ActualPrice float32 `bson:"actualPrice"` // Actual rate
	BaseAmount  float32 `bson:"baseAmount"`  // Base amount actually traded
	QuoteAmount float32 `bson:"quoteAmount"` // Quote amount actually traded
}

func (o OpResults) IsEmpty() bool {
	return reflect.DeepEqual(o, OpResults{})
}

type Operation struct {
	OpId      string    `bson:"opId"`      // Operation id
	ExeId     string    `bson:"exeId"`     // Execution id
	Type      OpType    `bson:"opType"`    // Manual vs Auto
	Base      string    `bson:"base"`      // Base crypto
	Quote     string    `bson:"quote"`     // Quote crypto
	Side      OpSide    `bson:"opSide"`    // Buy vs Sell
	Details   OpDetails `bson:"opDetails"` // Expected order details
	Results   OpResults `bson:"opResults"` // Actual order details
	Status    OpStatus  `bson:"status"`    // Status
	Spread    float32   `bson:"spread"`    // Spread percentage expected - actual
	Timestamp int64     `bson:"timestamp"` // Operation timestamp
}

func (o Operation) IsEmpty() bool {
	return reflect.DeepEqual(o, Operation{})
}
