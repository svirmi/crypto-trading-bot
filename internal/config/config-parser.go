package config

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	yaml "gopkg.in/yaml.v2"
)

type ExchangeConfig interface{}

type MongoDbConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
}

func (m MongoDbConfig) IsEmpty() bool {
	return reflect.DeepEqual(m, MongoDbConfig{})
}

type StrategyConfig struct {
	Type   string            `yaml:"type"`
	Config map[string]string `yaml:"config"`
}

func (s StrategyConfig) IsEmpty() bool {
	return reflect.DeepEqual(s, StrategyConfig{})
}

type Config struct {
	Exchange ExchangeConfig `yaml:"exchange"`
	MongoDb  MongoDbConfig  `yaml:"mongoDb"`
	Strategy StrategyConfig `yaml:"strategy"`
}

func (c Config) IsEmpty() bool {
	return reflect.DeepEqual(c, Config{})
}

var (
	appConfig              Config
	simulation_config_file = "config-simulation.yaml"
	testnet_config_file    = "config-testnet.yaml"
	mainnet_config_file    = "config.yaml"
	resource_folder        = "resources"
)

func Initialize(env model.Env) error {
	// Parsing config
	config, err := parse_config(env)
	if err != nil {
		return err
	}

	appConfig = config
	return nil
}

var GetExchangeConfig = func() ExchangeConfig {
	return appConfig.Exchange
}

var GetMongoDbConfig = func() MongoDbConfig {
	return appConfig.MongoDb
}

var GetStrategyConfig = func() StrategyConfig {
	return appConfig.Strategy
}

func parse_config(env model.Env) (config Config, err error) {
	configPath := get_config_filepath(env)
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	logrus.Infof(logger.CONFIG_PARSING, configPath)
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func get_config_filepath(env model.Env) string {
	if env == model.SIMULATION {
		return filepath.Join(resource_folder, simulation_config_file)
	} else if env == model.TESTNET {
		return filepath.Join(resource_folder, testnet_config_file)
	} else {
		return filepath.Join(resource_folder, mainnet_config_file)
	}
}
