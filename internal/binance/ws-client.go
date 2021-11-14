package binance

import (
	"log"

	binanceapi "github.com/adshao/go-binance/v2"
	abool "github.com/tevino/abool/v2"
	"github.com/valerioferretti92/trading-bot-demo/internal/handler"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"github.com/valerioferretti92/trading-bot-demo/internal/utils"
)

type mini_markets_stats_ctl struct {
	done, stop chan struct{}
	closed     *abool.AtomicBool
}

var mms_ctl mini_markets_stats_ctl = mini_markets_stats_ctl{}

func MiniMarketsStatsServe(symbols map[string]bool) error {
	errorHandler := func(err error) {
		log.Print(err.Error())
	}

	sentinel := abool.New()
	callback := func(rMiniMarketsStats binanceapi.WsAllMiniMarketsStatEvent) {
		// Avoid mini markets stats race conditions
		ok := sentinel.SetToIf(false, true)
		if !ok {
			log.Printf("skipping mini markets stats update...")
			return
		}

		// Processing mini markets stats update
		go func() {
			defer sentinel.UnSet()
			miniMarketsStats := make([]model.MiniMarketStats, 0, len(rMiniMarketsStats))
			for _, rMiniMarketStats := range rMiniMarketsStats {
				// Filter out symbols that are not in local wallet
				if !symbols[rMiniMarketStats.Symbol] {
					continue
				}
				// Filter out symbols whose numeric fields could not be parsed from string
				miniMarketStats, err := to_mini_market_stats(*rMiniMarketStats)
				if err != nil {
					log.Println(err.Error())
					log.Printf("skipping update for symbol %s", rMiniMarketStats.Symbol)
					continue
				}
				miniMarketsStats = append(miniMarketsStats, miniMarketStats)
			}
			handler.HandleMiniMarketsStats(miniMarketsStats)
		}()
	}

	// Opening web socket and intialising control structure
	done, stop, err := binanceapi.WsAllMiniMarketsStatServe(callback, errorHandler)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return err
	} else {
		mms_ctl.done = done
		mms_ctl.stop = stop
		mms_ctl.closed = abool.New()
	}
	return nil
}

func MiniMarketsStatsStop() {
	if mms_ctl.stop == nil || mms_ctl.done == nil || mms_ctl.closed.IsSet() {
		return
	}

	log.Printf("closing mini markets stats ws")
	mms_ctl.closed.Set()
	mms_ctl.stop <- struct{}{}
	<-mms_ctl.done
}

func Close() {
	MiniMarketsStatsStop()
}

/********************** Mapping to local representation **********************/

func to_mini_market_stats(rMiniMarketStat binanceapi.WsMiniMarketsStatEvent) (model.MiniMarketStats, error) {
	lastPrice, err := utils.ParseFloat32(rMiniMarketStat.LastPrice)
	if err != nil {
		return model.MiniMarketStats{}, err
	}
	openPrice, err := utils.ParseFloat32(rMiniMarketStat.OpenPrice)
	if err != nil {
		return model.MiniMarketStats{}, err
	}
	lowPrice, err := utils.ParseFloat32(rMiniMarketStat.LowPrice)
	if err != nil {
		return model.MiniMarketStats{}, err
	}
	highPrice, err := utils.ParseFloat32(rMiniMarketStat.HighPrice)
	if err != nil {
		return model.MiniMarketStats{}, err
	}
	baseVolume, err := utils.ParseFloat32(rMiniMarketStat.BaseVolume)
	if err != nil {
		return model.MiniMarketStats{}, err
	}
	quoteVolume, err := utils.ParseFloat32(rMiniMarketStat.QuoteVolume)
	if err != nil {
		return model.MiniMarketStats{}, err
	}

	return model.MiniMarketStats{
		Event:       rMiniMarketStat.Event,
		Time:        rMiniMarketStat.Time,
		Symbol:      rMiniMarketStat.Symbol,
		LastPrice:   lastPrice,
		OpenPrice:   openPrice,
		LowPrice:    lowPrice,
		HighPrice:   highPrice,
		BaseVolume:  baseVolume,
		QuoteVolume: quoteVolume}, nil
}
