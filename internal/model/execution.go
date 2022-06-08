package model

import "reflect"

type ExeStatus string

const (
	EXE_ACTIVE     ExeStatus = "EXE_ACTIVE"     // Execution started
	EXE_TERMINATED ExeStatus = "EXE_TERMINATED" // Execution terminated
)

type Execution struct {
	ExeId     string    `bson:"exeId"`     // Execution id
	Status    ExeStatus `bson:"status"`    // Execution status
	Assets    []string  `bson:"assets"`    // List of symbols to be traded
	Timestamp int64     `bson:"timestamp"` // Timestamp
}

func (e Execution) IsEmpty() bool {
	return reflect.DeepEqual(e, Execution{})
}
