package laccount

import (
	"context"

	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Inserts a new local account object
// Returns an error, if computation failed
func Insert(laccout model.LocalAccount) error {
	_, err := collection.InsertOne(context.TODO(), laccout)
	return err
}

// Finds latest version of a local wallet bound to a given
// execution id.
// Returns a local wallet or en empty wallet if nothing was
// found or an error was thrown.
// Returns an error if computation failed
func FindLatest(exeId string) (model.LocalAccount, error) {
	// Defining query
	filter := bson.D{{"exeId", exeId}}
	options := options.Find()
	options.SetSort(bson.D{{"timestamp", -1}})
	options.SetLimit(1)

	// Querying DB
	var results []model.LocalAccount
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return model.LocalAccount{}, err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return model.LocalAccount{}, err
	}

	// Returning results
	if len(results) == 0 {
		return model.LocalAccount{}, nil
	}
	return results[0], nil
}

// Returns all local wallet versions bound to an execution id.
// returns the list of local wallet versions, a nil slice if
// nothing was found or an error was thorwn
// Returns an error if computation failed
func FindAll(exeId string) ([]model.LocalAccount, error) {
	// Defining query
	filter := bson.D{{"exeId", exeId}}
	options := options.Find()
	options.SetSort(bson.D{{"timestamp", -1}})

	// Querying DB
	var results []model.LocalAccount
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, nil
	}

	// Returning query results
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
