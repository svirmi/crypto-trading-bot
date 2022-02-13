package testutils

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	_MONGODB_URI_TEST      string = "mongodb://localhost:27017"
	_MONGODB_DATABASE_TEST string = "ctb-unit-tests"
)

func GetMongoClientTest() *mongo.Client {
	clientOptions := options.Client().ApplyURI(_MONGODB_URI_TEST)
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Pinging db to test connection
	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Returning mongo client
	return mongoClient
}
