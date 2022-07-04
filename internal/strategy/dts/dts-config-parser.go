package dts

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

const (
	_BUY_THRESHOLD         = "buyThreshold"
	_SELL_THRESHOLD        = "sellThreshold"
	_STOP_LOSS_THRESHOLD   = "stopLossThreshold"
	_MISS_PROFIT_THRESHOLD = "missProfitThreshold"
)

type strategy_config_dts struct {
	BuyThreshold        decimal.Decimal
	SellThreshold       decimal.Decimal
	StopLossThreshold   decimal.Decimal
	MissProfitThreshold decimal.Decimal
}

func parse_config(props map[string]string) (s strategy_config_dts) {
	bt, found := props[_BUY_THRESHOLD]
	if !found {
		logrus.WithField("comp", "dts").Panicf(logger.XXX_ERR_MISSING_PROP_KEY, _BUY_THRESHOLD)
	}
	st, found := props[_SELL_THRESHOLD]
	if !found {
		logrus.WithField("comp", "dts").Panicf(logger.XXX_ERR_MISSING_PROP_KEY, _SELL_THRESHOLD)
	}
	mpt, found := props[_MISS_PROFIT_THRESHOLD]
	if !found {
		logrus.WithField("comp", "dts").Panicf(logger.XXX_ERR_MISSING_PROP_KEY, _MISS_PROFIT_THRESHOLD)
	}
	slt, found := props[_STOP_LOSS_THRESHOLD]
	if !found {
		logrus.WithField("comp", "dts").Panicf(logger.XXX_ERR_MISSING_PROP_KEY, _STOP_LOSS_THRESHOLD)
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

		logrus.WithField("comp", "dts").Panic(logger.DTS_ERR_NEGATIVE_THRESHOLDS)
	}

	return s
}
