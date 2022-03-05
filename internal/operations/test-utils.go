package operations

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func restore_operation_collection(old func() *mongo.Collection) {
	mongodb.GetOperationsCol = old
}

func mock_operation_collection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetOperationsCol
	mongodb.GetOperationsCol = func() *mongo.Collection {
		return mongoClient.
			Database(mongodb.MONGODB_DATABASE_TEST).
			Collection(mongodb.OP_COL_NAME)
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
		Amount:     decimal.NewFromFloat32(153.78),
		AmountSide: model.BASE_AMOUNT,
		Price:      decimal.NewFromFloat32(133.23),
		Results: model.OpResults{
			ActualPrice: decimal.NewFromFloat32(133.58),
			BaseAmount:  decimal.NewFromFloat32(153.78),
			QuoteAmount: decimal.NewFromFloat32(11224.56),
			Spread:      decimal.NewFromFloat32(12.1),
		},
		Status:    model.FILLED,
		Timestamp: time.Now().UnixMicro()}
}
