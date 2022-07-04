package pts

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

const (
	_BUY_PERCENTAGE         = "buyPercentage"
	_SELL_PERCENTAGE        = "sellPercentage"
	_BUY_AMOUNT_PERCENTAGE  = "buyAmountPercentage"
	_SELL_AMOUNT_PERCENTAGE = "sellAmountPercentage"
)

type strategy_config_pts struct {
	BuyPercentage        decimal.Decimal
	SellPercentage       decimal.Decimal
	BuyAmountPercentage  decimal.Decimal
	SellAmountPercentage decimal.Decimal
}

func parse_config(props map[string]string) (s strategy_config_pts) {
	bp, found := props[_BUY_PERCENTAGE]
	if !found {
		logrus.WithField("comp", "pts").Panicf(logger.XXX_ERR_MISSING_PROP_KEY, _BUY_PERCENTAGE)
	}
	sp, found := props[_SELL_PERCENTAGE]
	if !found {
		logrus.WithField("comp", "pts").Panicf(logger.XXX_ERR_MISSING_PROP_KEY, _SELL_PERCENTAGE)
	}
	bap, found := props[_BUY_AMOUNT_PERCENTAGE]
	if !found {
		logrus.WithField("comp", "pts").Panicf(logger.XXX_ERR_MISSING_PROP_KEY, _BUY_AMOUNT_PERCENTAGE)
	}
	sap, found := props[_SELL_AMOUNT_PERCENTAGE]
	if !found {
		logrus.WithField("comp", "dts").Panicf(logger.XXX_ERR_MISSING_PROP_KEY, _SELL_AMOUNT_PERCENTAGE)
	}

	s.BuyPercentage = utils.DecimalFromString(bp).Round(2)
	s.SellPercentage = utils.DecimalFromString(sp).Round(2)
	s.BuyAmountPercentage = utils.DecimalFromString(bap).Round(2)
	s.SellAmountPercentage = utils.DecimalFromString(sap).Round(2)

	// Checking config validity
	if s.BuyPercentage.LessThanOrEqual(decimal.Zero) ||
		s.SellPercentage.LessThanOrEqual(decimal.Zero) ||
		s.BuyAmountPercentage.LessThanOrEqual(decimal.Zero) ||
		s.SellAmountPercentage.LessThanOrEqual(decimal.Zero) {

		logrus.WithField("comp", "pts").Panic(logger.PTS_ERR_NEGATIVE_PERCENTAGES)
	}

	return s
}
