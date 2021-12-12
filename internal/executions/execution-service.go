package executions

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

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
		log.Printf("assets to be traded: %v", exe.Assets)
		return exe, nil
	}

	// No active execution found, starting a new one
	if raccount.IsEmpty() {
		return model.Execution{}, fmt.Errorf("empty remote account received")
	}

	exe = build_execution(raccount)
	log.Printf("starting execution %s", exe.ExeId)
	log.Printf("assets to be traded: %v", exe.Assets)
	if InsertOne(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// Gets active execution object by execution id.
// Returns the execution object, if found, an empty execution
// object if nothing was found, or an error was thrown.
// Returns an error if computation failed
func GetCurrentlyActiveByExeId(exeId string) (model.Execution, error) {
	exe, err := FindCurrentlyActiveByExeId(exeId)
	if err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

func GetCurrentlyActive() (model.Execution, error) {
	return FindCurrentlyActive()
}

// Changes execution status from ACTIVE to PAUSED
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

// Changes the execution status from PAUSED to ACTIVE
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
	if exe.Status == model.EXE_ACTIVE {
		err = fmt.Errorf("execution %s is already ACTIVE", exe.ExeId)
		return model.Execution{}, err
	}
	if exe.Status == model.EXE_TERMINATED {
		err = fmt.Errorf("execution %s is TERMINATED. Cannot be resumed", exe.ExeId)
		return model.Execution{}, err
	}

	exe.Status = model.EXE_ACTIVE
	if err := InsertOne(exe); err != nil {
		return model.Execution{}, err
	}
	return exe, nil
}

// Changes the execution status from ACTIVE or PAUSED to TERMINATED
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

func build_execution(account model.RemoteAccount) model.Execution {
	assets := make([]string, 0, len(account.Balances))
	for _, balance := range account.Balances {
		assets = append(assets, balance.Asset)
	}

	return model.Execution{
		ExeId:     uuid.NewString(),
		Status:    model.EXE_ACTIVE,
		Assets:    assets,
		Timestamp: time.Now().UnixMilli()}
}
