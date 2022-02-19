package fts

import (
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"
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

func get_fts_config(c interface{}) (s strategy_config_fts) {
	// Mapping interface{} to strategy_config_fts
	mapstructure.Decode(c, &s)

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
