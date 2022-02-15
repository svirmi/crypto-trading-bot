package config

import (
	"path/filepath"
	"strconv"
	"testing"
)

func TestParseConfig(t *testing.T) {
	test_parse_config(t, false)
}

func TestParseConfig_Testnet(t *testing.T) {
	test_parse_config(t, true)
}

func test_parse_config(t *testing.T, testnet bool) {
	// Restoring interanl status after test execution
	resource_folder_org := resource_folder
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
		test_config_not_blank(t, testnet)
	}
}

func test_config_not_blank(t *testing.T, testnet bool) {
	if appConfig.IsEmpty() {
		t.Error("empty config detected")
	}

	if GetBinanceApiConfig().IsEmpty() {
		t.Error("empty binance_api_config detected")
	}
	if GetBinanceApiConfig().ApiKey == "" {
		t.Error("empty binance api key detected")
	}
	if GetBinanceApiConfig().SecretKey == "" {
		t.Error("empty binance secret key detected")
	}
	if GetBinanceApiConfig().UseTestnet != testnet {
		expected := strconv.FormatBool(testnet)
		gotten := strconv.FormatBool(!testnet)
		t.Error("binance testnet property set to " + gotten + ", expected " + expected)
	}

	if GetMongoDbConfig().IsEmpty() {
		t.Error("empty mongo_db_config detected")
	}
	if GetMongoDbConfig().Database == "" {
		t.Error("empty mongo db database detected")
	}
	if GetMongoDbConfig().Uri == "" {
		t.Error("empty mongo db uri detected")
	}

	if GetStrategyConfig().IsEmpty() {
		t.Error("empty strategy_config detected")
	}
	if GetStrategyConfig().Config == nil {
		t.Error("uninitialised strategy config detected")
	}
	if GetStrategyConfig().Type == "" {
		t.Error("empty strategy type detected")
	}
}
