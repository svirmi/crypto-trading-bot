package binance

import (
	"fmt"
	"time"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/fatih/color"
	"github.com/valerioferretti92/trading-bot-demo/internal/config"
)

func BookTickerServe() {
	binanceapi.UseTestnet = false
	defer func() {
		binanceapi.UseTestnet = config.AppConfig.BinanceApi.UseTestnet
	}()

	wsDepthHandler := func(event *binanceapi.WsBookTickerEvent) {
		color.White("Symbol: %s", event.Symbol)
		color.Red("  - BestBidPrice: %s, BestBidQuantity: %s", event.BestBidPrice, event.BestBidQty)
		color.Green("  - BestAskPrice: %s, BestAskQty: %s", event.BestAskPrice, event.BestAskQty)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, stopC, err := binanceapi.WsBookTickerServe("ETHEUR", wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	// use stopC to exit
	go func() {
		time.Sleep(5 * time.Second)
		stopC <- struct{}{}
	}()
	<-doneC
}
