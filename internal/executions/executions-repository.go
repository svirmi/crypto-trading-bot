package executions

import (
	"context"

	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insert_one(exe model.Execution) errors.CtbError {
	// Inserting new execution object
	_, err := mongodb.GetExecutionsCol().InsertOne(context.TODO(), exe)
	return errors.WrapMongo(err)
}

func find_latest() (model.Execution, errors.CtbError) {
	collection := mongodb.GetExecutionsCol()

	opts := options.FindOne().SetSort(bson.D{{"timestamp", -1}})
	sr := collection.FindOne(context.TODO(), bson.D{{}}, opts)
	if sr.Err() == mongo.ErrNoDocuments {
		return model.Execution{}, nil
	}

	var result model.Execution
	err := sr.Decode(&result)
	if err != nil {
		return model.Execution{}, errors.WrapInternal(err)
	}

	return result, nil
}

func find_latest_by_exeId(exeId string) (model.Execution, errors.CtbError) {
	collection := mongodb.GetExecutionsCol()

	var result model.Execution
	opts := options.FindOne().SetSort(bson.D{{"timestamp", -1}})
	sr := collection.FindOne(context.TODO(), bson.D{{"exeId", exeId}}, opts)
	if sr.Err() == mongo.ErrNoDocuments {
		return model.Execution{}, nil
	}

	err := sr.Decode(&result)
	if err != nil {
		return model.Execution{}, errors.WrapInternal(err)
	}
	return result, nil
}

func find_by_exeId(exeId string) ([]model.Execution, errors.CtbError) {
	collection := mongodb.GetExecutionsCol()

	opts := options.Find().SetSort(bson.D{{"timestamp", 1}})
	cursor, err := collection.Find(context.TODO(), bson.D{{"exeId", exeId}}, opts)
	if err != nil {
		return nil, errors.WrapMongo(err)
	}

	var results []model.Execution
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, errors.WrapInternal(err)
	}
	return results, nil
}
