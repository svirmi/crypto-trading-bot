package binance

import (
	"fmt"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

type mini_markets_stats_ctl struct {
	mms_done, mms_stop chan struct{}
	mms                chan []model.MiniMarketStats
}

var mms_ctl mini_markets_stats_ctl = mini_markets_stats_ctl{}

func InitMmsChannel(mmsChannel chan []model.MiniMarketStats) {
	mms_ctl.mms = mmsChannel
}

func MiniMarketsStatsServe(assets []string) error {
	if mms_ctl.mms == nil {
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
			// Filter out symbols whose numeric fields could not be parsed from string
			miniMarketStats := to_mini_market_stats(*rMiniMarketStats)
			miniMarketsStats = append(miniMarketsStats, miniMarketStats)
		}

		// Return if no mini markets stats left after filtering
		if len(miniMarketsStats) == 0 {
			return
		}

		// Send a mini markets stats through the channel
		select {
		case mms_ctl.mms <- miniMarketsStats:
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
		mms_ctl.mms_done = done
		mms_ctl.mms_stop = stop
	}
	return nil
}

func MiniMarketsStatsStop() {
	if mms_ctl.mms_stop == nil || mms_ctl.mms_done == nil {
		return
	}

	logrus.WithField("comp", "binance").Info(logger.BINANACE_CLOSING_MMS)
	mms_ctl.mms_stop <- struct{}{}
	<-mms_ctl.mms_done

	if mms_ctl.mms != nil {
		close(mms_ctl.mms)
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
