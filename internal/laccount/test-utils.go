package laccount

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
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

func get_laccount_init_test(strategyType model.StrategyType) model.LocalAccountInit {
	return model.LocalAccountInit{
		ExeId: uuid.NewString(),
		RAccount: model.RemoteAccount{
			MakerCommission:  0,
			TakerCommission:  0,
			BuyerCommission:  0,
			SellerCommission: 0,
			Balances: []model.RemoteBalance{
				{Asset: "BTC", Amount: decimal.NewFromFloat32(11.34)},
				{Asset: "ETH", Amount: decimal.NewFromFloat32(29.12)},
				{Asset: "DOT", Amount: decimal.NewFromFloat32(13.67)},
				{Asset: "USDT", Amount: decimal.NewFromFloat32(155.67)},
				{Asset: "BUSD", Amount: decimal.NewFromFloat32(1232.45)}}},
		TradableAssetsPrice: map[string]model.AssetPrice{
			"BTC": {Asset: "BTC", Price: decimal.NewFromFloat32(39560.45)},
			"ETH": {Asset: "ETH", Price: decimal.NewFromFloat32(4500.45)},
			"DOT": {Asset: "DOT", Price: decimal.NewFromFloat32(49.45)}},
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
			"USDT": decimal.NewFromFloat32(155.67),
			"BUSD": decimal.NewFromFloat32(1232.45)},

		Assets: map[string]fts.AssetStatusFTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             decimal.NewFromFloat32(11.34),
				Usdt:               decimal.Zero,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: decimal.NewFromFloat32(39560.45),
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             decimal.NewFromFloat32(29.12),
				Usdt:               decimal.Zero,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: decimal.NewFromFloat32(4500.45)},
			"DOT": {
				Asset:              "DOT",
				Amount:             decimal.NewFromFloat32(13.67),
				Usdt:               decimal.Zero,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: decimal.NewFromFloat32(49.45)}}}
}
