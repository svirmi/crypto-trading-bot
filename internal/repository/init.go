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
	mongoClient *mongo.Client
	ctx         context.Context
)

func init() {
	log.Printf("Connecting to mongodb cluster %s", config.AppConfig.MongoDb.Uri)
	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoDb.Uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	mongoClient = client
}

func Disconnect() {
	log.Printf("Disconnecting from mongodb cluster %s", config.AppConfig.MongoDb.Uri)
	mongoClient.Disconnect(ctx)
}

func Ping() {
	err := mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
}
