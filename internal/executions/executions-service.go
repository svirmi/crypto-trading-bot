package executions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

func CreateOrRestore(raccount model.RemoteAccount) (model.Execution, error) {
	// Get current active execution from DB
	exe, err := find_currently_active()
	if err != nil {
		return model.Execution{}, err
	}

	// Active execution found, restoring it
	if !exe.IsEmpty() {
		logrus.Infof(logger.EXE_RESTORE, exe.ExeId, exe.Status, exe.Assets)
		return exe, nil
	}

	// No active execution found, starting a new one
	exe, err = build_execution(raccount)
	if err != nil {
		return model.Execution{}, err
	}

	logrus.Infof(logger.EXE_START, exe.ExeId, exe.Status, exe.Assets)
	if insert_one(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

func GetLatestByExeId(exeId string) (model.Execution, error) {
	exe, err := find_latest_by_exeId(exeId)
	if err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

func GetCurrentlyActive() (model.Execution, error) {
	return find_currently_active()
}

func Pause(exeId string) (model.Execution, error) {
	exe, err := find_latest_by_exeId(exeId)
	if err != nil {
		return model.Execution{}, err
	}
	if exe.IsEmpty() {
		err = fmt.Errorf(logger.EXE_ERR_NOT_FOUND, exeId)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_PAUSED {
		err = fmt.Errorf(logger.EXE_ERR_STATUS_TRANSITION_NOT_ALLOWED,
			exe.ExeId, model.EXE_PAUSED, model.EXE_PAUSED)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_TERMINATED {
		err = fmt.Errorf(logger.EXE_ERR_STATUS_TRANSITION_NOT_ALLOWED,
			exeId, model.EXE_TERMINATED, model.EXE_PAUSED)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}

	exe.Status = model.EXE_PAUSED
	exe.Timestamp = time.Now().UnixMicro()
	if err := insert_one(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

func Resume(exeId string) (model.Execution, error) {
	exe, err := find_latest_by_exeId(exeId)
	if err != nil {
		return model.Execution{}, err
	}
	if exe.IsEmpty() {
		err = fmt.Errorf(logger.EXE_ERR_NOT_FOUND, exeId)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_ACTIVE {
		err = fmt.Errorf(logger.EXE_ERR_STATUS_TRANSITION_NOT_ALLOWED,
			exe.ExeId, model.EXE_ACTIVE, model.EXE_ACTIVE)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_TERMINATED {
		err = fmt.Errorf(logger.EXE_ERR_STATUS_TRANSITION_NOT_ALLOWED,
			exe.ExeId, model.EXE_TERMINATED, model.EXE_ACTIVE)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}

	exe.Status = model.EXE_ACTIVE
	exe.Timestamp = time.Now().UnixMicro()
	if err := insert_one(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

func Terminate(exeId string) (model.Execution, error) {
	exe, err := find_latest_by_exeId(exeId)
	if err != nil {
		return model.Execution{}, err
	}
	if exe.IsEmpty() {
		err = fmt.Errorf(logger.EXE_ERR_NOT_FOUND, exeId)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_TERMINATED {
		err = fmt.Errorf(logger.EXE_ERR_STATUS_TRANSITION_NOT_ALLOWED,
			exe.ExeId, model.EXE_TERMINATED, model.EXE_TERMINATED)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}

	exe.Status = model.EXE_TERMINATED
	exe.Timestamp = time.Now().UnixMicro()
	if err := insert_one(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

func build_execution(raccount model.RemoteAccount) (model.Execution, error) {
	if len(raccount.Balances) == 0 {
		err := fmt.Errorf(logger.EXE_ERR_EMPTY_RACC)
		logrus.Error(err.Error())
		return model.Execution{}, err
	}

	assets := make([]string, 0, len(raccount.Balances))
	for _, balance := range raccount.Balances {
		assets = append(assets, balance.Asset)
	}

	return model.Execution{
		ExeId:     uuid.NewString(),
		Status:    model.EXE_ACTIVE,
		Assets:    assets,
		Timestamp: time.Now().UnixMicro()}, nil
}
