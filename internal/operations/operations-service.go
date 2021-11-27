package operations

import "github.com/valerioferretti92/crypto-trading-bot/internal/model"

// Creates an operation
// Returns an error if computation failed
func Create(op model.Operation) error {
	return Insert(op)
}

// Creates many operations
// Returns an error if computation failed
func CreateMany(ops []model.Operation) error {
	return InsertMany(ops)
}

// Gets operations by execution id
// Returns a slice of operations or nil if an error occurred
// Returns an error if computation failed
func GetByExecutionId(exeId string) ([]model.Operation, error) {
	return FindByExeId(exeId)
}
