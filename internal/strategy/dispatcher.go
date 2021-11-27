package strategy

import (
	"fmt"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
	"go.mongodb.org/mongo-driver/bson"
)

func InitLocalAccount(creationRequest model.LocalAccountInit) (model.ILocalAccount, error) {
	if creationRequest.StrategyType == model.FIXED_THRESHOLD_STRATEGY {
		return fts.InitLocalAccountFTS(creationRequest)
	} else {
		err := fmt.Errorf("unknwon strategy type %s", creationRequest.StrategyType)
		return nil, err
	}
}

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
