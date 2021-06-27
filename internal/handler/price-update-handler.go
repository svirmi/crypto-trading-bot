package handler

import (
	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/valerioferretti92/trading-bot-demo/internal/repository"
)

func HandlePriceUpdate(marketEvent binanceapi.WsAllMiniMarketsStatEvent) {
	repository.UpsertMiniMarketsStat(marketEvent)
}
