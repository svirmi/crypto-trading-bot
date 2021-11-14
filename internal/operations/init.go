package operations

import (
	"context"
	"log"

	"github.com/valerioferretti92/trading-bot-demo/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	operationsColName = "operations"
	collection        *mongo.Collection
)

func init() {
	log.Printf("connecting to %s/%s", config.AppConfig.MongoDbConfig.Uri, operationsColName)
	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoDbConfig.Uri)
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Pinging db to test connection
	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Getting mongo collection instance
	collection = mongoClient.
		Database(config.AppConfig.MongoDbConfig.Database).
		Collection(operationsColName)
}

func Close() {
	if collection == nil {
		return
	}
	collection.Database().Client().Disconnect(context.TODO())
}
