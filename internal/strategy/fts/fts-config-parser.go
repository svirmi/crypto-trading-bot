package fts

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
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

func get_fts_config(strategyConfig config.StrategyConfig) (s strategy_config_fts, err error) {
	if strategyConfig.Type != string(model.FIXED_THRESHOLD_STRATEGY) {
		msg := fmt.Sprintf(logger.FTS_ERR_MISMATCHING_STRATEGY,
			model.FIXED_THRESHOLD_STRATEGY, strategyConfig.Type)
		logrus.Error(msg)
		return strategy_config_fts{}, model.NewCtbError(msg, false)
	}

	// Mapping interface{} to strategy_config_fts
	tmp := struct {
		BuyThreshold        string
		SellThreshold       string
		StopLossThreshold   string
		MissProfitThreshold string
	}{}
	mapstructure.Decode(strategyConfig.Config, &tmp)
	s.BuyThreshold = utils.DecimalFromString(tmp.BuyThreshold).Round(2)
	s.SellThreshold = utils.DecimalFromString(tmp.SellThreshold).Round(2)
	s.StopLossThreshold = utils.DecimalFromString(tmp.StopLossThreshold).Round(2)
	s.MissProfitThreshold = utils.DecimalFromString(tmp.MissProfitThreshold).Round(2)

	// Checking config validity
	if s.is_empty() {
		msg := fmt.Sprintf(logger.FTS_ERR_FAILED_TO_PARSE_CONFIG, strategyConfig.Config)
		logrus.Error(msg)
		return strategy_config_fts{}, model.NewCtbError(msg, true)
	}
	if s.BuyThreshold.LessThanOrEqual(decimal.Zero) ||
		s.SellThreshold.LessThanOrEqual(decimal.Zero) ||
		s.MissProfitThreshold.LessThanOrEqual(decimal.Zero) ||
		s.StopLossThreshold.LessThanOrEqual(decimal.Zero) {
		msg := fmt.Sprintf(logger.FTS_ERR_NEGATIVE_THRESHOLDS)
		logrus.Error(msg)
		return strategy_config_fts{}, model.NewCtbError(msg, true)
	}

	return s, nil
}
