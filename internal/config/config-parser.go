package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

type BinanceApiConfig struct {
	ApiKey     string `yaml:"apiKey"`
	SecretKey  string `yaml:"secretKey"`
	UseTestnet bool   `yaml:"useTestnet"`
}

func (b BinanceApiConfig) IsEmpty() bool {
	return reflect.DeepEqual(b, BinanceApiConfig{})
}

type MongoDbConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
}

func (m MongoDbConfig) IsEmpty() bool {
	return reflect.DeepEqual(m, MongoDbConfig{})
}

type StrategyConfig struct {
	Type   string      `yaml:"type"`
	Config interface{} `yaml:"config"`
}

func (s StrategyConfig) IsEmpty() bool {
	return reflect.DeepEqual(s, StrategyConfig{})
}

type Config struct {
	BinanceApi BinanceApiConfig `yaml:"binanceApi"`
	MongoDb    MongoDbConfig    `yaml:"mongoDb"`
	Strategy   StrategyConfig   `yaml:"strategy"`
}

func (c Config) IsEmpty() bool {
	return reflect.DeepEqual(c, Config{})
}

var (
	appConfig           Config
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

var GetBinanceApiConfig = func() BinanceApiConfig {
	return appConfig.BinanceApi
}

var GetMongoDbConfig = func() MongoDbConfig {
	return appConfig.MongoDb
}

var GetStrategyConfig = func() StrategyConfig {
	return appConfig.Strategy
}

func parse_config(testnet bool) (config Config, err error) {
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
