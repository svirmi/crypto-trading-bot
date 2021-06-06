package config

import (
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var (
	config_testnet_path = "resources/config-testnet.yaml"
	config_path         = "resources/config.yaml"
)

var (
	AppConfig Config
)

type Config struct {
	BinanceApi struct {
		ApiKey     string `yaml:"apiKey"`
		SecretKey  string `yaml:"secretKey"`
		UseTestnet bool   `yaml:"useTestnet"`
	} `yaml:"binanceApi"`
}

func ParseConfig(testnet bool) (Config, error) {
	if (Config{}) != AppConfig {
		return AppConfig, nil
	}

	var configPath string
	if testnet {
		configPath = config_testnet_path
	} else {
		configPath = config_path
	}

	f, err := os.Open(configPath)
	if err != nil {
		return (Config{}), fmt.Errorf("could not open %s", configPath)
	}
	defer f.Close()

	log.Printf("Parsing config file %s", configPath)
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		return (Config{}), fmt.Errorf("could not parse %s", configPath)
	}
	return AppConfig, nil
}
