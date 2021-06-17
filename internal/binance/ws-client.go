package binance

import (
	"fmt"
	"time"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/valerioferretti92/trading-bot-demo/internal/repository"
)

func MiniMarketsStatServe() {
	wsBookTickerEventHandler := func(marketEvent binanceapi.WsAllMiniMarketsStatEvent) {
		repository.Insert(marketEvent)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, stopC, err := binanceapi.WsAllMiniMarketsStatServe(wsBookTickerEventHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	// use stopC to exit
	go func() {
		time.Sleep(10 * time.Second)
		stopC <- struct{}{}
	}()
	<-doneC
}
