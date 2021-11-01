package laccount

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

// Creates a local account based on the remote account, or restores
// account already in DB linked to a currently active execution
// identified by exeId.
// Returns local account object or an empty local account if an
// error was thrown.
// Returns an error if computation failed
// TRUSTS that exeId correspond to an active execution
func CreateOrRestore(exeId string, raccount model.RemoteAccount) (model.LocalAccount, error) {
	// Get current local account from DB by execution id
	laccount, err := FindLatest(exeId)
	if err != nil {
		return model.LocalAccount{}, err
	}

	// Restore existing local account
	if !laccount.IsEmpty() {
		log.Printf("restoring local account %s", laccount.AccountId)
		return laccount, nil
	}

	// Create new local account
	laccount, err = buildLocalAccount(exeId, raccount)
	if err != nil {
		return model.LocalAccount{}, nil
	}
	if laccount.IsEmpty() {
		err = fmt.Errorf("empty local account built from remote account %v", raccount)
		return laccount, err
	}
	if len(laccount.Balances) == 0 {
		log.Printf("local account with empty symbol map")
	}
	if err := Insert(laccount); err != nil {
		return model.LocalAccount{}, err
	}
	log.Printf("registering local account %s", laccount.AccountId)
	return laccount, nil
}

// Updates local account identified by the execution id exeId
// with an executed operation
// Returns local account object or an empty local account if an
// error was thrown or checks did not succeed.
// Returns an error if computation failed or checks did not succeed
// TRUSTS that exeId correspond to an active execution
func RecordTradingOperation(exeId string, operation model.Operation) (model.LocalAccount, error) {
	// Getting local account from DB
	laccount, err := FindLatest(exeId)
	if err != nil {
		return model.LocalAccount{}, err
	}

	// Checking prerequisites
	if laccount.IsEmpty() {
		err = fmt.Errorf("no local account for execution id %s could be found", exeId)
		return model.LocalAccount{}, err
	}
	if operation.Type != model.OP_BUY && operation.Type != model.OP_SELL {
		err = fmt.Errorf("operation type %s unsupported", operation.Type)
		return model.LocalAccount{}, err
	}
	if operation.Base == "USDT" || operation.Quote != "USDT" {
		err = fmt.Errorf("bad operation: %v", operation)
		return model.LocalAccount{}, err
	}

	// Updating local account
	laccount, err = recordTradingOperation(operation, laccount)
	if err != nil {
		return model.LocalAccount{}, err
	}

	// Storing updated local account in DB
	err = Insert(laccount)
	if err != nil {
		return model.LocalAccount{}, err
	}
	return laccount, nil
}

// Gets latest version of local wallet linked to the execution
// identified by exeId
// Returns local account object or an empty local account if an
// error was thrown
// Returns an error if computation failed or checks did not succeed
// TRUSTS that exeId correspond to an active execution
func GetLatest(exeId string) (model.LocalAccount, error) {
	laccount, err := FindLatest(exeId)

	if err != nil {
		return model.LocalAccount{}, err
	}
	if laccount.IsEmpty() {
		err = fmt.Errorf("no local account for execution id %s could be found", exeId)
		return model.LocalAccount{}, err
	}
	return laccount, nil
}

func recordTradingOperation(operation model.Operation, laccount model.LocalAccount) (model.LocalAccount, error) {
	asset := operation.Base
	lbalance := laccount.Balances[asset]
	if lbalance.IsEmpty() {
		err := fmt.Errorf("asset %s not managed in local wallet", asset)
		return model.LocalAccount{}, err
	}

	if operation.Type == model.OP_BUY {
		lbalance.Amount = lbalance.Amount + operation.Actual.BaseQty
		lbalance.Usdt = lbalance.Usdt - operation.Actual.QuoteQty
	} else {
		lbalance.Amount = lbalance.Amount - operation.Actual.BaseQty
		lbalance.Usdt = lbalance.Usdt + operation.Actual.QuoteQty
	}
	lbalance.OperationIds = append(lbalance.OperationIds, operation.OpId)

	if lbalance.Amount < 0 || lbalance.Usdt < 0 {
		err := fmt.Errorf("trying to update balance with negative values: %v", lbalance)
		return model.LocalAccount{}, err
	}
	laccount.Balances[asset] = lbalance
	return laccount, nil
}

func buildLocalAccount(exeId string, raccount model.RemoteAccount) (model.LocalAccount, error) {
	balances := make(map[string]model.LocalBalance)
	for _, balance := range raccount.Balances {
		stramount := balance.Amount
		amount, err := strconv.ParseFloat(stramount, 32)
		if err != nil {
			return model.LocalAccount{}, err
		}

		lbalance := model.LocalBalance{
			Asset:        balance.Asset,
			Amount:       float32(amount),
			Usdt:         0,
			OperationIds: []string{}}
		balances[lbalance.Asset] = lbalance
	}

	laccount := model.LocalAccount{
		AccountId: uuid.NewString(),
		ExeId:     exeId,
		Balances:  balances,
		Timestamp: time.Now().UnixMilli()}
	return laccount, nil
}
