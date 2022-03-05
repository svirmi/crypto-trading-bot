package laccount

import (
	"context"
	"fmt"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Inserts a new local account object
// Returns an error, if computation failed
func insert(laccout model.ILocalAccount) error {
	_, err := mongodb.GetLocalAccountsCol().InsertOne(context.TODO(), laccout)
	return err
}

// Finds latest version of a local wallet bound to a given
// execution id.
// Returns a local wallet or en empty wallet if nothing was
// found or an error was thrown.
// Returns an error if computation failed
func find_latest_by_exeId(exeId string) (model.ILocalAccount, error) {
	collection := mongodb.GetLocalAccountsCol()

	// Defining query
	filter := bson.D{{"metadata.exeId", exeId}}
	options := options.FindOne()
	options.SetSort(bson.D{{"metadata.timestamp", -1}})

	// Querying DB
	result := collection.FindOne(context.TODO(), filter, options)
	return decode(result)
}

func decode(sr *mongo.SingleResult) (model.ILocalAccount, error) {
	payload := struct {
		model.LocalAccountMetadata `bson:"metadata"`
	}{}

	err := sr.Decode(&payload)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	strategyType := model.StrategyType(payload.GetStrategyType())
	if strategyType == model.FIXED_THRESHOLD_STRATEGY {
		laccount_fts := fts.LocalAccountFTS{}
		err := sr.Decode(&laccount_fts)
		return laccount_fts, err
	} else {
		return nil, fmt.Errorf("unrecognized strategy type %s", strategyType)
	}
}
