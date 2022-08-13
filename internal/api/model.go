package api

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

// Errors
type ctb_error_dto struct {
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e ctb_error_dto) is_empty() bool {
	return reflect.DeepEqual(e, ctb_error_dto{})
}

func to_ctb_error_dto(err errors.CtbError) ctb_error_dto {
	var httpStatus int
	switch err.Code() {
	case errors.BAD_REQUEST_ERROR:
		httpStatus = http.StatusBadRequest
	case errors.DUPLICATE_ERROR:
		httpStatus = http.StatusConflict
	case errors.NOT_FOUND_ERROR:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	return ctb_error_dto{
		Message: err.Error(),
		Status:  httpStatus,
	}
}

type ctb_response_dto struct {
	Body   interface{}
	Status int
}

func exe_to_ctb_response_dto(httpStatus int, exe model.Execution) ctb_response_dto {
	var exe_dto exe_res_dto
	exe_dto.ExeId = exe.ExeId
	exe_dto.Assets = exe.Assets
	exe_dto.Status = string(exe.Status)
	exe_dto.StrategyType = string(exe.StrategyType)
	exe_dto.Props = exe.Props
	exe_dto.Timestamp = exe.Timestamp

	var response ctb_response_dto
	response.Status = httpStatus
	response.Body = exe_dto
	return response
}

// Ping pong
type ping_pong_res_dto struct {
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
