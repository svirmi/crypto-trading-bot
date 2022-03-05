package executions

import (
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func restore_execution_collection(old func() *mongo.Collection) {
	mongodb.GetExecutionsCol = old
}

func mock_execution_collection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetExecutionsCol
	mongodb.GetExecutionsCol = func() *mongo.Collection {
		return mongoClient.
			Database(mongodb.MONGODB_DATABASE_TEST).
			Collection(mongodb.EXE_COL_NAME)
	}
	return old
}
