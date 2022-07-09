package mongodb

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// Collection names
	_EXE_COL_NAME       = "executions"
	_OP_COL_NAME        = "operations"
	_LACC_COL_NAME      = "laccounts"
	_PRICE_COL_NAME     = "prices"
	_ANALITYCS_COL_NAME = "analytics"

	// Index names
	_ACTIVE_EXECUTION_INDEX = "active-execution-index"
	_LATEST_LACCOUNT_INDEX  = "latest-laccount-index"
)

type mongo_connection struct {
	mongoClient   *mongo.Client
	executionsCol *mongo.Collection
	operationsCol *mongo.Collection
	laccountsCol  *mongo.Collection
	priceCol      *mongo.Collection
	analyticsCol  *mongo.Collection
}

var mongoConnection mongo_connection

func Initialize() error {
	mongoDbConfig := config.GetMongoDbConfig()
	logrus.Infof(logger.MONGO_CONNECTING, mongoDbConfig.Uri)
	clientOptions := options.Client().
		ApplyURI(mongoDbConfig.Uri).
		SetRegistry(GetCustomRegistry())
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
		mongoClient:   mongoClient,
		executionsCol: get_collection_handle(mongoClient, _EXE_COL_NAME),
		operationsCol: get_collection_handle(mongoClient, _OP_COL_NAME),
		laccountsCol:  get_collection_handle(mongoClient, _LACC_COL_NAME),
		priceCol:      get_collection_handle(mongoClient, _PRICE_COL_NAME),
		analyticsCol:  get_collection_handle(mongoClient, _ANALITYCS_COL_NAME)}

	executionIndexes := get_execution_indexes()
	laccountIndexes := get_laccount_indexes()
	err = create_indexes(mongoConnection.executionsCol, executionIndexes)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	err = create_indexes(mongoConnection.laccountsCol, laccountIndexes)
	if err != nil {
		logrus.Error(err.Error())
		return err
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
	return mongoConnection.laccountsCol
}

var GetPriceCol = func() *mongo.Collection {
	return mongoConnection.priceCol
}

var GetAnalyticsCol = func() *mongo.Collection {
	return mongoConnection.analyticsCol
}

func Disconnect() {
	logrus.Info(logger.MONGO_DISCONNECTING)

	if mongoConnection.mongoClient != nil {
		mongoConnection.mongoClient.Disconnect(context.TODO())
	}
}

func get_execution_indexes() []mongo.IndexModel {
	return []mongo.IndexModel{{
		Keys:    bson.D{{"timestamp", -1}},
		Options: options.Index().SetName(_ACTIVE_EXECUTION_INDEX)}}
}

func get_laccount_indexes() []mongo.IndexModel {
	return []mongo.IndexModel{{
		Keys:    bson.D{{"metadata.exeId", -1}, {"metadata.timestamp", -1}},
		Options: options.Index().SetName(_LATEST_LACCOUNT_INDEX)}}
}

func create_indexes(coll *mongo.Collection, indexes []mongo.IndexModel) error {
	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	names, err := coll.Indexes().CreateMany(context.TODO(), indexes, opts)
	if err != nil {
		return err
	}

	logrus.Infof(logger.MONGO_INDEXES_CREATION, names)
	return nil
}

func get_collection_handle(mongoClient *mongo.Client, collection string) *mongo.Collection {
	database := config.GetMongoDbConfig().Database
	logrus.Infof(logger.MONGO_COLLECTION_HANLDE, database, collection)
	return mongoClient.
		Database(database).
		Collection(collection)
}
