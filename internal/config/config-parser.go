package config

import (
	"flag"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type binance_api_config struct {
	ApiKey     string `yaml:"apiKey"`
	SecretKey  string `yaml:"secretKey"`
	UseTestnet bool   `yaml:"useTestnet"`
}

type mongo_db_config struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type strategy_config struct {
	Type   string      `yaml:"type"`
	Config interface{} `yaml:"config"`
}

type config struct {
	BinanceApi binance_api_config `yaml:"binanceApi"`
	MongoDb    mongo_db_config    `yaml:"mongoDb"`
	Strategy   strategy_config    `yaml:"strategy"`
}

var (
	appConfig           config
	config_testnet_path = "resources/config-testnet.yaml"
	config_path         = "resources/config.yaml"
)

func init() {
	// Parsing command line
	testnet := flag.Bool("testnet", false, "if present, application runs on testnet")
	flag.Parse()

	// Parsing config
	config, err := parse_config(testnet)
	if err != nil {
		log.Fatalf(err.Error())
	}
	appConfig = config
}

func GetBinanceApiConfig() binance_api_config {
	return appConfig.BinanceApi
}

func GetMongoDbConfig() mongo_db_config {
	return appConfig.MongoDb
}

func GetStrategyConfig() strategy_config {
	return appConfig.Strategy
}

func parse_config(testnet *bool) (config config, err error) {
	var configPath string
	if *testnet {
		configPath = config_testnet_path
	} else {
		configPath = config_path
	}

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("could not open %s", configPath)
	}
	defer f.Close()

	log.Printf("parsing config file %s", configPath)
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("could not parse %s", configPath)
	}
	return config, nil
}
