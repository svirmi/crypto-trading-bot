package api

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

// Errors
type error_dto struct {
	ErrorMessage string `json:"errorMessage"`
}

func (a error_dto) is_empty() bool {
	return reflect.DeepEqual(a, error_dto{})
}

// Ping pong
type ping_pong_dto struct {
	Ping string `json:"ping"`
}

// Executions
type exe_create_req_dto struct {
	StrategyType   string            `json:"strategyType"`
	StrategyConfig map[string]string `json:"strategyConfig"`
}

func (e exe_create_req_dto) Validate() error {
	if e.StrategyType == "" {
		err := fmt.Errorf(logger.API_ERR_FIELD_REQUIRED, "strategyType")
		logrus.Errorf(err.Error())
		return err
	}
	if e.StrategyConfig == nil {
		err := fmt.Errorf(logger.API_ERR_FIELD_REQUIRED, "strategyConfig")
		logrus.Errorf(err.Error())
		return err
	}
	if len(e.StrategyConfig) == 0 {
		err := fmt.Errorf(logger.API_ERR_COLLECTION_EMPTY, "strategyConfig")
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}

type exe_update_req_dto struct {
	Status string `json:"status"`
}

func (e exe_update_req_dto) Validate() error {
	if e.Status == "" {
		err := fmt.Errorf(logger.API_ERR_FIELD_REQUIRED, "status")
		logrus.Errorf(err.Error())
		return err
	}
	status := model.ExeStatus(e.Status)
	if status != model.EXE_ACTIVE && status != model.EXE_TERMINATED {
		err := fmt.Errorf(logger.API_ERR_VALUE_UNKNOWN, "status", status)
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}

type exe_res_dto struct {
	ExeId        string            `json:"exeId"`
	Status       string            `json:"status"`
	Assets       []string          `json:"assets"`
	StrategyType string            `json:"strategyType"`
	Props        map[string]string `json:"props"`
	Timestamp    int64             `json:"timestamp"`
}

func to_exe_res_dto(exe model.Execution) exe_res_dto {
	var e exe_res_dto
	e.ExeId = exe.ExeId
	e.Assets = exe.Assets
	e.Status = string(exe.Status)
	e.StrategyType = string(exe.StrategyType)
	e.Props = exe.Props
	e.Timestamp = exe.Timestamp
	return e
}
