package binance

import (
	"log"

	binanceapi "github.com/adshao/go-binance/v2"
)

var (
	doneC, stopC chan struct{}
	err          error
)

func MiniMarketsStatServe(handler func(binanceapi.WsAllMiniMarketsStatEvent)) error {
	errHandler := func(err error) {
		log.Printf("%s\n", err.Error())
	}
	doneC, stopC, err = binanceapi.WsAllMiniMarketsStatServe(handler, errHandler)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return err
	}
	return nil
}

func Close() {
	log.Printf("closing price update web socket")
	stopC <- struct{}{}
	<-doneC
}
