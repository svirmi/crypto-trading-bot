package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

var (
	fts_test_resource_folder = "test-resources"
)

func TestMain(m *testing.M) {
	logger.Initialize(true, logrus.TraceLevel)
	code := m.Run()
	os.Exit(code)
}

func TestInitialize(t *testing.T) {
	test_parse_config(t, model.MAINNET, fts_test_resource_folder)

	got := appConfig
	got_sconfig := make(map[string]string)
	mapstructure.Decode(got.Strategy.Config, &got_sconfig)
	got.Strategy.Config = got_sconfig

	testutils.AssertEq(t, get_config(), got, "config")
}

func TestInitialize_Testnet(t *testing.T) {
	test_parse_config(t, model.TESTNET, fts_test_resource_folder)

	got := appConfig
	got_sconfig := make(map[string]string)
	mapstructure.Decode(got.Strategy.Config, &got_sconfig)
	got.Strategy.Config = got_sconfig

	testutils.AssertEq(t, get_testnet_config(), got, "config")
}

func test_parse_config(t *testing.T, env model.Env, test_resource_folder string) {
	// Restoring interanl status after test execution
	resource_folder_org := resource_folder
	resource_folder = test_resource_folder
	defer func() {
		resource_folder = resource_folder_org
	}()

	// Testing testnet config parsing
	resource_folder = filepath.Join("..", "..", resource_folder)
	config, _ := parse_config(env)
	appConfig = config
}

func get_config() Config {
	return Config{
		BinanceApi: BinanceApiConfig{
			ApiKey:     "HTqza54XTX09uBANVQOvMO78N478MhDxLbEiBfSRR8Yc7MBIlXGxG2cwK4Ok3KvI",
			SecretKey:  "vOZJYqQrYjgwSL5EDUxLYTv7Gh8nQvRqX5IefmnySqSAUdVvgOTfTe6HJsO9tvTY",
			UseTestnet: false},
		MongoDb: MongoDbConfig{
			Uri:      "mongodb://localhost:27017",
			Database: "ctb"},
		Strategy: StrategyConfig{
			Type: "TEST_STRATEGY_TYPE",
			Config: map[string]string{
				"prop1": "prop1",
				"prop2": "prop2"}}}
}

func get_testnet_config() Config {
	return Config{
		BinanceApi: BinanceApiConfig{
			ApiKey:     "fkAHgTpxMVBWXueZfAyYK2NnR4SZdTNPR45mlJivVg4dNEnoWSbODTUQHDkiNjN6",
			SecretKey:  "4yFKwURuMG7onlVoqFeV4Fz3I7ZNcNFMmDTRrlUk45IbbudFEWJtXAQGhqJEJtPg",
			UseTestnet: true},
		MongoDb: MongoDbConfig{
			Uri:      "mongodb://localhost:27017",
			Database: "ctb-testnet"},
		Strategy: StrategyConfig{
			Type: "TEST_STRATEGY_TYPE",
			Config: map[string]string{
				"prop1": "prop1",
				"prop2": "prop2"}}}
}
