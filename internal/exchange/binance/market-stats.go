package binance

import (
	"fmt"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

var mini_markets_stats_serve = func() error {
	if mmsCh == nil {
		err := fmt.Errorf(logger.BINEX_ERR_NIL_MMS_CH)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return err
	}

	errorHandler := func(err error) {
		logrus.WithField("comp", "binancex").
			Errorf(logger.BINEX_ERR_FAILED_TO_HANLDE_MMS, err.Error())
	}

	callback := func(rMiniMarketsStats binanceapi.WsAllMiniMarketsStatEvent) {
		miniMarketsStats := make([]model.MiniMarketStats, 0, len(rMiniMarketsStats))
		for _, rMiniMarketStats := range rMiniMarketsStats {
			miniMarketStats := to_mini_market_stats(*rMiniMarketStats)
			miniMarketsStats = append(miniMarketsStats, miniMarketStats)
		}

		// Return if no mini markets stats left after filtering
		if len(miniMarketsStats) == 0 {
			return
		}

		// Send a mini markets stats through the channel
		select {
		case mmsCh <- miniMarketsStats:
		default:
			logrus.WithField("comp", "binancex").
				Warnf(logger.BINEX_DROP_MMS_UPDATE, len(miniMarketsStats))
		}

	}

	// Opening web socket and intialising control structure
	done, stop, err := binanceapi.WsAllMiniMarketsStatServe(callback, errorHandler)
	if err != nil {
		logrus.WithField("comp", "binancex").Error(err.Error())
		return err
	} else {
		mmsDoneCh = done
		mmsStopCh = stop
	}
	return nil
}

var mini_markets_stats_stop = func() {
	if mmsStopCh == nil || mmsDoneCh == nil {
		return
	}

	logrus.WithField("comp", "binancex").Info(logger.BINEX_CLOSING_MMS)
	mmsStopCh <- struct{}{}
	<-mmsDoneCh

	if mmsCh != nil {
		close(mmsCh)
	}
}

/********************** Mapping to local representation **********************/

func to_mini_market_stats(rMiniMarketStat binanceapi.WsMiniMarketsStatEvent) model.MiniMarketStats {
	lastPrice := utils.DecimalFromString(rMiniMarketStat.LastPrice)
	openPrice := utils.DecimalFromString(rMiniMarketStat.OpenPrice)
	lowPrice := utils.DecimalFromString(rMiniMarketStat.LowPrice)
	highPrice := utils.DecimalFromString(rMiniMarketStat.HighPrice)
	baseVolume := utils.DecimalFromString(rMiniMarketStat.BaseVolume)
	quoteVolume := utils.DecimalFromString(rMiniMarketStat.QuoteVolume)

	return model.MiniMarketStats{
		Event:       rMiniMarketStat.Event,
		Time:        rMiniMarketStat.Time,
		Asset:       utils.GetAssetFromSymbol(rMiniMarketStat.Symbol),
		LastPrice:   lastPrice,
		OpenPrice:   openPrice,
		LowPrice:    lowPrice,
		HighPrice:   highPrice,
		BaseVolume:  baseVolume,
		QuoteVolume: quoteVolume}
}
