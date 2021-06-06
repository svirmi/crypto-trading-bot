package binance

import (
	"context"
	"fmt"

	binanceapi "github.com/adshao/go-binance/v2"
)

func GetAccout() {
	res, err := httpClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("MakerCommission: %d\n", res.MakerCommission)
	fmt.Printf("TakerCommission: %d\n", res.TakerCommission)
	fmt.Printf("BuyerCommission: %d\n", res.BuyerCommission)
	fmt.Printf("SellerCommission: %d\n", res.SellerCommission)
	fmt.Printf("CanTrade: %t\n", res.CanTrade)
	fmt.Printf("CanWithdraw: %t\n", res.CanWithdraw)
	fmt.Printf("CanDeposit: %t\n", res.CanDeposit)
	fmt.Printf("Balances: %q\n", res.Balances)
}

func SendMarketOrder(base, quote string, qty float64) error {
	dsymbol, dfound := symbols[base+quote]
	isymbol, ifound := symbols[quote+base]
	if !dfound && !ifound {
		return fmt.Errorf("currency pair %s - %s not found", base, quote)
	}

	if dfound {
		sendMarketOrder(base, quote, qty, binanceapi.SideTypeBuy, dsymbol)
	} else {
		sendMarketOrder(quote, base, qty, binanceapi.SideTypeSell, isymbol)
	}
	return nil
}

func sendMarketOrder(base, quote string, qty float64, side binanceapi.SideType, symbol binanceapi.Symbol) {
	ordersvc := httpClient.NewCreateOrderService().
		Symbol(base + quote).
		Type(binanceapi.OrderTypeMarket).
		Side(side).
		Quantity(fmt.Sprintf("%f", qty))

	order, err := ordersvc.Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Symbol: %s, Side: %s, Qty: %f, Status: %s\n", order.Symbol, order.Side, qty, order.Status)
}
