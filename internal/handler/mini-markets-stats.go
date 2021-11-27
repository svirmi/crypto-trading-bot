package handler

import (
	"log"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

func HandleMiniMarketsStats(miniMarketsStats []model.MiniMarketStats) {
	msg := "received price update: %s %f"
	for _, miniMarketStats := range miniMarketsStats {
		log.Printf(msg, miniMarketStats.Asset, miniMarketStats.LastPrice)
	}
}
