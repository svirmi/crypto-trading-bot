package handler

import (
	"log"

	binanceapi "github.com/adshao/go-binance/v2"
)

func HandlePriceUpdate(marketEvent binanceapi.WsAllMiniMarketsStatEvent) {
	msg := "received price update: %s = %s UDS"
	for i := range marketEvent {
		log.Printf(msg, marketEvent[i].Symbol, marketEvent[i].LastPrice)
	}
}
