package executions

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
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
