package laccount

import (
	"context"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Inserts a new local account object
// Returns an error, if computation failed
func Insert(laccout model.ILocalAccount) error {
	_, err := collection.InsertOne(context.TODO(), laccout)
	return err
}

// Finds latest version of a local wallet bound to a given
// execution id.
// Returns a local wallet or en empty wallet if nothing was
// found or an error was thrown.
// Returns an error if computation failed
func FindLatest(exeId string) (model.ILocalAccount, error) {
	// Defining query
	filter := bson.D{{"metadata.exeId", exeId}}
	options := options.Find()
	options.SetSort(bson.D{{"metadata.timestamp", -1}})
	options.SetLimit(1)

	// Querying DB
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	results, err := decode_many(cursor)
	if err != nil {
		return nil, err
	}

	// Returning results
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

// Returns all local wallet versions bound to an execution id.
// returns the list of local wallet versions, a nil slice if
// nothing was found or an error was thorwn
// Returns an error if computation failed
func FindAll(exeId string) ([]model.ILocalAccount, error) {
	// Defining query
	filter := bson.D{{"metadata.exeId", exeId}}
	options := options.Find()
	options.SetSort(bson.D{{"metadata.timestamp", -1}})

	// Querying DB
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, nil
	}

	// Returning query results
	results, err := decode_many(cursor)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func decode_many(cursor *mongo.Cursor) ([]model.ILocalAccount, error) {
	laccounts := make([]model.ILocalAccount, 0)
	for cursor.Next(context.TODO()) {
		raw := cursor.Current
		laccount, err := decode_one(raw)
		if err != nil {
			return nil, err
		}
		laccounts = append(laccounts, laccount)
	}
	return laccounts, nil
}

func decode_one(raw bson.Raw) (model.ILocalAccount, error) {
	payload := struct {
		model.LocalAccountMetadata `bson:"metadata"`
	}{}

	err := bson.Unmarshal(raw, &payload)
	if err != nil {
		return nil, err
	}
	return strategy.DecodeLocalAccount(raw, payload.StrategyType)
}
