package executions

import (
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
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

func get_execution_init() model.ExecutionInit {
	balances := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("5.0")},
		{Asset: "ETH", Amount: utils.DecimalFromString("10.45")}}
	raccount := model.RemoteAccount{
		MakerCommission:  0,
		TakerCommission:  1,
		BuyerCommission:  2,
		SellerCommission: 1,
		Balances:         balances}
	return model.ExecutionInit{
		Raccount:     raccount,
		StrategyType: model.DTS_STRATEGY,
		Props: map[string]string{
			"prop1": "value1",
			"prop2": "value2"}}

}

func get_execution() model.Execution {
	return model.Execution{
		ExeId:        uuid.NewString(),
		Status:       model.EXE_ACTIVE,
		Assets:       []string{"BTC", "ETH"},
		StrategyType: model.DTS_STRATEGY,
		Props: map[string]string{
			"prop1": "value1",
			"prop2": "value2"},
		Timestamp: time.Now().UnixMicro()}
}
