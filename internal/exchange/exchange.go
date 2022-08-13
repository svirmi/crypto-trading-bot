package exchange

import (
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

type iexchange interface {
	initialize(chan []model.MiniMarketStats, chan model.MiniMarketStatsAck) errors.CtbError
	can_spot_trade(string) bool
	get_spot_market_limits(string) (model.SpotMarketLimits, errors.CtbError)
	filter_tradable_assets([]string) []string
	get_assets_value([]string) (map[string]model.AssetPrice, errors.CtbError)
	get_account() (model.RemoteAccount, errors.CtbError)
	send_spot_market_order(model.Operation) (model.Operation, errors.CtbError)
	mini_markets_stats_serve() errors.CtbError
	mini_markets_stats_stop()
}

var (
	exchange iexchange
	mmsCh    chan []model.MiniMarketStats
	cllCh    chan model.MiniMarketStatsAck
)

func Initialize(extype model.ExchangeType, mmsch chan []model.MiniMarketStats, cllch chan model.MiniMarketStatsAck) errors.CtbError {
	if exchange != nil {
		return nil
	}

	if extype == model.LOCALEX {
		exchange = local_exchange{}
	} else if extype == model.BINANCEX {
		exchange = binance_exchange{}
	} else {
		err := errors.Internal(logger.EX_ERR_UNKNOWN_EXTYPE, extype)
		logrus.Error(err.Error())
		return err
	}
	return exchange.initialize(mmsch, cllch)
}

func CanSpotTrade(symbol string) bool {
	if exchange == nil {
		err := errors.Internal(logger.EX_ERR_UNINITIALIZED)
		logrus.Panic(err.Error())
	}
	return exchange.can_spot_trade(symbol)
}

func GetSpotMarketLimits(symbol string) (model.SpotMarketLimits, errors.CtbError) {
	if exchange == nil {
		err := errors.Internal(logger.EX_ERR_UNINITIALIZED)
		logrus.Panic(err.Error())
	}
	return exchange.get_spot_market_limits(symbol)
}

func FilterTradableAssets(bases []string) []string {
	if exchange == nil {
		err := errors.Internal(logger.EX_ERR_UNINITIALIZED)
		logrus.Panic(err.Error())
	}
	return exchange.filter_tradable_assets(bases)
}

func GetAssetsValue(bases []string) (map[string]model.AssetPrice, errors.CtbError) {
	if exchange == nil {
		err := errors.Internal(logger.EX_ERR_UNINITIALIZED)
		logrus.Panic(err.Error())
	}
	return exchange.get_assets_value(bases)
}

func GetAccount() (model.RemoteAccount, errors.CtbError) {
	if exchange == nil {
		err := errors.Internal(logger.EX_ERR_UNINITIALIZED)
		logrus.Panic(err.Error())
	}
	return exchange.get_account()
}

func SendSpotMarketOrder(op model.Operation) (model.Operation, errors.CtbError) {
	if exchange == nil {
		err := errors.Internal(logger.EX_ERR_UNINITIALIZED)
		logrus.Panic(err.Error())
	}
	return exchange.send_spot_market_order(op)
}

func MiniMarketsStatsServe() errors.CtbError {
	if exchange == nil {
		err := errors.Internal(logger.EX_ERR_UNINITIALIZED)
		logrus.Panic(err.Error())
	}
	return exchange.mini_markets_stats_serve()
}

func MiniMarketsStatsStop() {
	if exchange == nil {
		err := errors.Internal(logger.EX_ERR_UNINITIALIZED)
		logrus.Panic(err.Error())
	}
	exchange.mini_markets_stats_stop()
}
