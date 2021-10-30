package executions

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

type ExecutionCache struct {
	valid bool
	exe   model.Execution
}

var cache ExecutionCache

func GetExecution() (model.Execution, error) {
	if cache.valid {
		return cache.exe, nil
	}

	exe, err := getActiveExecution()
	if err != nil {
		return model.Execution{}, err
	}
	cache.valid = true
	cache.exe = exe
	return exe, nil
}

func CreateOrRestoreExecution() (model.Execution, error) {
	// Invaidating cache
	cache.valid = false

	// Get current active execution from DB
	exe, err := getActiveExecution()
	if err != nil {
		return model.Execution{}, err
	}

	// Active execution found, restoring it
	if !exe.IsEmpty() {
		log.Printf("restoring execution %s", exe.ExeId)
		log.Printf("exe status: %s, exe symbols: %v", exe.Status, exe.Symbols)
		return exe, nil
	}

	// No active execution found, starting a new one
	account, err := binance.GetAccout()
	if err != nil {
		return model.Execution{}, err
	}

	exe = buildExecution(account)
	log.Printf("starting execution %s", exe.ExeId)
	log.Printf("crypto to be traded: %v", exe.Symbols)
	if InsertOneExecution(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// STARTED --> PAUSED allowed
// RESUMED --> PAUSED allowed
// PAUSED --> PAUSED forbidden
// TERMINATED --> PAUSED forbidden
// Once the execution is paused, the bot will stop automatic
// trading of cryptocurrencies and awill llow manual operations.
func PauseExecution() (model.Execution, error) {
	// Invalidating cache
	cache.valid = false

	exe, err := getActiveExecution()
	if err != nil {
		return model.Execution{}, err
	}
	if exe.IsEmpty() {
		err = fmt.Errorf("no active execution found")
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_PAUSED {
		err = fmt.Errorf("execution %s is already PAUSED", exe.ExeId)
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_TERMINATED {
		err = fmt.Errorf("execution %s is TERMINATED. Cannot be paused", exe.ExeId)
		return model.Execution{}, err
	}

	exe.Status = model.EXE_PAUSED
	if err := InsertOneExecution(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// PAUSED --> RESUMED allowed
// STARTED --> RESUMED allowed
// RESUMED --> RESUMED allowed
// TERMINATED --> RESUMED forbidden
// Once the execution is resumed, the bot will start trading
// cryptocurrencies and manual intervention will be no longer
// allowed.
func ResumeExecution() (model.Execution, error) {
	// Invaidating cache
	cache.valid = false

	exe, err := getActiveExecution()
	if err != nil {
		return model.Execution{}, err
	}
	if exe.IsEmpty() {
		err = fmt.Errorf("no active execution found")
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_RESUMED {
		err = fmt.Errorf("execution %s is already RESMUED", exe.ExeId)
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_TERMINATED {
		err = fmt.Errorf("execution %s is TERMINATED. Cannot be resumed", exe.ExeId)
		return model.Execution{}, err
	}

	exe.Status = model.EXE_RESUMED
	if err := InsertOneExecution(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// STARTED --> TERMINATED allowed
// RESUMED --> TERMINATED allowed
// PAUSED --> TERMINATED allowed
// TERMINATED --> TERMINATED forbidden
// Once the execution is terminated, it can not be resumed.
// Cryptocurrencies are sold into USDT and to resume
// automatic trading, a new execution will have to be created.
func TerminateExecution() (model.Execution, error) {
	// Invaidating cache
	cache.valid = false

	exe, err := getActiveExecution()
	if err != nil {
		return model.Execution{}, err
	}
	if exe.IsEmpty() {
		err = fmt.Errorf("no active execution found")
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_TERMINATED {
		err = fmt.Errorf("execution %s is already TERMINATED", exe.ExeId)
		return model.Execution{}, err
	}

	exe.Status = model.EXE_TERMINATED
	if err := InsertOneExecution(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// Builds and returns an execution struct based on
// the account object and whose status is EXE_STARTED
func buildExecution(account model.Account) model.Execution {
	symbols := make([]string, 0, len(account.Balances))
	for i := range account.Balances {
		symbols = append(symbols, account.Balances[i].Asset)
	}

	return model.Execution{
		ExeId:     uuid.NewString(),
		Status:    model.EXE_STARTED,
		Symbols:   symbols,
		Timestamp: time.Now().UnixMilli()}
}

// Returns active execution found read from DB. Empty
// execution struct, if nothing is found
func getActiveExecution() (model.Execution, error) {
	exes, err := FindAllLatestExecution()
	if err != nil {
		return model.Execution{}, err
	}

	var current model.Execution = model.Execution{}
	for i := range exes {
		if exes[i].Status == model.EXE_TERMINATED {
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
