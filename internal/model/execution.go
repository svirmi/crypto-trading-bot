package model

import "reflect"

const (
	EXE_ACTIVE     = "EXE_ACTIVE"     // Execution started
	EXE_PAUSED     = "EXE_PAUSED"     // Execution paused, manual operations enabled
	EXE_TERMINATED = "EXE_TERMINATED" // Execution terminated and cryptos sold off
)

type Execution struct {
	ExeId     string   `bson:"exeId"`     // Execution id
	Status    string   `bson:"status"`    // Execution status
	Symbols   []string `bson:"symbols"`   // List of symbols to be traded
	Timestamp int64    `bson:"timestamp"` // Timestamp
}

func (e Execution) IsEmpty() bool {
	return reflect.DeepEqual(e, Execution{})
}
