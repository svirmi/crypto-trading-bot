package operations

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func TestMain(m *testing.M) {
	logger.Initialize(true, logrus.TraceLevel)
	code := m.Run()
	os.Exit(code)
}

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
		Amount:     utils.DecimalFromString("153.78"),
		AmountSide: model.BASE_AMOUNT,
		Price:      utils.DecimalFromString("133.23"),
		Results: model.OpResults{
			ActualPrice: utils.DecimalFromString("133.58"),
			BaseDiff:    utils.DecimalFromString("153.78"),
			QuoteDiff:   utils.DecimalFromString("11224.56"),
			Spread:      utils.DecimalFromString("12.1"),
		},
		Status:    model.FILLED,
		Timestamp: time.Now().UnixMicro()}
}
