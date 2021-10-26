package repository

import (
	"context"
	"log"

	"github.com/valerioferretti92/trading-bot-demo/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx context.Context

	mongoClient       *mongo.Client
	operationsColName string = "operations"
	operationsCol     *mongo.Collection
	executionsColName string = "executions"
	executionsCol     *mongo.Collection
)

func init() {
	// Connecting to db
	log.Printf("connecting to %s", config.AppConfig.MongoDbConfig.Uri)
	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoDbConfig.Uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}
	mongoClient = client

	// Pinging db to test connection
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	// Bilding collection handles
	operationsCol = client.
		Database(config.AppConfig.MongoDbConfig.Database).
		Collection(operationsColName)
	// Bilding collection handles
	executionsCol = client.
		Database(config.AppConfig.MongoDbConfig.Database).
		Collection(executionsColName)
}

func Disconnect() {
	if mongoClient == nil {
		return
	}

	log.Printf("disconnecting from %s", config.AppConfig.MongoDbConfig.Uri)
	mongoClient.Disconnect(ctx)
}
