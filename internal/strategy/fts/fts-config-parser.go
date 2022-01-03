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

func get_fts_config() (s strategy_config_fts) {
	strategyConfig := config.GetStrategyConfig()
	strategyType := model.StrategyType(strategyConfig.Type)

	// Check strategy type
	if strategyType != model.FIXED_THRESHOLD_STRATEGY {
		log.Fatalf("wrong startegy type %s", strategyConfig.Type)
	}

	// Mapping interface{} to strategy_config_fts
	mapstructure.Decode(strategyConfig.Config, &s)

	// Checking config validity
	if s.is_empty() {
		log.Fatalf("failed to parse fts config")
	}
	if s.BuyThreshold <= 0 ||
		s.SellThreshold <= 0 ||
		s.MissProfitThreshold <= 0 ||
		s.StopLossThreshold <= 0 {
		log.Fatalf("fts thresholds must be strictly positive")
	}

	log.Printf("fts strategy config: %+v", s)
	return s
}
