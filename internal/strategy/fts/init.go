package fts

import (
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

type strategy_config_fts struct {
	BuyThreshold        float32
	SellThreshold       float32
	StopLossThreshold   float32
	MissProfitThreshold float32
}

func (a strategy_config_fts) is_empty() bool {
	return reflect.DeepEqual(a, strategy_config_fts{})
}

var strategy_config strategy_config_fts

func init() {
	strategyConfig := config.GetStrategyConfig()

	// Check strategy type
	if strategyConfig.Type != model.FIXED_THRESHOLD_STRATEGY {
		log.Fatalf("wrong startegy type %s", strategyConfig.Type)
	}

	// Mapping interface{} to strategy_config_fts
	mapstructure.Decode(strategyConfig.Config, &strategy_config)

	// Checking config validity
	if strategy_config.is_empty() {
		log.Fatalf("failed to parse fts config")
	}
	if strategy_config.BuyThreshold <= 0 ||
		strategy_config.SellThreshold <= 0 ||
		strategy_config.MissProfitThreshold <= 0 ||
		strategy_config.StopLossThreshold <= 0 {
		log.Fatalf("fts thresholds must be strictly positive")
	}

	log.Printf("fts strategy config: %+v", strategy_config)
}
