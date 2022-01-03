package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

type binance_api_config struct {
	ApiKey     string `yaml:"apiKey"`
	SecretKey  string `yaml:"secretKey"`
	UseTestnet bool   `yaml:"useTestnet"`
}

func (b binance_api_config) IsEmpty() bool {
	return reflect.DeepEqual(b, binance_api_config{})
}

type mongo_db_config struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
}

func (m mongo_db_config) IsEmpty() bool {
	return reflect.DeepEqual(m, mongo_db_config{})
}

type strategy_config struct {
	Type   string      `yaml:"type"`
	Config interface{} `yaml:"config"`
}

func (s strategy_config) IsEmpty() bool {
	return reflect.DeepEqual(s, strategy_config{})
}

type config struct {
	BinanceApi binance_api_config `yaml:"binanceApi"`
	MongoDb    mongo_db_config    `yaml:"mongoDb"`
	Strategy   strategy_config    `yaml:"strategy"`
}

func (c config) IsEmpty() bool {
	return reflect.DeepEqual(c, config{})
}

var (
	appConfig           config
	testnet_config_file = "config-testnet.yaml"
	mainnet_config_file = "config.yaml"
	resource_folder     = "resources"
)

func ParseConfig() {
	// Parsing command line
	testnet := flag.Bool("testnet", false, "if present, application runs on testnet")
	flag.Parse()

	// Parsing config
	config, err := parse_config(*testnet)
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

func parse_config(testnet bool) (config config, err error) {
	configPath := get_config_filepath(testnet)
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

func get_config_filepath(testnet bool) string {
	if testnet {
		return filepath.Join(resource_folder, testnet_config_file)
	} else {
		return filepath.Join(resource_folder, mainnet_config_file)
	}
}
