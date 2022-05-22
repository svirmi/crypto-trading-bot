package pts

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

type strategy_config_pts struct {
	BuyPercentage        decimal.Decimal
	SellPercentage       decimal.Decimal
	BuyAmountPercentage  decimal.Decimal
	SellAmountPercentage decimal.Decimal
}

func (a strategy_config_pts) is_empty() bool {
	return reflect.DeepEqual(a, strategy_config_pts{})
}

func get_pts_config(strategyConfig config.StrategyConfig) (s strategy_config_pts) {
	if strategyConfig.Type != string(model.PTS_STRATEGY) {
		msg := fmt.Sprintf(logger.XXX_ERR_MISMATCHING_STRATEGY,
			model.PTS_STRATEGY, strategyConfig.Type)
		logrus.WithField("comp", "pts").Panic(msg)
	}

	// Mapping interface{} to strategy_config_pts
	tmp := struct {
		BuyPercentage        string
		SellPercentage       string
		BuyAmountPercentage  string
		SellAmountPercentage string
	}{}
	err := mapstructure.Decode(strategyConfig.Config, &tmp)
	if err != nil {
		logrus.WithField("comp", "pts").Panic(err.Error())
	}

	s.BuyPercentage = utils.DecimalFromString(tmp.BuyPercentage).Round(2)
	s.SellPercentage = utils.DecimalFromString(tmp.SellPercentage).Round(2)
	s.BuyAmountPercentage = utils.DecimalFromString(tmp.BuyAmountPercentage).Round(2)
	s.SellAmountPercentage = utils.DecimalFromString(tmp.SellAmountPercentage).Round(2)

	// Checking config validity
	if s.is_empty() {
		logrus.WithField("comp", "pts").Panicf(logger.XXX_ERR_FAILED_TO_PARSE_CONFIG, strategyConfig.Config)
	}
	if s.BuyPercentage.LessThanOrEqual(decimal.Zero) ||
		s.SellPercentage.LessThanOrEqual(decimal.Zero) ||
		s.BuyAmountPercentage.LessThanOrEqual(decimal.Zero) ||
		s.SellAmountPercentage.LessThanOrEqual(decimal.Zero) {

		logrus.WithField("comp", "pts").Panic(logger.PTS_ERR_NEGATIVE_PERCENTAGES)
	}

	return s
}
