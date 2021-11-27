package operations

import (
	"context"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Inserts operations in DB
// Returns an error if computation failed
func InsertMany(ops []model.Operation) error {
	collection := mongodb.GetOperationsCol()

	var payload []interface{}
	for i := range ops {
		payload = append(payload, ops[i])
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err := collection.InsertMany(context.TODO(), payload, opts)
	return err
}

// Inserts an operation in DB
// Returns an error if computation failed
func Insert(op model.Operation) error {
	_, err := mongodb.GetOperationsCol().InsertOne(context.TODO(), op)
	return err
}

// Finds all operation by execution id exeId
// Returns list of operation, if any was found, nil
// oterwise
// Returns an error if computation failed
func FindByExeId(exeId string) ([]model.Operation, error) {
	collection := mongodb.GetOperationsCol()

	// Defining query
	options := options.Find().SetSort(bson.D{{"timestamp", 1}})
	filter := bson.D{{"exeId", exeId}}

	// Querying DB
	var results []model.Operation
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	// Returning results
	return results, nil
}
