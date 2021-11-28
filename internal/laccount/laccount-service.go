package laccount

import (
	"fmt"
	"log"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
)

// Creates a local account based on the remote account, or restores
// account already in DB linked to a currently active execution
// identified by exeId.
// Returns local account object or an empty local account if an
// error was thrown.
// Returns an error if computation failed
// TRUSTS that exeId corresponds to an active execution
func CreateOrRestore(creationRequest model.LocalAccountInit) (model.ILocalAccount, error) {
	// Get current local account from DB by execution id
	laccount, err := FindLatest(creationRequest.ExeId)
	if err != nil {
		return nil, err
	}

	// Restore existing local account
	if laccount != nil && laccount.GetStrategyType() != creationRequest.StrategyType {
		err = fmt.Errorf("strategy type mismatch for exeId %s", creationRequest.ExeId)
		return nil, err
	}
	if laccount != nil && laccount.GetStrategyType() == creationRequest.StrategyType {
		log.Printf("restoring local account %s", laccount.GetAccountId())
		return laccount, nil
	}

	// Initialise new local account
	laccount, err = initialise_local_account(creationRequest)
	if err != nil {
		return nil, err
	}
	if laccount == nil {
		err = fmt.Errorf("failed to build local account from remote account")
		return nil, err
	}

	// Inseting in mongo db and returning value
	if err := Insert(laccount); err != nil {
		return nil, err
	}
	log.Printf("registering local account %s", laccount.GetAccountId())
	return laccount, nil
}

func initialise_local_account(creationRequest model.LocalAccountInit) (model.ILocalAccount, error) {
	var laccount model.ILocalAccount = nil

	if creationRequest.StrategyType == model.FIXED_THRESHOLD_STRATEGY {
		laccount = fts.LocalAccountFTS{}
	} else {
		err := fmt.Errorf("unknwon strategy type %s", creationRequest.StrategyType)
		return nil, err
	}
	laccount, err := laccount.Initialize(creationRequest)
	if err != nil {
		return nil, err
	}
	return laccount, nil
}
