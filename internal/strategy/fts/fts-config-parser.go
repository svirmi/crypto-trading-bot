package fts

import (
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type strategy_config_fts struct {
	BuyThreshold        decimal.Decimal
	SellThreshold       decimal.Decimal
	StopLossThreshold   decimal.Decimal
	MissProfitThreshold decimal.Decimal
}

func (a strategy_config_fts) is_empty() bool {
	return reflect.DeepEqual(a, strategy_config_fts{})
}

func get_fts_config(c interface{}) (s strategy_config_fts) {
	// Mapping interface{} to strategy_config_fts
	tmp := struct {
		BuyThreshold        float32
		SellThreshold       float32
		StopLossThreshold   float32
		MissProfitThreshold float32
	}{}
	mapstructure.Decode(c, &tmp)
	s.BuyThreshold = decimal.NewFromFloat32(tmp.BuyThreshold)
	s.SellThreshold = decimal.NewFromFloat32(tmp.SellThreshold)
	s.StopLossThreshold = decimal.NewFromFloat32(tmp.StopLossThreshold)
	s.MissProfitThreshold = decimal.NewFromFloat32(tmp.MissProfitThreshold)

	// Checking config validity
	if s.is_empty() {
		log.Fatalf("failed to parse fts config")
	}
	if s.BuyThreshold.LessThanOrEqual(decimal.Zero) ||
		s.SellThreshold.LessThanOrEqual(decimal.Zero) ||
		s.MissProfitThreshold.LessThanOrEqual(decimal.Zero) ||
		s.StopLossThreshold.LessThanOrEqual(decimal.Zero) {
		log.Fatalf("fts thresholds must be strictly positive")
	}

	log.Printf("fts strategy config: %+v", s)
	return s
}
