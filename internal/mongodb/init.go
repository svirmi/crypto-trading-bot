package mongodb

import (
	"context"
	"log"

	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	executionsColName    = "executions"
	operationsColName    = "operations"
	localAccountsColName = "laccounts"
)

type mongo_connection struct {
	mongoClient      *mongo.Client
	executionsCol    *mongo.Collection
	operationsCol    *mongo.Collection
	localAccountsCol *mongo.Collection
}

var mongoConnection mongo_connection

func init() {
	log.Printf("connecting to mongo instance: %s", config.AppConfig.MongoDbConfig.Uri)
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

	// Getting collection handles
	mongoConnection = mongo_connection{
		mongoClient:      mongoClient,
		executionsCol:    get_collection_handle(mongoClient, executionsColName),
		operationsCol:    get_collection_handle(mongoClient, operationsColName),
		localAccountsCol: get_collection_handle(mongoClient, localAccountsColName),
	}
}

func GetExecutionsCol() *mongo.Collection {
	return mongoConnection.executionsCol
}

func GetOperationsCol() *mongo.Collection {
	return mongoConnection.operationsCol
}

func GetLocalAccountsCol() *mongo.Collection {
	return mongoConnection.localAccountsCol
}

func Disconnect() {
	log.Printf("disconnecting from mongodb")

	if mongoConnection.mongoClient != nil {
		mongoConnection.mongoClient.Disconnect(context.TODO())
	}
}

func get_collection_handle(mongoClient *mongo.Client, collection string) *mongo.Collection {
	log.Printf("getting handle to %s collection", collection)
	return mongoClient.
		Database(config.AppConfig.MongoDbConfig.Database).
		Collection(collection)
}