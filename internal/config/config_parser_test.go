package config

import (
	"path/filepath"
	"testing"

	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

var (
	fts_test_resource_folder = "test-resources"
)

func TestParseConfig_Fts(t *testing.T) {
	test_parse_config(t, false, fts_test_resource_folder)
	assert_config(t, get_config())
}

func TestParseConfig_Testnet_Fts(t *testing.T) {
	test_parse_config(t, true, fts_test_resource_folder)
	assert_config(t, get_testnet_config())
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

func assert_config(t *testing.T, expected config) {
	if utils.Xor(appConfig.IsEmpty(), expected.IsEmpty()) {
		t.Fatalf("config: expected = %v, gotten = %v", expected, appConfig)
	}

	got_bconf := GetBinanceApiConfig()
	exp_bconf := expected.BinanceApi
	if utils.Xor(got_bconf.IsEmpty(), exp_bconf.IsEmpty()) {
		t.Fatalf("binance config: expected = %v, gotten = %v", exp_bconf, got_bconf)
	}
	if got_bconf.ApiKey != exp_bconf.ApiKey {
		t.Fatalf("binance api key: gotten = %v, expected = %v",
			got_bconf.ApiKey, exp_bconf.ApiKey)
	}
	if got_bconf.SecretKey != exp_bconf.SecretKey {
		t.Fatalf("binance api secret: gotten = %v, expected = %v",
			got_bconf.SecretKey, exp_bconf.SecretKey)
	}
	if got_bconf.UseTestnet != exp_bconf.UseTestnet {
		t.Fatalf("testnet: expetd = %v, gotten = %v",
			got_bconf.UseTestnet, exp_bconf.UseTestnet)
	}

	got_mconf := GetMongoDbConfig()
	exp_mconf := expected.MongoDb
	if utils.Xor(got_mconf.IsEmpty(), exp_mconf.IsEmpty()) {
		t.Fatalf("mongo config: expected = %v, gotten = %v", exp_mconf, got_mconf)
	}
	if got_mconf.Database != exp_mconf.Database {
		t.Fatalf("mongo database: expected = %v, gotten = %v",
			exp_mconf.Database, got_mconf.Database)
	}
	if got_mconf.Uri != exp_mconf.Uri {
		t.Fatalf("mongo uri: expected = %v, gotten = %v",
			exp_mconf.Uri, exp_mconf.Uri)
	}

	got_sconf := GetStrategyConfig()
	exp_sconf := expected.Strategy
	if utils.Xor(got_sconf.IsEmpty(), exp_sconf.IsEmpty()) {
		t.Fatalf("strategy config: expected = %v, gotten = %v", exp_sconf, got_sconf)
	}
	if got_sconf.Type != exp_sconf.Type {
		t.Fatalf("strategy type: expected = %s, gotten = %s",
			exp_sconf.Type, got_sconf.Type)
	}
	if got_sconf.Config == nil {
		t.Fatal("strategy config: expected != nil, gotten = nil")
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
			Config: struct {
				prop1 string
				prop2 string
			}{
				prop1: "prop1",
				prop2: "prop2"}}}
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
			Config: struct {
				prop1 string
				prop2 string
			}{
				prop1: "prop1",
				prop2: "prop2"}}}
}
