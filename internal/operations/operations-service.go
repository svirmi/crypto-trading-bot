package operations

import "github.com/valerioferretti92/crypto-trading-bot/internal/model"

func Create(op model.Operation) error {
	return insert(op)
}

func CreateMany(ops []model.Operation) error {
	return insert_many(ops)
}

func GetByExeId(exeId string) ([]model.Operation, error) {
	return find_by_exe_id(exeId)
}
