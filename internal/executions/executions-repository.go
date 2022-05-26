package executions

import (
	"context"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insert_one(exe model.Execution) error {
	// Inserting new execution object
	_, err := mongodb.GetExecutionsCol().InsertOne(context.TODO(), exe)
	return err
}

func find_latest_by_exeId(exeId string) (model.Execution, error) {
	collection := mongodb.GetExecutionsCol()

	var result model.Execution
	opts := options.FindOne().SetSort(bson.D{{"timestamp", -1}})
	err := collection.FindOne(context.TODO(), bson.D{{"exeId", exeId}}, opts).Decode(&result)
	return result, err
}

func find_currently_active() (model.Execution, error) {
	collection := mongodb.GetExecutionsCol()

	var result model.Execution
	opts := options.FindOne().SetSort(bson.D{{"timestamp", -1}})
	err := collection.FindOne(context.TODO(), bson.D{{}}, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return model.Execution{}, nil
	}
	if err != nil {
		return model.Execution{}, err
	}

	if result.Status != model.EXE_ACTIVE {
		return model.Execution{}, nil
	}
	return result, nil
}
