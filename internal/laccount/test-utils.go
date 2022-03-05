package laccount

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
	"go.mongodb.org/mongo-driver/mongo"
)

func restore_laccount_collection(old func() *mongo.Collection) {
	mongodb.GetLocalAccountsCol = old
}

func mock_laccount_collection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetLocalAccountsCol
	mongodb.GetLocalAccountsCol = func() *mongo.Collection {
		return mongoClient.
			Database(mongodb.MONGODB_DATABASE_TEST).
			Collection(mongodb.LACC_COL_NAME)
	}
	return old
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
