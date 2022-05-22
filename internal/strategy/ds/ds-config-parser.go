package ds

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

type strategy_config_ds struct {
	BuyThreshold        decimal.Decimal
	SellThreshold       decimal.Decimal
	StopLossThreshold   decimal.Decimal
	MissProfitThreshold decimal.Decimal
}

func (a strategy_config_ds) is_empty() bool {
	return reflect.DeepEqual(a, strategy_config_ds{})
}

func get_ds_config(strategyConfig config.StrategyConfig) (s strategy_config_ds) {
	if strategyConfig.Type != string(model.DEMO_STRATEGY) {
		msg := fmt.Sprintf(logger.DS_ERR_MISMATCHING_STRATEGY,
			model.DEMO_STRATEGY, strategyConfig.Type)
		logrus.WithField("comp", "ds").Panic(msg)
	}

	// Mapping interface{} to strategy_config_ds
	tmp := struct {
		BuyThreshold        string
		SellThreshold       string
		StopLossThreshold   string
		MissProfitThreshold string
	}{}
	err := mapstructure.Decode(strategyConfig.Config, &tmp)
	if err != nil {
		logrus.Panic(err.Error())
	}

	s.BuyThreshold = utils.DecimalFromString(tmp.BuyThreshold).Round(2)
	s.SellThreshold = utils.DecimalFromString(tmp.SellThreshold).Round(2)
	s.StopLossThreshold = utils.DecimalFromString(tmp.StopLossThreshold).Round(2)
	s.MissProfitThreshold = utils.DecimalFromString(tmp.MissProfitThreshold).Round(2)

	// Checking config validity
	if s.is_empty() {
		logrus.WithField("comp", "ds").Panicf(logger.DS_ERR_FAILED_TO_PARSE_CONFIG, strategyConfig.Config)
	}
	if s.BuyThreshold.LessThanOrEqual(decimal.Zero) ||
		s.SellThreshold.LessThanOrEqual(decimal.Zero) ||
		s.MissProfitThreshold.LessThanOrEqual(decimal.Zero) ||
		s.StopLossThreshold.LessThanOrEqual(decimal.Zero) {

		logrus.WithField("comp", "ds").Panic(logger.DS_ERR_NEGATIVE_THRESHOLDS)
	}

	return s
}
