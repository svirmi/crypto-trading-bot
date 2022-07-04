package model

import "reflect"

type ExeStatus string

const (
	EXE_ACTIVE     ExeStatus = "EXE_ACTIVE"
	EXE_TERMINATED ExeStatus = "EXE_TERMINATED"
)

type Execution struct {
	ExeId        string            `bson:"exeId"`
	Status       ExeStatus         `bson:"status"`
	Assets       []string          `bson:"assets"`
	StrategyType StrategyType      `bson:"strategyType"`
	Props        map[string]string `bson:"props"`
	Timestamp    int64             `bson:"timestamp"`
}

func (e Execution) IsEmpty() bool {
	return reflect.DeepEqual(e, Execution{})
}

type ExecutionInit struct {
	Raccount     RemoteAccount
	StrategyType StrategyType
	Props        map[string]string
}

func (e ExecutionInit) IsEmpty() bool {
	return reflect.DeepEqual(e, ExecutionInit{})
}
