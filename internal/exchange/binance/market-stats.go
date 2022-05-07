package binance

import (
	"fmt"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

var mini_markets_stats_serve = func(assets []string) error {
	if mms == nil {
		err := fmt.Errorf(logger.BINANCE_ERR_NIL_MMS_CH)
		logrus.WithField("comp", "binance").Error(err.Error())
		return err
	}

	symbolsMap := make(map[string]bool)
	for _, asset := range assets {
		symbolsMap[utils.GetSymbolFromAsset(asset)] = true
	}

	errorHandler := func(err error) {
		logrus.WithField("comp", "binance").
			Errorf(logger.BINANCE_ERR_FAILED_TO_HANLDE_MMS, err.Error())
	}

	callback := func(rMiniMarketsStats binanceapi.WsAllMiniMarketsStatEvent) {
		miniMarketsStats := make([]model.MiniMarketStats, 0, len(assets))
		for _, rMiniMarketStats := range rMiniMarketsStats {
			// Filter out symbols that are not in local wallet
			if !symbolsMap[rMiniMarketStats.Symbol] {
				continue
			}
			miniMarketStats := to_mini_market_stats(*rMiniMarketStats)
			miniMarketsStats = append(miniMarketsStats, miniMarketStats)
		}

		// Return if no mini markets stats left after filtering
		if len(miniMarketsStats) == 0 {
			return
		}

		// Send a mini markets stats through the channel
		select {
		case mms <- miniMarketsStats:
		default:
			logrus.WithField("comp", "binance").
				Warnf(logger.BINANCE_DROP_MMS_UPDATE, len(miniMarketsStats))
		}

	}

	// Opening web socket and intialising control structure
	done, stop, err := binanceapi.WsAllMiniMarketsStatServe(callback, errorHandler)
	if err != nil {
		logrus.WithField("comp", "binance").Error(err.Error())
		return err
	} else {
		mms_done = done
		mms_stop = stop
	}
	return nil
}

var mini_markets_stats_stop = func() {
	if mms_stop == nil || mms_done == nil {
		return
	}

	logrus.WithField("comp", "binance").Info(logger.BINANACE_CLOSING_MMS)
	mms_stop <- struct{}{}
	<-mms_done

	if mms != nil {
		close(mms)
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
