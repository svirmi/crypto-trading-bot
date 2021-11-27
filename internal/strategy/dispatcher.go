package strategy

import (
	"fmt"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
	"go.mongodb.org/mongo-driver/bson"
)

// TODO: create mongo package that will haold one mongo connection only
// TODO: move this function into mongo package using reflection to avoid
// if-elses over the strategy type
func DecodeLocalAccount(data bson.Raw, strategyType string) (model.ILocalAccount, error) {
	var ilaccount model.ILocalAccount = nil
	var err error = nil

	if strategyType == model.FIXED_THRESHOLD_STRATEGY {
		laccount := fts.LocalAccountFTS{}
		err = bson.Unmarshal(data, &laccount)
		ilaccount = laccount
	} else {
		err = fmt.Errorf("unknown strategy type %s", strategyType)
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return ilaccount, nil

}
