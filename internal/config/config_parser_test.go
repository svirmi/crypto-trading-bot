package config

import (
	"path/filepath"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

var (
	fts_test_resource_folder = "test-resources"
)

func TestParseConfig(t *testing.T) {
	test_parse_config(t, false, fts_test_resource_folder)

	gotten := appConfig
	got_sconfig := make(map[string]string)
	mapstructure.Decode(gotten.Strategy.Config, &got_sconfig)
	gotten.Strategy.Config = got_sconfig
	testutils.AssertStructEq(t, get_config(), gotten)
}

func TestParseConfig_Testnet(t *testing.T) {
	test_parse_config(t, true, fts_test_resource_folder)

	gotten := appConfig
	got_sconfig := make(map[string]string)
	mapstructure.Decode(gotten.Strategy.Config, &got_sconfig)
	gotten.Strategy.Config = got_sconfig
	testutils.AssertStructEq(t, get_testnet_config(), gotten)
}

func test_parse_config(t *testing.T, testnet bool, test_resource_folder string) {
	// Restoring interanl status after test execution
	resource_folder_org := resource_folder
	resource_folder = test_resource_folder
	defer func() {
		resource_folder = resource_folder_org
	}()

	// Testing testnet config parsing
	resource_folder = filepath.Join("..", "..", resource_folder)
	config, err := parse_config(testnet)

	// Asserting config
	if err != nil {
		t.Fatalf(err.Error())
	} else {
		appConfig = config
	}
}

func get_config() config {
	return config{
		BinanceApi: binance_api_config{
			ApiKey:     "HTqza54XTX09uBANVQOvMO78N478MhDxLbEiBfSRR8Yc7MBIlXGxG2cwK4Ok3KvI",
			SecretKey:  "vOZJYqQrYjgwSL5EDUxLYTv7Gh8nQvRqX5IefmnySqSAUdVvgOTfTe6HJsO9tvTY",
			UseTestnet: false},
		MongoDb: mongo_db_config{
			Uri:      "mongodb://localhost:27017",
			Database: "ctb"},
		Strategy: strategy_config{
			Type: "TEST_STRATEGY_TYPE",
			Config: map[string]string{
				"prop1": "prop1",
				"prop2": "prop2"}}}
}

func get_testnet_config() config {
	return config{
		BinanceApi: binance_api_config{
			ApiKey:     "fkAHgTpxMVBWXueZfAyYK2NnR4SZdTNPR45mlJivVg4dNEnoWSbODTUQHDkiNjN6",
			SecretKey:  "4yFKwURuMG7onlVoqFeV4Fz3I7ZNcNFMmDTRrlUk45IbbudFEWJtXAQGhqJEJtPg",
			UseTestnet: true},
		MongoDb: mongo_db_config{
			Uri:      "mongodb://localhost:27017",
			Database: "ctb-testnet"},
		Strategy: strategy_config{
			Type: "TEST_STRATEGY_TYPE",
			Config: map[string]string{
				"prop1": "prop1",
				"prop2": "prop2"}}}
}
