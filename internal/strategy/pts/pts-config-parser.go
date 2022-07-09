package pts

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

type strategy_config_pts struct {
	BuyPercentage        decimal.Decimal
	SellPercentage       decimal.Decimal
	BuyAmountPercentage  decimal.Decimal
	SellAmountPercentage decimal.Decimal
}

func (l LocalAccountPTS) ValidateConfig(props map[string]string) error {
	_, err := parse_config(props)
	return err
}

func parse_config(props map[string]string) (s strategy_config_pts, err error) {
	bp, found := props[_BUY_PERCENTAGE]
	if !found {
		err = fmt.Errorf(logger.XXX_ERR_MISSING_PROP_KEY, _BUY_PERCENTAGE)
		logrus.WithField("comp", "pts").Error(err.Error())
		return strategy_config_pts{}, err
	}
	sp, found := props[_SELL_PERCENTAGE]
	if !found {
		err = fmt.Errorf(logger.XXX_ERR_MISSING_PROP_KEY, _SELL_PERCENTAGE)
		logrus.WithField("comp", "pts").Error(err.Error())
		return strategy_config_pts{}, err
	}
	bap, found := props[_BUY_AMOUNT_PERCENTAGE]
	if !found {
		err = fmt.Errorf(logger.XXX_ERR_MISSING_PROP_KEY, _BUY_AMOUNT_PERCENTAGE)
		logrus.WithField("comp", "pts").Error(err.Error())
		return strategy_config_pts{}, err
	}
	sap, found := props[_SELL_AMOUNT_PERCENTAGE]
	if !found {
		err = fmt.Errorf(logger.XXX_ERR_MISSING_PROP_KEY, _SELL_AMOUNT_PERCENTAGE)
		logrus.WithField("comp", "pts").Error(err.Error())
		return strategy_config_pts{}, err
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

		err = fmt.Errorf(logger.PTS_ERR_NEGATIVE_PERCENTAGES)
		logrus.WithField("comp", "pts").Error(err.Error())
		return strategy_config_pts{}, err
	}

	return s, nil
}
