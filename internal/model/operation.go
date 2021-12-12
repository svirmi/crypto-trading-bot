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

type OpResults struct {
	ActualPrice float32 `bson:"actualPrice"` // Actual rate
	BaseAmount  float32 `bson:"baseAmount"`  // Base amount actually traded
	QuoteAmount float32 `bson:"quoteAmount"` // Quote amount actually traded
	Spread      float32 `bson:"spread"`      // Spread percentage expected - actual
}

func (o OpResults) IsEmpty() bool {
	return reflect.DeepEqual(o, OpResults{})
}

type Operation struct {
	OpId       string     `bson:"opId"`       // Operation id
	ExeId      string     `bson:"exeId"`      // Execution id
	Type       OpType     `bson:"type"`       // Manual vs Auto
	Base       string     `bson:"base"`       // Base crypto
	Quote      string     `bson:"quote"`      // Quote crypto
	Side       OpSide     `bson:"side"`       // Buy vs Sell
	Amount     float32    `bson:"amount"`     // Amount to be bought or sold
	AmountSide AmountSide `bson:"amountSide"` // What amount refers to, base or quote
	Price      float32    `bson:"price"`      // How much of "quote" to get one unit of "base"
	Results    OpResults  `bson:"results"`    // Results
	Status     OpStatus   `bson:"status"`     // Status
	Timestamp  int64      `bson:"timestamp"`  // Operation timestamp
}

func (o Operation) IsEmpty() bool {
	return reflect.DeepEqual(o, Operation{})
}

func (o Operation) Flip() Operation {
	o.Base, o.Quote = o.Quote, o.Base
	o.Side = o.Side.Invert()
	o.AmountSide = o.AmountSide.Invert()
	return o
}
