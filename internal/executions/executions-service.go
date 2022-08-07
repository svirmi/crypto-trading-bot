package executions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

func Create(req model.ExecutionInit) (model.Execution, error) {
	// Get current active execution from DB
	exe, err := find_latest()
	if err != nil {
		logrus.Error(err.Error())
		return model.Execution{}, err
	}

	// Active execution found
	if !exe.IsEmpty() && exe.Status == model.EXE_ACTIVE {
		err := fmt.Errorf(logger.EXE_ERR_FAILED_TO_CREATE, exe.ExeId)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}

	// No active execution found, starting a new one
	exe, err = build_execution(req)
	if err != nil {
		logrus.Error(err.Error())
		return model.Execution{}, err
	}

	logrus.Infof(logger.EXE_START, exe.ExeId, exe.Status, exe.Assets)
	if err = insert_one(exe); err != nil {
		logrus.Error(err.Error())
		return model.Execution{}, err
	}
	return exe, nil
}

func GetLatest() (model.Execution, error) {
	exe, err := find_latest()
	if err != nil {
		logrus.Error(err.Error())
	}
	return exe, err
}

func GetLastestByExeId(exeId string) (model.Execution, error) {
	exe, err := find_latest_by_exeId(exeId)
	if err != nil {
		logrus.Error(err.Error())
	}
	return exe, err
}

func GetByExeId(exeId string) ([]model.Execution, error) {
	exe, err := find_by_exeId(exeId)
	if err != nil {
		logrus.Error(err.Error())
	}
	return exe, err
}

func Update(update model.Execution) (model.Execution, error) {
	exe, err := find_latest_by_exeId(update.ExeId)
	if err != nil {
		logrus.Error(err.Error())
		return model.Execution{}, err
	}
	if exe.IsEmpty() {
		err = fmt.Errorf(logger.EXE_ERR_NOT_FOUND, update.ExeId)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}

	if update.Status == exe.Status {
		return exe, nil
	}

	if update.Status == model.EXE_ACTIVE {
		err = fmt.Errorf(logger.EXE_ERR_STATUS_TRANSITION_NOT_ALLOWED,
			exe.ExeId, exe.Status, model.EXE_ACTIVE)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}

	exe.Status = update.Status
	exe.Timestamp = time.Now().UnixMicro()
	if err := insert_one(exe); err != nil {
		logrus.Error(err.Error())
		return model.Execution{}, err
	}
	return exe, nil
}

func build_execution(req model.ExecutionInit) (model.Execution, error) {
	raccount := req.Raccount

	if len(raccount.Balances) == 0 {
		err := fmt.Errorf(logger.EXE_ERR_EMPTY_RACC)
		return model.Execution{}, err
	}

	var usdt bool = false
	assets := make([]string, 0, len(raccount.Balances))
	for _, balance := range raccount.Balances {
		assets = append(assets, balance.Asset)
		if balance.Asset == "USDT" {
			usdt = true
		}
	}
	if !usdt {
		assets = append(assets, "USDT")
	}

	return model.Execution{
		ExeId:        uuid.NewString(),
		Status:       model.EXE_ACTIVE,
		Assets:       assets,
		StrategyType: req.StrategyType,
		Props:        req.Props,
		Timestamp:    time.Now().UnixMicro()}, nil
}
