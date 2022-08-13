package laccount

import (
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy"
)

func Create(req model.LocalAccountInit) (model.ILocalAccount, errors.CtbError) {
	// Get current local account from DB by execution id
	laccount, err := find_latest_by_exeId(req.ExeId)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if laccount != nil {
		err := errors.Duplicate(logger.LACC_ERR_FAILED_TO_CREATE, req.ExeId, laccount.GetAccountId())
		logrus.Error(err.Error())
		return nil, err
	}

	// Initialise new local account
	laccount, err = initialise_local_account(req)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if laccount == nil {
		err = errors.Internal(logger.LACC_ERR_BUILD_FAILURE)
		logrus.Error(err.Error())
		return nil, err
	}

	// Inseting in mongo db and returning value
	if err := insert(laccount); err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	logrus.Infof(logger.LACC_REGISTER, laccount.GetAccountId())
	return laccount, nil
}

func Update(update model.ILocalAccount) (model.ILocalAccount, errors.CtbError) {
	lacc, err := find_latest_by_exeId(update.GetExeId())
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if lacc == nil {
		err := errors.NotFound(logger.LACC_ERR_NOT_FOUND, update.GetExeId())
		logrus.Error(err.Error())
		return nil, err
	}

	err = insert(update)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return update, err
}

func GetLatestByExeId(exeId string) (model.ILocalAccount, errors.CtbError) {
	lacc, err := find_latest_by_exeId(exeId)
	if err != nil {
		logrus.Error(err.Error())
	}
	return lacc, err
}

func GetByExeId(exeId string) ([]model.ILocalAccount, errors.CtbError) {
	laccs, err := find_by_exeId(exeId)
	if err != nil {
		logrus.Error(err.Error())
	}
	return laccs, err
}

func initialise_local_account(req model.LocalAccountInit) (model.ILocalAccount, errors.CtbError) {
	if len(req.RAccount.Balances) == 0 {
		err := errors.Internal(logger.LACC_ERR_EMPTY_RACC)
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
