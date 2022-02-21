package laccount

import (
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/mongo"
)

var _LACC_COLLECTION_TEST string = "laccounts"

func restore_laccount_collection(old func() *mongo.Collection) {
	mongodb.GetLocalAccountsCol = old
}

func mock_laccount_collection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetLocalAccountsCol
	mongodb.GetLocalAccountsCol = func() *mongo.Collection {
		return mongoClient.
			Database(testutils.MONGODB_DATABASE_TEST).
			Collection(_LACC_COLLECTION_TEST)
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
				{Asset: "BTC", Amount: 11.34},
				{Asset: "ETH", Amount: 29.12},
				{Asset: "DOT", Amount: 13.67},
				{Asset: "USDT", Amount: 155.67},
				{Asset: "BUSD", Amount: 1232.45}}},
		TradableAssetsPrice: map[string]model.AssetPrice{
			"BTC": {Asset: "BTC", Price: 39560.45},
			"ETH": {Asset: "ETH", Price: 4500.45},
			"DOT": {Asset: "DOT", Price: 49.45}},
		StrategyType: strategyType}
}

func get_laccount_test_FTS() fts.LocalAccountFTS {
	return fts.LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

		Ignored: map[string]float32{
			"USDT": 155.67,
			"BUSD": 1232.45},

		Assets: map[string]fts.AssetStatusFTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             11.34,
				Usdt:               0,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: 39560.45,
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             29.12,
				Usdt:               0,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: 4500.45},
			"DOT": {
				Asset:              "DOT",
				Amount:             13.67,
				Usdt:               0,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: 49.45}}}
}
