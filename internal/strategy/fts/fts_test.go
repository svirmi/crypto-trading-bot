package fts

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

func TestInitialize(t *testing.T) {
	gotten, err := LocalAccountFTS{}.Initialize(get_laccount_init_test())
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()

	testutils.AssertStructEq(t, expected, gotten)
}

func get_laccount_init_test() model.LocalAccountInit {
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
		StrategyType: model.FIXED_THRESHOLD_STRATEGY}
}

func get_laccount_test() LocalAccountFTS {
	return LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

		Ignored: map[string]float32{
			"USDT": 155.67,
			"BUSD": 1232.45},

		Assets: map[string]AssetStatusFTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             11.34,
				Usdt:               0,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: 39560.45,
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             29.12,
				Usdt:               0,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: 4500.45},
			"DOT": {
				Asset:              "DOT",
				Amount:             13.67,
				Usdt:               0,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: 49.45}}}
}
