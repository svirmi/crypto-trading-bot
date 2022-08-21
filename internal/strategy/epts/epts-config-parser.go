package epts

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

type strategy_config_epts struct {
	BuyPercentage            decimal.Decimal
	SellPercentage           decimal.Decimal
	InitBuyAmountPercentage  decimal.Decimal
	InitSellAmountPercentage decimal.Decimal
	ExponentialBase          decimal.Decimal
}

func (l LocalAccountEPTS) ValidateConfig(props map[string]string) errors.CtbError {
	_, err := parse_config(props)
	return err
}

func parse_config(props map[string]string) (s strategy_config_epts, err errors.CtbError) {
	bp, found := props[_BUY_PERCENTAGE]
	if !found {
		err = errors.BadRequest(logger.XXX_ERR_MISSING_PROP_KEY, _BUY_PERCENTAGE)
		logrus.WithField("comp", "epts").Error(err.Error())
		return strategy_config_epts{}, err
	}
	sp, found := props[_SELL_PERCENTAGE]
	if !found {
		err = errors.BadRequest(logger.XXX_ERR_MISSING_PROP_KEY, _SELL_PERCENTAGE)
		logrus.WithField("comp", "epts").Error(err.Error())
		return strategy_config_epts{}, err
	}
	bap, found := props[_INIT_BUY_AMOUNT_PERCENTAGE]
	if !found {
		err = errors.BadRequest(logger.XXX_ERR_MISSING_PROP_KEY, _INIT_BUY_AMOUNT_PERCENTAGE)
		logrus.WithField("comp", "epts").Error(err.Error())
		return strategy_config_epts{}, err
	}
	sap, found := props[_INIT_SELL_AMOUNT_PERCENTAGE]
	if !found {
		err = errors.BadRequest(logger.XXX_ERR_MISSING_PROP_KEY, _INIT_SELL_AMOUNT_PERCENTAGE)
		logrus.WithField("comp", "epts").Error(err.Error())
		return strategy_config_epts{}, err
	}
	eb, found := props[_EXPONENTIAL_BASE]
	if !found {
		err = errors.BadRequest(logger.XXX_ERR_MISSING_PROP_KEY, _EXPONENTIAL_BASE)
		logrus.WithField("comp", "epts").Error(err.Error())
		return strategy_config_epts{}, err
	}

	s.BuyPercentage = utils.DecimalFromString(bp).Round(2)
	s.SellPercentage = utils.DecimalFromString(sp).Round(2)
	s.InitBuyAmountPercentage = utils.DecimalFromString(bap).Round(2)
	s.InitSellAmountPercentage = utils.DecimalFromString(sap).Round(2)
	s.ExponentialBase = utils.DecimalFromString(eb).Round(2)

	// Checking config validity
	if s.BuyPercentage.LessThanOrEqual(decimal.Zero) ||
		s.SellPercentage.LessThanOrEqual(decimal.Zero) ||
		s.InitBuyAmountPercentage.LessThanOrEqual(decimal.Zero) ||
		s.InitSellAmountPercentage.LessThanOrEqual(decimal.Zero) {

		err = errors.BadRequest(logger.XXX_ERR_NEGATIVE_PERCENTAGES)
		logrus.WithField("comp", "epts").Error(err.Error())
		return strategy_config_epts{}, err
	}

	if s.ExponentialBase.LessThanOrEqual(utils.DecimalFromString("1")) {
		err = errors.BadRequest(logger.EPTS_ERR_EXPONENTIAL_BASE, s.ExponentialBase)
		logrus.WithField("comp", "epts").Error(err.Error())
		return strategy_config_epts{}, err
	}

	return s, nil
}
