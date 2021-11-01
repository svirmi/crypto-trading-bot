package operations

import (
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

// Inserts an init operation for each owned crypto currency,
// if the it finds no operations linked to the execution exeId.
// If it finds a non empty list operations, it assumes that
// they match remote and local wallets nad it does nothing.
// Returns the list of operations it creates, nil if none was
// created or an error occurred.
// Returns an error if computation failed.
// TRUSTS that exeId referes to a non terminated execution
// TRUSTS that, if operations are found in DB, they match
// local and remote wallet
func Initialize(exeId string, raccount model.RemoteAccount) ([]model.Operation, error) {
	// Getting operations linked to execution id exeId
	ops, err := FindByExeId(exeId)
	if err != nil {
		return nil, err
	}

	// Returing empty slice if operations linked to exeId were found
	if len(ops) > 0 {
		return []model.Operation{}, nil
	}

	// Creating OP_INIT operations
	for _, rbalance := range raccount.Balances {
		op, err := buildOperation(exeId, rbalance)

		if err != nil {
			return nil, err
		} else {
			log.Printf("regitering init operation, symbol: %s, balance: %f",
				op.Base, op.Actual.BaseQty)
			ops = append(ops, op)
		}
	}

	// Inserting OP_INIT operations in DB
	err = InsertMany(ops)
	if err != nil {
		return nil, err
	}
	return ops, nil
}

func buildOperation(exeId string, rbalance model.RemoteBalance) (model.Operation, error) {
	amount, err := strconv.ParseFloat(rbalance.Amount, 32)
	if err != nil {
		return model.Operation{}, err
	}

	order := model.OrderDetails{
		Rate:     0,
		BaseQty:  float32(amount),
		QuoteQty: 0,
	}

	return model.Operation{
		OpId:      uuid.NewString(),
		ExeId:     exeId,
		Type:      model.OP_INIT,
		Base:      rbalance.Asset,
		Quote:     "USDT",
		Expected:  order,
		Actual:    order,
		Spread:    0,
		Timestamp: time.Now().UnixMilli()}, nil
}
