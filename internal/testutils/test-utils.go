package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MONGODB_URI_TEST      string = "mongodb://localhost:27017"
	MONGODB_DATABASE_TEST string = "ctb-unit-tests"
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

func AssertStructEq(t *testing.T, exp, got interface{}) {
	bexp, err := json.MarshalIndent(exp, "", "  ")
	if err != nil {
		t.Fatalf("failed to enocode payload: %v", exp)
	}

	bgot, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("failed to enocode payload: %v", got)
	}

	res := bytes.Compare(bexp, bgot)
	if res != 0 {
		t.Errorf("exp = %s", string(bexp[:]))
		t.Errorf("got = %s", string(bgot[:]))
		t.Fatal("exp and got structs are not equivalent")
	}
}
