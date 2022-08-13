package laccount

import (
	"context"

	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insert(laccout model.ILocalAccount) errors.CtbError {
	_, err := mongodb.GetLocalAccountsCol().InsertOne(context.TODO(), laccout)
	return errors.WrapMongo(err)
}

func find_latest_by_exeId(exeId string) (model.ILocalAccount, errors.CtbError) {
	collection := mongodb.GetLocalAccountsCol()

	// Defining query
	filter := bson.D{{"metadata.exeId", exeId}}
	options := options.FindOne()
	options.SetSort(bson.D{{"metadata.timestamp", -1}})

	// Querying DB
	result := collection.FindOne(context.TODO(), filter, options)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}

	// Decoding result
	raw, _ := result.DecodeBytes()
	return strategy.DecodeLaccount(raw, mongodb.GetCustomRegistry())
}

func find_by_exeId(exeId string) ([]model.ILocalAccount, errors.CtbError) {
	collection := mongodb.GetLocalAccountsCol()

	// Defining query
	filter := bson.D{{"metadata.exeId", exeId}}
	options := options.Find()
	options.SetSort(bson.D{{"metadata.timestamp", 1}})

	// Querying DB
	results, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, errors.WrapMongo(err)
	}

	// Decoding results
	laccounts := make([]model.ILocalAccount, 0)
	for results.Next(context.TODO()) {
		raw := results.Current
		laccount, err := strategy.DecodeLaccount(raw, mongodb.GetCustomRegistry())
		if err != nil {
			return nil, err
		}
		laccounts = append(laccounts, laccount)
	}
	return laccounts, nil
}
