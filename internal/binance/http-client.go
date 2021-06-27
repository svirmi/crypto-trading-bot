package binance

import (
	"context"
	"fmt"
	"log"
	"strconv"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/valerioferretti92/trading-bot-demo/internal/repository"
)

func GetAccout() (*binanceapi.Account, error) {
	account, err := httpClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Printf("%s\n", err.Error())
		return nil, fmt.Errorf("failed to retrieve account information")
	}
	return account, nil
}

func SendMarketOrder(base, quote string, qty float64) (err error) {
	dsymbol, dfound := symbols[base+quote]
	isymbol, ifound := symbols[quote+base]
	if !dfound && !ifound {
		err = fmt.Errorf("currency pair %s%s missing from symbol map", base, quote)
		return err
	}

	if dfound {
		sendMarketOrder(base, quote, qty, binanceapi.SideTypeBuy, dsymbol)
		return nil
	}

	symbol, err := repository.FindSymbolByPair(quote + base)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return fmt.Errorf("unable to get a price for %s%s", quote, base)
	}
	iprice, err := strconv.ParseFloat(symbol.LastPrice, 64)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return fmt.Errorf("unable to parse price %s into a float", symbol.LastPrice)
	}
	return sendMarketOrder(quote, base, qty*(1/iprice), binanceapi.SideTypeSell, isymbol)
}

func sendMarketOrder(base, quote string, qty float64, side binanceapi.SideType, symbol binanceapi.Symbol) error {
	ordersvc := httpClient.NewCreateOrderService().
		Symbol(base + quote).
		Type(binanceapi.OrderTypeMarket).
		Side(side).
		Quantity(fmt.Sprintf("%f", qty))

	order, err := ordersvc.Do(context.Background())
	if err != nil {
		log.Printf("%s\n", err.Error())
		return fmt.Errorf("failed to place market order %s%s", base, quote)
	}
	log.Printf("symbol: %s, side: %s, qty: %f, status: %s\n", order.Symbol, order.Side, qty, order.Status)
	return nil
}
