package laccount

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

func CreateOrRestore(exeId string, raccount model.RemoteAccount) (model.LocalAccount, error) {
	// Get current local account from DB by execution id
	laccount, err := FindLatest(exeId)
	if err != nil {
		return model.LocalAccount{}, err
	}

	// Restore existing local account
	if !laccount.IsEmpty() {
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
	return laccount, nil
}

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

func Get(exeId string) (model.LocalAccount, error) {
	laccount, err := FindLatest(exeId)

	if err != nil {
		return model.LocalAccount{}, err
	}
	if laccount.IsEmpty() {
		err = fmt.Errorf("no local account for execution id %s could be found", exeId)
		return laccount, err
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
		ExeId:     exeId,
		Balances:  balances,
		Timestamp: time.Now().UnixMilli()}
	return laccount, nil
}
