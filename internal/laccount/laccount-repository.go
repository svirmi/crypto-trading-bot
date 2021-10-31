package laccount

import (
	"context"

	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Insert(laccout model.LocalAccount) error {
	_, err := collection.InsertOne(context.TODO(), laccout)
	return err
}

func FindLatest(exeId string) (model.LocalAccount, error) {
	filter := bson.D{{"exeId", exeId}}

	options := options.Find()
	options.SetSort(bson.D{{"timestamp", -1}})
	options.SetLimit(1)

	var results []model.LocalAccount
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return model.LocalAccount{}, nil
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return model.LocalAccount{}, err
	}
	if len(results) == 0 {
		return model.LocalAccount{}, nil
	}
	return results[0], nil
}

func FindAll(exeId string) ([]model.LocalAccount, error) {
	filter := bson.D{{"exeId", exeId}}

	options := options.Find()
	options.SetSort(bson.D{{"timestamp", -1}})

	var results []model.LocalAccount
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, nil
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
