package config

import (
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	yaml "gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type ExchangeConfig interface{}

type MongoDbConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
}

func (m MongoDbConfig) IsEmpty() bool {
	return reflect.DeepEqual(m, MongoDbConfig{})
}

type Config struct {
	Server   ServerConfig
	Exchange ExchangeConfig `yaml:"exchange"`
	MongoDb  MongoDbConfig  `yaml:"mongoDb"`
}

func (c Config) IsEmpty() bool {
	return reflect.DeepEqual(c, Config{})
}

var (
	appConfig Config
)

func Initialize(configFilepath string) errors.CtbError {
	config, err := parse_config(configFilepath)
	appConfig = config
	return err
}

var GetServerConfig = func() ServerConfig {
	return appConfig.Server
}

var GetExchangeConfig = func() ExchangeConfig {
	return appConfig.Exchange
}

var GetMongoDbConfig = func() MongoDbConfig {
	return appConfig.MongoDb
}

func parse_config(configFilepath string) (Config, errors.CtbError) {
	f, err := os.Open(configFilepath)
	if err != nil {
		return Config{}, errors.WrapBadRequest(err)
	}
	defer f.Close()

	logrus.Infof(logger.CONFIG_PARSING, configFilepath)
	decoder := yaml.NewDecoder(f)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, errors.WrapBadRequest(err)
	}
	return config, nil
}
