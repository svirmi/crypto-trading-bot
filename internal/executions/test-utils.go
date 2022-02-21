package executions

import (
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/mongo"
)

var _EXE_COLLECTION_TEST string = "executions"

func restore_execution_collection(old func() *mongo.Collection) {
	mongodb.GetExecutionsCol = old
}

func mock_execution_collection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetExecutionsCol
	mongodb.GetExecutionsCol = func() *mongo.Collection {
		return mongoClient.
			Database(testutils.MONGODB_DATABASE_TEST).
			Collection(_EXE_COLLECTION_TEST)
	}
	return old
}
