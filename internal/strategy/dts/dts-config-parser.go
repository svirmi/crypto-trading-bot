package dts

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

type strategy_config_dts struct {
	BuyThreshold        decimal.Decimal
	SellThreshold       decimal.Decimal
	StopLossThreshold   decimal.Decimal
	MissProfitThreshold decimal.Decimal
}

func (a strategy_config_dts) is_empty() bool {
	return reflect.DeepEqual(a, strategy_config_dts{})
}

func get_dts_config(strategyConfig config.StrategyConfig) (s strategy_config_dts) {
	if strategyConfig.Type != string(model.DTS_STRATEGY) {
		msg := fmt.Sprintf(logger.XXX_ERR_MISMATCHING_STRATEGY,
			model.DTS_STRATEGY, strategyConfig.Type)
		logrus.WithField("comp", "dts").Panic(msg)
	}

	// Mapping interface{} to strategy_config_dts
	tmp := struct {
		BuyThreshold        string
		SellThreshold       string
		StopLossThreshold   string
		MissProfitThreshold string
	}{}
	err := mapstructure.Decode(strategyConfig.Config, &tmp)
	if err != nil {
		logrus.WithField("comp", "dts").Panic(err.Error())
	}

	s.BuyThreshold = utils.DecimalFromString(tmp.BuyThreshold).Round(2)
	s.SellThreshold = utils.DecimalFromString(tmp.SellThreshold).Round(2)
	s.StopLossThreshold = utils.DecimalFromString(tmp.StopLossThreshold).Round(2)
	s.MissProfitThreshold = utils.DecimalFromString(tmp.MissProfitThreshold).Round(2)

	// Checking config validity
	if s.is_empty() {
		logrus.WithField("comp", "dts").Panicf(logger.XXX_ERR_FAILED_TO_PARSE_CONFIG, strategyConfig.Config)
	}
	if s.BuyThreshold.LessThanOrEqual(decimal.Zero) ||
		s.SellThreshold.LessThanOrEqual(decimal.Zero) ||
		s.MissProfitThreshold.LessThanOrEqual(decimal.Zero) ||
		s.StopLossThreshold.LessThanOrEqual(decimal.Zero) {

		logrus.WithField("comp", "dts").Panic(logger.DTS_ERR_NEGATIVE_THRESHOLDS)
	}

	return s
}
