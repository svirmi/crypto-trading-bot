package repository

import (
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertManyOperations(ops []model.Operation) error {
	var payload []interface{}
	for i := range ops {
		payload = append(payload, ops[i])
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err := operationsCol.InsertMany(ctx, payload, opts)
	return err
}

func FindLatestOperations(exeId string, symbols []string) ([]model.Operation, error) {
	match := bson.D{{"$match", bson.D{{"exeId", exeId}, {"symbol", bson.D{{"$in", symbols}}}}}}
	sort := bson.D{{"$sort", bson.D{{"timestamp", 1}}}}
	group := bson.D{{"$group", bson.D{
		{"_id", "$symbol"},
		{"opId", bson.D{{"$last", "$opId"}}},
		{"exeId", bson.D{{"$last", "$exeId"}}},
		{"type", bson.D{{"$last", "$type"}}},
		{"price", bson.D{{"$last", "$price"}}},
		{"forcastedQty", bson.D{{"$last", "$forcastedQty"}}},
		{"qty", bson.D{{"$last", "$qty"}}},
		{"timestamp", bson.D{{"$last", "$timestamp"}}}}}}
	project := bson.D{{"$project", bson.D{
		{"opId", 1},
		{"exeId", 1},
		{"type", 1},
		{"price", 1},
		{"forcastedQty", 1},
		{"qty", 1},
		{"timestamp", 1},
		{"symbol", "$_id"},
		{"_id", 0}}}}

	// Parsing results
	var results []model.Operation
	cursor, err := operationsCol.Aggregate(ctx, mongo.Pipeline{match, sort, group, project})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
