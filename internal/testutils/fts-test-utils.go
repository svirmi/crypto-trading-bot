package testutils

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
)

func GetLocalAccountTest_FTS() fts.LocalAccountFTS {
	return fts.LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

		Ignored: map[string]float32{
			"USDT": 1000.01,
			"BUSD": 145.75},

		Assets: map[string]fts.AssetStatusFTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             13.45,
				Usdt:               0,
				LastOperationType:  fts.OP_BUY_FTS,
				LastOperationPrice: 27834.85,
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             0,
				Usdt:               145000.34,
				LastOperationType:  fts.OP_SELL_FTS,
				LastOperationPrice: 3998.45}}}
}

func AssertLocalAccount_FTS(t *testing.T, expected, gotten fts.LocalAccountFTS) {
	if expected.AccountId != gotten.AccountId {
		t.Errorf("AccountId: expected = %s, gotten = %s", expected.AccountId, gotten.AccountId)
	}
	if expected.ExeId != gotten.ExeId {
		t.Errorf("ExeId: expected = %s, gotten = %s", expected.ExeId, gotten.ExeId)
	}
	if expected.StrategyType != gotten.StrategyType {
		t.Errorf("StrategyType: expected = %s, gotten = %s", expected.StrategyType, gotten.StrategyType)
	}
	if expected.Timestamp != gotten.Timestamp {
		t.Errorf("StrategyType: expected = %s, gotten = %s", expected.StrategyType, gotten.StrategyType)
	}
	if !reflect.DeepEqual(expected.Ignored, gotten.Ignored) {
		t.Errorf("Ignored: expected = %v, gotten = %v", expected.Ignored, gotten.Ignored)
	}
	if !reflect.DeepEqual(expected.Assets, gotten.Assets) {
		t.Errorf("Assets: expected = %v, gotten = %v", expected.Assets, gotten.Assets)
	}
}
