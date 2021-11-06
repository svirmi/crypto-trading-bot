package laccount

import (
	"fmt"
	"log"

	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

// Creates a local account based on the remote account, or restores
// account already in DB linked to a currently active execution
// identified by exeId.
// Returns local account object or an empty local account if an
// error was thrown.
// Returns an error if computation failed
// TRUSTS that exeId corresponds to an active execution
func CreateOrRestore(exeId string, raccount model.RemoteAccount, strategyType string) (model.ILocalAccount, error) {
	// Get current local account from DB by execution id
	laccount, err := FindLatest(exeId)
	if err != nil {
		return nil, err
	}

	// Restore existing local account
	if laccount != nil && laccount.GetStrategyType() != strategyType {
		err = fmt.Errorf("strategy type mismatch for exeId %s", exeId)
		return nil, err
	}
	if laccount != nil && laccount.GetStrategyType() == strategyType {
		log.Printf("restoring local account %s", laccount.GetAccountId())
		return laccount, nil
	}

	// Create new local account
	laccount, err = buildLocalAccount(exeId, raccount, strategyType)
	if err != nil {
		return nil, err
	}
	if laccount == nil {
		err = fmt.Errorf("failed to build local account from remote account %v", raccount)
		return nil, err
	}
	if err := Insert(laccount); err != nil {
		return nil, err
	}
	log.Printf("registering local account %s", laccount.GetAccountId())
	return laccount, nil
}

func buildLocalAccount(exeId string, raccount model.RemoteAccount, strategyType string) (model.ILocalAccount, error) {
	if strategyType == model.FIXED_THRESHOLD_STRATEGY {
		return buildLocalAccountFTS(exeId, raccount)
	} else {
		err := fmt.Errorf("unknwon strategy type %s", strategyType)
		return nil, err
	}
}
