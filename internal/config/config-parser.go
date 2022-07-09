package config

import (
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
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

type Config struct {
	Exchange ExchangeConfig `yaml:"exchange"`
	MongoDb  MongoDbConfig  `yaml:"mongoDb"`
}

func (c Config) IsEmpty() bool {
	return reflect.DeepEqual(c, Config{})
}

var (
	appConfig Config
)

func Initialize(configFilepath string) error {
	config, err := parse_config(configFilepath)
	appConfig = config
	return err
}

var GetExchangeConfig = func() ExchangeConfig {
	return appConfig.Exchange
}

var GetMongoDbConfig = func() MongoDbConfig {
	return appConfig.MongoDb
}

func parse_config(configFilepath string) (config Config, err error) {
	f, err := os.Open(configFilepath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	logrus.Infof(logger.CONFIG_PARSING, configFilepath)
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
