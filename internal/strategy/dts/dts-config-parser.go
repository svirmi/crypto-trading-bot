package dts

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

type strategy_config_dts struct {
	BuyThreshold        decimal.Decimal
	SellThreshold       decimal.Decimal
	StopLossThreshold   decimal.Decimal
	MissProfitThreshold decimal.Decimal
}

func (l LocalAccountDTS) ValidateConfig(props map[string]string) error {
	_, err := parse_config(props)
	return err
}

func parse_config(props map[string]string) (s strategy_config_dts, err error) {
	bt, found := props[_BUY_THRESHOLD]
	if !found {
		err = fmt.Errorf(logger.XXX_ERR_MISSING_PROP_KEY, _BUY_THRESHOLD)
		logrus.WithField("comp", "dts").Error(err.Error())
		return strategy_config_dts{}, err
	}
	st, found := props[_SELL_THRESHOLD]
	if !found {
		err = fmt.Errorf(logger.XXX_ERR_MISSING_PROP_KEY, _SELL_THRESHOLD)
		logrus.WithField("comp", "dts").Error(err.Error())
		return strategy_config_dts{}, err
	}
	mpt, found := props[_MISS_PROFIT_THRESHOLD]
	if !found {
		err = fmt.Errorf(logger.XXX_ERR_MISSING_PROP_KEY, _MISS_PROFIT_THRESHOLD)
		logrus.WithField("comp", "dts").Error(err.Error())
		return strategy_config_dts{}, err
	}
	slt, found := props[_STOP_LOSS_THRESHOLD]
	if !found {
		err = fmt.Errorf(logger.XXX_ERR_MISSING_PROP_KEY, _STOP_LOSS_THRESHOLD)
		logrus.WithField("comp", "dts").Error(err.Error())
		return strategy_config_dts{}, err
	}

	s.BuyThreshold = utils.DecimalFromString(bt).Round(2)
	s.SellThreshold = utils.DecimalFromString(st).Round(2)
	s.StopLossThreshold = utils.DecimalFromString(slt).Round(2)
	s.MissProfitThreshold = utils.DecimalFromString(mpt).Round(2)

	// Checking config validity
	if s.BuyThreshold.LessThanOrEqual(decimal.Zero) ||
		s.SellThreshold.LessThanOrEqual(decimal.Zero) ||
		s.MissProfitThreshold.LessThanOrEqual(decimal.Zero) ||
		s.StopLossThreshold.LessThanOrEqual(decimal.Zero) {

		err = fmt.Errorf(logger.DTS_ERR_NEGATIVE_THRESHOLDS)
		logrus.WithField("comp", "dts").Error(err.Error())
		return strategy_config_dts{}, err
	}

	return s, nil
}
