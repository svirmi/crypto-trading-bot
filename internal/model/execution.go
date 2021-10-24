package model

import "reflect"

const (
	EXE_INIT   = "EXE_INIT"   // Execution started
	EXE_UPDATE = "EXE_UPDATE" // Execution modification peformed in manual mode
	EXE_DONE   = "EXE_DONE"   // Execution terminated
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
