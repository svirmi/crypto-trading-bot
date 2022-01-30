package utilstest

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MONGODB_URI_TEST      string = "mongodb://localhost:27017"
	MONGODB_DATABASE_TEST string = "ctb-unit-tests"
	EXE_COLLECTION_TEST   string = "executions"
)

func GetMongoClientTest() *mongo.Client {
	clientOptions := options.Client().ApplyURI(MONGODB_URI_TEST)
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
