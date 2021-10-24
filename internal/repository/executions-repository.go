package repository

import (
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOneExecution(exe model.Execution) error {
	_, err := executionsCol.InsertOne(ctx, exe)
	return err
}

func FindAllLatestExecution() ([]model.Execution, error) {
	// Querying executions collection
	sort := bson.D{{"$sort", bson.D{{"timestamp", 1}}}}
	group := bson.D{{"$group", bson.D{
		{"_id", "$exeId"},
		{"status", bson.D{{"$last", "$status"}}},
		{"symbols", bson.D{{"$last", "$symbols"}}},
		{"timestamp", bson.D{{"$last", "$timestamp"}}}}}}
	project := bson.D{{"$project", bson.D{
		{"timestamp", 1},
		{"symbols", 1},
		{"status", 1},
		{"exeId", "$_id"},
		{"_id", 0}}}}

	// Parsing results
	var err error
	var cursor *mongo.Cursor
	var results []model.Execution
	if cursor, err = executionsCol.Aggregate(ctx, mongo.Pipeline{sort, group, project}); err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
