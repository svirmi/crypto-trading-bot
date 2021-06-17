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

	mongoClient            *mongo.Client
	miniMarketsStatColName string = "mini-markets-stat"
	miniMarketsStatCol     *mongo.Collection
)

func init() {
	// Connecting to db
	log.Printf("Connecting to mongodb cluster %s", config.AppConfig.MongoDbConfig.Uri)
	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoDbConfig.Uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	mongoClient = client

	// Pinging db to test connection
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// Bilding collection handles
	miniMarketsStatCol = client.
		Database(config.AppConfig.MongoDbConfig.Database).
		Collection(miniMarketsStatColName)
}

func Disconnect() {
	log.Printf("Disconnecting from mongodb cluster %s", config.AppConfig.MongoDbConfig.Uri)
	mongoClient.Disconnect(ctx)
}
