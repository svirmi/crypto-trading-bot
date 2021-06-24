package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var (
	config_testnet_path = "resources/config-testnet.yaml"
	config_path         = "resources/config.yaml"

	AppConfig Config
)

type Config struct {
	BinanceApi struct {
		ApiKey     string `yaml:"apiKey"`
		SecretKey  string `yaml:"secretKey"`
		UseTestnet bool   `yaml:"useTestnet"`
	} `yaml:"binanceApi"`
	MongoDbConfig struct {
		Uri      string `yaml:"uri"`
		Database string `yaml:"database"`
	} `yaml:"mongoDbConfig"`
}

func init() {
	// Parsing command line
	testnet := flag.Bool("testnet", false, "if present, application runs on testnet")
	flag.Parse()

	// Parsing config
	_, err := parseConfig(*testnet)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
}

func parseConfig(testnet bool) (Config, error) {
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

	log.Printf("parsing config file %s", configPath)
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		return (Config{}), fmt.Errorf("could not parse %s", configPath)
	}
	return AppConfig, nil
}
