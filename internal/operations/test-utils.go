package operations

import (
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/mongo"
)

var _OP_COLLECTION_TEST string = "operations"

func restore_operation_collection(old func() *mongo.Collection) {
	mongodb.GetOperationsCol = old
}

func mock_operation_collection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetOperationsCol
	mongodb.GetOperationsCol = func() *mongo.Collection {
		return mongoClient.
			Database(testutils.MONGODB_DATABASE_TEST).
			Collection(_OP_COLLECTION_TEST)
	}
	return old
}

func get_operation_test() model.Operation {
	return model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      uuid.NewString(),
		Type:       model.AUTO,
		Base:       "BTC",
		Quote:      "USDT",
		Side:       model.BUY,
		Amount:     153.78,
		AmountSide: model.BASE_AMOUNT,
		Price:      133.23,
		Results: model.OpResults{
			ActualPrice: 133.58,
			BaseAmount:  153.78,
			QuoteAmount: 11224.56,
			Spread:      12.1,
		},
		Status:    model.FILLED,
		Timestamp: time.Now().UnixMicro()}
}
