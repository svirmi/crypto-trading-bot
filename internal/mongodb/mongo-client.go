package mongodb

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// Collection names
	_EXE_COL_NAME  = "executions"
	_OP_COL_NAME   = "operations"
	_LACC_COL_NAME = "laccounts"
)

type mongo_connection struct {
	mongoClient      *mongo.Client
	executionsCol    *mongo.Collection
	operationsCol    *mongo.Collection
	localAccountsCol *mongo.Collection
}

var mongoConnection mongo_connection

func Initialize() error {
	mongoDbConfig := config.GetMongoDbConfig()
	logrus.Infof(logger.MONGO_CONNECTING, mongoDbConfig.Uri)
	clientOptions := options.Client().
		ApplyURI(mongoDbConfig.Uri).
		SetRegistry(build_custom_registry())
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Pinging db to test connection
	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return err
	}

	// Getting collection handles
	mongoConnection = mongo_connection{
		mongoClient:      mongoClient,
		executionsCol:    get_collection_handle(mongoClient, _EXE_COL_NAME),
		operationsCol:    get_collection_handle(mongoClient, _OP_COL_NAME),
		localAccountsCol: get_collection_handle(mongoClient, _LACC_COL_NAME),
	}
	return nil
}

var GetExecutionsCol = func() *mongo.Collection {
	return mongoConnection.executionsCol
}

var GetOperationsCol = func() *mongo.Collection {
	return mongoConnection.operationsCol
}

var GetLocalAccountsCol = func() *mongo.Collection {
	return mongoConnection.localAccountsCol
}

func Disconnect() {
	logrus.Info(logger.MONGO_DISCONNECTING)

	if mongoConnection.mongoClient != nil {
		mongoConnection.mongoClient.Disconnect(context.TODO())
	}
}

func get_collection_handle(mongoClient *mongo.Client, collection string) *mongo.Collection {
	database := config.GetMongoDbConfig().Database
	logrus.Infof(logger.MONGO_COLLECTION_HANLDE, database, collection)
	return mongoClient.
		Database(database).
		Collection(collection)
}
