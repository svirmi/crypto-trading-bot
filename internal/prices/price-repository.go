package prices

import (
	"context"

	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insert_many(prices []model.SymbolPrice) errors.CtbError {
	collection := mongodb.GetPriceCol()

	var payload []interface{}
	for i := range prices {
		payload = append(payload, prices[i])
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err := collection.InsertMany(context.TODO(), payload, opts)
	return errors.WrapMongo(err)
}

func find(symbols []string, start, end int64) ([]model.SymbolPrice, errors.CtbError) {
	collection := mongodb.GetPriceCol()

	// Defining query
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"symbol", bson.D{{"$in", symbols}}}},
			bson.D{{"timestamp", bson.D{{"$gt", start}}}},
			bson.D{{"timestamp", bson.D{{"$lt", end}}}}}}}

	// Defining query options
	options := options.Find()
	options.SetSort(bson.D{{"timestamp", 1}})

	// Executing query
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, errors.WrapMongo(err)
	}

	// parsing results
	prices := make([]model.SymbolPrice, 0)
	err = cursor.All(context.TODO(), &prices)
	if err != nil {
		return nil, errors.WrapInternal(err)
	}
	return prices, nil
}

func find_by_timestamp(symbols []string, start, end int64) ([]model.SymbolPriceByTimestamp, errors.CtbError) {
	collection := mongodb.GetPriceCol()

	// Defining query
	filter := bson.A{
		bson.D{{"$match", bson.D{
			{"symbol", bson.D{{"$in", symbols}}},
			{"timestamp", bson.D{{"$gt", start}}},
			{"timestamp", bson.D{{"$lt", end}}}}}},
		bson.D{{"$group", bson.D{
			{"_id", "$timestamp"},
			{"symbolPrices", bson.D{{"$push", bson.D{{"symbol", "$symbol"}, {"price", "$price"}}}}}}}},
		bson.D{{"$project", bson.D{
			{"timestamp", "$_id"},
			{"symbolPrices", "$symbolPrices"},
			{"_id", 0}}}},
		bson.D{{"$sort", bson.D{
			{"timestamp", 1}}}}}

	// Executing query
	cursor, err := collection.Aggregate(context.TODO(), filter)
	if err != nil {
		return nil, errors.WrapMongo(err)
	}

	// parsing results
	pricesByTimestamp := make([]model.SymbolPriceByTimestamp, 0)
	err = cursor.All(context.TODO(), &pricesByTimestamp)
	if err != nil {
		return nil, errors.WrapInternal(err)
	}
	return pricesByTimestamp, nil
}
