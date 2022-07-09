package config

import (
	"path/filepath"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

var (
	test_resource_folder            = "test-resources"
	test_simulation_config_filename = "config-simulation.yaml"
	test_testnet_config_filename    = "config-testnet.yaml"
	test_mainnet_config_filename    = "config.yaml"
)

func TestInitialize_Mainnet(t *testing.T) {
	logger.Initialize(false, true, true)
	test_parse_config(test_resource_folder, test_mainnet_config_filename)

	got := appConfig
	got_econfig := make(map[string]string)
	mapstructure.Decode(got.Exchange, &got_econfig)
	got.Exchange = got_econfig

	testutils.AssertEq(t, get_mainnet_config(), got, "config")
}

func TestInitialize_Testnet(t *testing.T) {
	logger.Initialize(false, true, true)
	test_parse_config(test_resource_folder, test_testnet_config_filename)

	got := appConfig
	got_econfig := make(map[string]string)
	mapstructure.Decode(got.Exchange, &got_econfig)
	got.Exchange = got_econfig

	testutils.AssertEq(t, get_testnet_config(), got, "config")
}

func TestInitialize_Simulation(t *testing.T) {
	logger.Initialize(false, true, true)
	test_parse_config(test_resource_folder, test_simulation_config_filename)

	got := appConfig
	got_econfig := make(map[string]string)
	mapstructure.Decode(got.Exchange, &got_econfig)
	got.Exchange = got_econfig

	testutils.AssertEq(t, get_simulation_config(), got, "config")
}

func test_parse_config(test_resource_folder, test_config_filename string) {
	// Testing testnet config parsing
	test_config_filepath := filepath.Join("..", "..", test_resource_folder, test_config_filename)
	config, _ := parse_config(test_config_filepath)
	appConfig = config
}

func get_mainnet_config() Config {
	return Config{
		Exchange: map[string]string{
			"propA": "propA",
			"propB": "propB"},
		MongoDb: MongoDbConfig{
			Uri:      "mongodb://localhost:27017",
			Database: "ctb"}}
}

func get_testnet_config() Config {
	return Config{
		Exchange: map[string]string{
			"propA": "propA",
			"propB": "propB"},
		MongoDb: MongoDbConfig{
			Uri:      "mongodb://localhost:27017",
			Database: "ctb-testnet"}}
}

func get_simulation_config() Config {
	return Config{
		Exchange: map[string]string{
			"propA": "propA",
			"propB": "propB"},
		MongoDb: MongoDbConfig{
			Uri:      "mongodb://localhost:27017",
			Database: "ctb-simulation"}}
}
