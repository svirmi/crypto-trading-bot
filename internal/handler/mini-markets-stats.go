package handler

import (
	"log"
	"time"

	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

func HandleMiniMarketsStats(miniMarketsStats []model.MiniMarketStats) {
	msg := "received price update: %s %f"
	for _, miniMarketStats := range miniMarketsStats {
		log.Printf(msg, miniMarketStats.Symbol, miniMarketStats.LastPrice)
	}
	time.Sleep(5 * time.Second)
}
