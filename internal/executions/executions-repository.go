package executions

import (
	"context"

	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOne(exe model.Execution) error {
	_, err := collection.InsertOne(context.TODO(), exe)
	return err
}

func FindActive() ([]model.Execution, error) {
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
	filter := bson.D{{"$match", bson.D{{"status", bson.D{{"$ne", "TERMINATED"}}}}}}

	// Parsing results
	var results []model.Execution
	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{sort, group, project, filter})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
