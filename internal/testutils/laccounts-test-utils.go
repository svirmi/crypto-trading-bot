package testutils

import (
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

var _LACC_COLLECTION_TEST string = "laccounts"

func RestoreLaccountCollection(old func() *mongo.Collection) {
	mongodb.GetLocalAccountsCol = old
}

func MockLaccountCollection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetLocalAccountsCol
	mongodb.GetLocalAccountsCol = func() *mongo.Collection {
		return mongoClient.
			Database(_MONGODB_DATABASE_TEST).
			Collection(_LACC_COLLECTION_TEST)
	}
	return old
}
