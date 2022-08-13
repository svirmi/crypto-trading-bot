package operations

import (
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

func Create(op model.Operation) errors.CtbError {
	err := insert(op)
	if err != nil {
		logrus.Error(err.Error())
	}
	return err
}

func CreateMany(ops []model.Operation) errors.CtbError {
	err := insert_many(ops)
	if err != nil {
		logrus.Error(err.Error())
	}
	return err
}

func GetByExeId(exeId string) ([]model.Operation, errors.CtbError) {
	ops, err := find_by_exe_id(exeId)
	if err != nil {
		logrus.Error(err.Error())
	}
	return ops, err
}
