package executions

import (
	"context"
	"fmt"

	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Inserts a new execution object.
// Checks that exe.ExeId does not reference an execution that
// was already terminated
// Returns an error if computation failed or checks did not succeed
func InsertOne(exe model.Execution) error {
	// Defining query
	query := bson.D{{"exeId", exe.ExeId}}
	options := options.Find().
		SetSort(bson.D{{"timestamp", -1}}).
		SetLimit(1)

	// Querying DB
	var results []model.Execution
	cursor, err := collection.Find(context.TODO(), query, options)
	if err != nil {
		return err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return err
	}

	//Checking results
	if len(results) == 1 && results[0].Status == model.EXE_TERMINATED {
		return fmt.Errorf("execution %s has already been terminated", exe.ExeId)
	}

	// Inserting new execution object
	_, err = collection.InsertOne(context.TODO(), exe)
	return err
}

// Finds currently execution object by execution id.
// Returns the execution object, if found, an empty execution
// object if nothing was found or an error was thrown.
// Returns an error if computation failed
func FindCurrentlyActiveByExeId(exeId string) (model.Execution, error) {
	// Defining query stages
	filter1 := bson.D{{"exeId", exeId}}
	sort := bson.D{{"$sort", bson.D{{"timestamp", -1}}}}
	limit := bson.D{{"$limit", 1}}
	filter2 := bson.D{{"$match", bson.D{{"status", bson.D{{"$ne", "TERMINATED"}}}}}}

	// Querying DB
	var results []model.Execution
	pipeline := mongo.Pipeline{filter1, sort, limit, filter2}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return model.Execution{}, err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return model.Execution{}, err
	}

	// Returning results
	if len(results) == 0 {
		return model.Execution{}, nil
	}
	return results[0], nil
}

// Finds currently active execution object.
// Returns the execution object, if found, an empty execution
// object if nothing was found or an error was thrown.
// Returns an error if computation failed
func FindCurrentlyActive() (model.Execution, error) {
	// Defining query stages
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

	// Querying DB
	var results []model.Execution
	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{sort, group, project, filter})
	if err != nil {
		return model.Execution{}, err
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return model.Execution{}, err
	}

	// Returning results
	if len(results) > 1 {
		err = fmt.Errorf("more then one active executions found")
		return model.Execution{}, err
	}
	if len(results) == 0 {
		return model.Execution{}, nil
	}
	return results[0], nil
}
