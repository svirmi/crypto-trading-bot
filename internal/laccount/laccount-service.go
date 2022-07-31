package laccount

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy"
)

func Create(req model.LocalAccountInit) (model.ILocalAccount, error) {
	// Get current local account from DB by execution id
	laccount, err := find_latest_by_exeId(req.ExeId)
	if err != nil {
		return nil, err
	}
	if laccount != nil {
		err := fmt.Errorf(logger.LACC_ERR_FAILED_TO_CREATE, req.ExeId, laccount.GetAccountId())
		return nil, err
	}

	// Initialise new local account
	laccount, err = initialise_local_account(req)
	if err != nil {
		return nil, err
	}
	if laccount == nil {
		err = fmt.Errorf(logger.LACC_ERR_BUILD_FAILURE)
		logrus.Error(err.Error())
		return nil, err
	}

	// Inseting in mongo db and returning value
	if err := insert(laccount); err != nil {
		return nil, err
	}
	logrus.Infof(logger.LACC_REGISTER, laccount.GetAccountId())
	return laccount, nil
}

func Update(laccount model.ILocalAccount) error {
	return insert(laccount)
}

func GetLatestByExeId(exeId string) (model.ILocalAccount, error) {
	return find_latest_by_exeId(exeId)
}

func GetByExeId(exeId string) ([]model.ILocalAccount, error) {
	return find_by_exeId(exeId)
}

func initialise_local_account(req model.LocalAccountInit) (model.ILocalAccount, error) {
	if len(req.RAccount.Balances) == 0 {
		err := fmt.Errorf(logger.LACC_ERR_EMPTY_RACC)
		logrus.Error(err.Error())
		return nil, err
	}

	laccount, err := strategy.InstanciateLocalAccount(req.StrategyType)
	if err != nil {
		return nil, err
	}

	laccount, err = laccount.Initialize(req)
	if err != nil {
		return nil, err
	}
	return laccount, nil
}
