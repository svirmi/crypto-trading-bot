package executions

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

// Gets active execution object by execution id.
// Returns the execution object, if found, an empty execution
// object if nothing was found, or an error was thrown.
// Returns an error if computation failed
func GetCurrentlyActive(exeId string) (model.Execution, error) {
	exe, err := FindCurrentlyActiveByExeId(exeId)
	if err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// Creates an new execution based on a remote wallet object, or
// restors a previous active execution.
// Returns the newly created execution object, or the active execution
// object found in DB.
// Returns an error if computation failed or an error was thrown.
func CreateOrRestore(raccount model.RemoteAccount) (model.Execution, error) {
	// Get current active execution from DB
	exe, err := FindCurrentlyActive()
	if err != nil {
		return model.Execution{}, err
	}

	// Active execution found, restoring it
	if !exe.IsEmpty() {
		log.Printf("restoring execution %s", exe.ExeId)
		log.Printf("execution status: %s", exe.Status)
		log.Printf("assets to be traded: %v", exe.Symbols)
		return exe, nil
	}

	// No active execution found, starting a new one
	if raccount.IsEmpty() {
		return model.Execution{}, fmt.Errorf("empty remote account received")
	}

	exe = buildExecution(raccount)
	log.Printf("starting execution %s", exe.ExeId)
	log.Printf("assets to be traded: %v", exe.Symbols)
	if InsertOne(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// Changes execution status to PAUSED
// STARTED --> PAUSED allowed
// RESUMED --> PAUSED allowed
// PAUSED --> PAUSED forbidden
// TERMINATED --> PAUSED forbidden
// Once the execution is paused, the bot will stop automatic
// trading of cryptocurrencies and will allow manual operations.
// Returns the modified execution object or an empty execution
// object if checks failed, or an error was thrown.
// Returns an error if computation failed or checks did not
// succeed
func Pause(exeId string) (model.Execution, error) {
	exe, err := FindCurrentlyActiveByExeId(exeId)
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
	if err := InsertOne(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// Changes the execution statues to RESUMED
// PAUSED --> RESUMED allowed
// STARTED --> RESUMED allowed
// RESUMED --> RESUMED allowed
// TERMINATED --> RESUMED forbidden
// Once the execution is resumed, the bot will start trading
// cryptocurrencies and manual intervention will be no longer
// allowed.
// Returns the modified execution object or an empty execution
// object if checks failed, or an error was thrown.
// Returns an error if computation failed or checks did not
// succeed
func Resume(exeId string) (model.Execution, error) {
	exe, err := FindCurrentlyActiveByExeId(exeId)
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
	if err := InsertOne(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// Changes the execution status to TERMINATED
// STARTED --> TERMINATED allowed
// RESUMED --> TERMINATED allowed
// PAUSED --> TERMINATED allowed
// TERMINATED --> TERMINATED forbidden
// Once the execution is terminated, it can not be resumed.
// Cryptocurrencies are sold into USDT and to resume
// automatic trading, a new execution will have to be created.
// Returns the modified execution object or an empty execution
// object if checks failed, or an error was thrown.
// Returns an error if computation failed or checks did not
// succeed
func Terminate(exeId string) (model.Execution, error) {
	exe, err := FindCurrentlyActiveByExeId(exeId)
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
	if err := InsertOne(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

func buildExecution(account model.RemoteAccount) model.Execution {
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
