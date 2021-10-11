package binance

import (
	"context"
	"fmt"
	"log"

	binanceapi "github.com/adshao/go-binance/v2"
)

// GetAccount returns account inforamtion
func GetAccout() (*binanceapi.Account, error) {
	account, err := httpClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Printf("%s\n", err.Error())
		return nil, fmt.Errorf("failed to retrieve account information")
	}
	return account, nil
}

// SendMarketOrder places a market order to obtain qty units of target
// currency, paying with exchange currency. Internally, it will map
// the target - exchange pair into a binance base - quote pair.
func SendMarketOrder(target, exchange string, qty float64) (err error) {
	_, dfound := symbols[target+exchange]
	_, ifound := symbols[exchange+target]
	if !dfound && !ifound {
		err = fmt.Errorf("neither %s%s nor %s%s is a valid exchange symbol",
			target, exchange, exchange, target)
		return err
	}

	if dfound {
		return sendMarketOrder(target, exchange, qty, true, binanceapi.SideTypeBuy)
	} else {
		return sendMarketOrder(exchange, target, qty, false, binanceapi.SideTypeSell)
	}
}

func sendMarketOrder(base, quote string, qty float64, regular bool, side binanceapi.SideType) error {
	ordersvc := httpClient.NewCreateOrderService().
		Symbol(base + quote).
		Type(binanceapi.OrderTypeMarket).
		Side(side)
	if regular {
		ordersvc.Quantity(fmt.Sprintf("%f", qty))
	} else {
		ordersvc.QuoteOrderQty(fmt.Sprintf("%f", qty))
	}

	order, err := ordersvc.Do(context.Background())
	if err != nil {
		log.Printf("%s\n", err.Error())
		return fmt.Errorf("failed to place market order %s%s", base, quote)
	}
	log.Printf("symbol: %s, side: %s, qty: %f, status: %s\n", order.Symbol, order.Side, qty, order.Status)
	return nil
}
