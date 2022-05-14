package laccount

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func TestMain(m *testing.M) {
	logger.Initialize(false, logrus.TraceLevel)
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

func get_laccount_init_test(strategyType model.StrategyType) model.LocalAccountInit {
	return model.LocalAccountInit{
		ExeId: uuid.NewString(),
		RAccount: model.RemoteAccount{
			MakerCommission:  0,
			TakerCommission:  0,
			BuyerCommission:  0,
			SellerCommission: 0,
			Balances: []model.RemoteBalance{
				{Asset: "BTC", Amount: utils.DecimalFromString("11.34")},
				{Asset: "ETH", Amount: utils.DecimalFromString("29.12")},
				{Asset: "DOT", Amount: utils.DecimalFromString("13.67")},
				{Asset: "LUNA", Amount: utils.DecimalFromString("90.67")},
				{Asset: "USDT", Amount: utils.DecimalFromString("155.67")},
				{Asset: "BUSD", Amount: utils.DecimalFromString("1232.45")}}},
		TradableAssetsPrice: map[string]model.AssetPrice{
			"BTC": {Asset: "BTC", Price: utils.DecimalFromString("39560.45")},
			"ETH": {Asset: "ETH", Price: utils.DecimalFromString("4500.45")},
			"DOT": {Asset: "DOT", Price: utils.DecimalFromString("49.45")}},
		StrategyType: strategyType}
}

func get_laccount_test_FTS() fts.LocalAccountFTS {
	return fts.LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

		Ignored: map[string]decimal.Decimal{
			"USDT": utils.DecimalFromString("155.67"),
			"LUNA": utils.DecimalFromString("90.67"),
			"BUSD": utils.DecimalFromString("1232.45")},

		Assets: map[string]fts.AssetStatusFTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             utils.DecimalFromString("11.34"),
				Usdt:               decimal.Zero,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: utils.DecimalFromString("39560.45"),
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             utils.DecimalFromString("29.12"),
				Usdt:               decimal.Zero,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: utils.DecimalFromString("4500.45")},
			"DOT": {
				Asset:              "DOT",
				Amount:             utils.DecimalFromString("13.67"),
				Usdt:               decimal.Zero,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: utils.DecimalFromString("49.45")}}}
}
