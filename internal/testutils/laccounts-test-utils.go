package testutils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

var _LACC_COLLECTION_TEST string = "laccounts"

func RestoreLaccountCollection(old func() *mongo.Collection) {
	mongodb.GetLocalAccountsCol = old
}

func MockLaccountCollection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetLocalAccountsCol
	mongodb.GetLocalAccountsCol = func() *mongo.Collection {
		return mongoClient.
			Database(_MONGODB_DATABASE_TEST).
			Collection(_LACC_COLLECTION_TEST)
	}
	return old
}

func GetLocalAccountInitTest(strategyType model.StrategyType) model.LocalAccountInit {
	return model.LocalAccountInit{
		ExeId: uuid.NewString(),
		RAccount: model.RemoteAccount{
			MakerCommission:  0,
			TakerCommission:  0,
			BuyerCommission:  0,
			SellerCommission: 0,
			Balances: []model.RemoteBalance{
				{Asset: "BTC", Amount: 11.34},
				{Asset: "ETH", Amount: 29.12}}},
		TradableAssetsPrice: map[string]model.AssetPrice{
			"ETH": {Asset: "ETH", Price: 4500.45},
			"BTC": {Asset: "BTC", Price: 39560.45},
			"DOT": {Asset: "DOT", Price: 49.45}},
		StrategyType: strategyType}
}

func AssertInitLocalAccount(t *testing.T, init model.LocalAccountInit, gotten model.ILocalAccount) {
	if gotten.GetExeId() != init.ExeId {
		t.Errorf("ExeId: expected = %s, gotten = %s", init.ExeId, gotten.GetExeId())
	}
	if gotten.GetAccountId() == "" {
		t.Error("AccountId: expected != nil, gotten = nil")
	}
	if gotten.GetStrategyType() != init.StrategyType {
		t.Errorf("StrategyType: expected = %s, gotten = %s", init.StrategyType, gotten.GetStrategyType())
	}
	if gotten.GetTimestamp() == 0 {
		t.Error("Timestamp: expected != 0, gotten = 0")
	}
}
