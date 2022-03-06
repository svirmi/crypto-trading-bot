package operations

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

func mock_mongo_config() func() config.MongoDbConfig {
	old := config.GetMongoDbConfig
	config.GetMongoDbConfig = func() config.MongoDbConfig {
		return config.MongoDbConfig{
			Uri:      testutils.MONGODB_URI_TEST,
			Database: testutils.MONGODB_DATABASE_TEST,
		}
	}
	return old
}

func restore_mongo_config(old func() config.MongoDbConfig) {
	config.GetMongoDbConfig = old
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
