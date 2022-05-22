package laccount

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/ds"
)

// Creates a local account based on the remote account, or restores
// account already in DB linked to a currently active execution
// identified by exeId.
// Returns local account object or an empty local account if an
// error was thrown.
// Returns an error if computation failed
// TRUSTS that exeId corresponds to an active execution
func CreateOrRestore(req model.LocalAccountInit) (model.ILocalAccount, error) {
	// Get current local account from DB by execution id
	laccount, err := find_latest_by_exeId(req.ExeId)
	if err != nil {
		return nil, err
	}

	// Restore existing local account
	if laccount != nil && laccount.GetStrategyType() != req.StrategyType {
		err = fmt.Errorf(logger.LACC_ERR_STRATEGY_MISMATCH,
			req.ExeId, req.StrategyType, laccount.GetAccountId(), laccount.GetStrategyType())
		logrus.Error(err.Error())
		return nil, err
	}
	if laccount != nil && laccount.GetStrategyType() == req.StrategyType {
		logrus.Infof(logger.LACC_RESTORE, laccount.GetAccountId())
		return laccount, nil
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

func Create(laccout model.ILocalAccount) error {
	return insert(laccout)
}

func GetLatestByExeId(exeId string) (model.ILocalAccount, error) {
	return find_latest_by_exeId(exeId)
}

func initialise_local_account(req model.LocalAccountInit) (model.ILocalAccount, error) {
	if len(req.RAccount.Balances) == 0 {
		err := fmt.Errorf(logger.LACC_ERR_EMPTY_RACC)
		logrus.Error(err.Error())
		return nil, err
	}

	var laccount model.ILocalAccount = nil
	if req.StrategyType == model.DEMO_STRATEGY {
		laccount = ds.LocalAccountDS{}
	} else {
		err := fmt.Errorf(logger.LACC_ERR_UNKNOWN_STRATEGY, req.StrategyType)
		logrus.Error(err.Error())
		return nil, err
	}
	laccount, err := laccount.Initialize(req)
	if err != nil {
		return nil, err
	}
	return laccount, nil
}
