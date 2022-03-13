package fts

import (
	"log"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
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
		BuyThreshold        string
		SellThreshold       string
		StopLossThreshold   string
		MissProfitThreshold string
	}{}
	mapstructure.Decode(c, &tmp)
	s.BuyThreshold = utils.DecimalFromString(tmp.BuyThreshold).Round(2)
	s.SellThreshold = utils.DecimalFromString(tmp.SellThreshold).Round(2)
	s.StopLossThreshold = utils.DecimalFromString(tmp.StopLossThreshold).Round(2)
	s.MissProfitThreshold = utils.DecimalFromString(tmp.MissProfitThreshold).Round(2)

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
