package executions

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

func CreateOrResumeExecution() (model.Execution, error) {
	// Check if execution needs to be resumed
	exes, err := FindAllLatestExecution()
	if err != nil {
		return model.Execution{}, err
	}
	exe, err := getActiveExecution(exes)
	if err != nil {
		return model.Execution{}, err
	}
	if !exe.IsEmpty() {
		log.Printf("resuming execution %s\n", exe.ExeId)
		log.Printf("crypto to be traded: %v", exe.Symbols)
		return exe, nil
	}

	// Getting accout details from bianance
	account, err := binance.GetAccout()
	if err != nil {
		return model.Execution{}, err
	}

	// Starting new execution
	exe = buildExecution(account)
	log.Printf("starting execution %s", exe.ExeId)
	log.Printf("crypto to be traded: %v", exe.Symbols)
	if InsertOneExecution(exe); err != nil {
		return model.Execution{}, err
	}

	return exe, nil
}

func buildExecution(account model.Account) model.Execution {
	symbols := make([]string, 0, len(account.Balances))
	for i := range account.Balances {
		symbols = append(symbols, account.Balances[i].Asset)
	}

	return model.Execution{
		ExeId:     uuid.NewString(),
		Status:    model.EXE_INIT,
		Symbols:   symbols,
		Timestamp: time.Now().UnixMilli()}
}

func getActiveExecution(exes []model.Execution) (model.Execution, error) {
	var current model.Execution = model.Execution{}

	for i := range exes {
		if exes[i].Status == model.EXE_DONE {
			continue
		}
		if current.IsEmpty() {
			current = exes[i]
		} else {
			errorTemplate := "found two executions concurrently active: %s, %s"
			return model.Execution{}, fmt.Errorf(errorTemplate, current.ExeId, exes[i].ExeId)
		}
	}
	return current, nil
}
