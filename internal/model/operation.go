package model

import (
	"reflect"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
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
	}
	logrus.Panicf(logger.MODEL_ERR_UNKNOWN_OP_SIDE, s)
	return ""
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
	}
	logrus.Panicf(logger.MODEL_ERR_UNKNOWN_AMOUNT_SIDE, s)
	return ""
}

type OpResults struct {
	ActualPrice decimal.Decimal `bson:"actualPrice"` // Actual rate
	BaseDiff    decimal.Decimal `bson:"baseAmount"`  // Base amount actually traded
	QuoteDiff   decimal.Decimal `bson:"quoteAmount"` // Quote amount actually traded
	Spread      decimal.Decimal `bson:"spread"`      // Spread percentage exp - actual
}

func (o OpResults) IsEmpty() bool {
	return reflect.DeepEqual(o, OpResults{})
}

type Operation struct {
	OpId       string          `bson:"opId"`       // Operation id
	ExeId      string          `bson:"exeId"`      // Execution id
	Type       OpType          `bson:"type"`       // Manual vs Auto
	Base       string          `bson:"base"`       // Base crypto
	Quote      string          `bson:"quote"`      // Quote crypto
	Side       OpSide          `bson:"side"`       // Buy vs Sell
	Amount     decimal.Decimal `bson:"amount"`     // Amount to be bought or sold
	AmountSide AmountSide      `bson:"amountSide"` // What amount refers to, base or quote
	Price      decimal.Decimal `bson:"price"`      // How much of "quote" to get one unit of "base"
	Results    OpResults       `bson:"results"`    // Results
	Status     OpStatus        `bson:"status"`     // Status
	Timestamp  int64           `bson:"timestamp"`  // Operation timestamp
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
