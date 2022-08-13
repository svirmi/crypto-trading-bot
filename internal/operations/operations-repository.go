package operations

import (
	"context"

	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insert_many(ops []model.Operation) errors.CtbError {
	collection := mongodb.GetOperationsCol()

	var payload []interface{}
	for i := range ops {
		payload = append(payload, ops[i])
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err := collection.InsertMany(context.TODO(), payload, opts)
	return errors.WrapMongo(err)
}

func insert(op model.Operation) errors.CtbError {
	_, err := mongodb.GetOperationsCol().InsertOne(context.TODO(), op)
	return errors.WrapMongo(err)
}

func find_by_exe_id(exeId string) ([]model.Operation, errors.CtbError) {
	collection := mongodb.GetOperationsCol()

	// Defining query
	options := options.Find().SetSort(bson.D{{"timestamp", 1}})
	filter := bson.D{{"exeId", exeId}}

	// Querying DB
	var results []model.Operation
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, errors.WrapMongo(err)
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, errors.WrapInternal(err)
	}

	// Returning results
	return results, nil
}
