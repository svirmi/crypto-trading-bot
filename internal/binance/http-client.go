package binance

import (
	"context"
	"fmt"
	"log"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"github.com/valerioferretti92/trading-bot-demo/internal/utils"
)

func GetAssetsValueUsdt(bases []string) (map[string]model.SymbolPrice, error) {
	lprices := make(map[string]model.SymbolPrice)

	pricesService := httpClient.NewListPricesService()
	for _, base := range bases {
		symbol := utils.GetSymbolFromAsset(base)
		_, found := symbols[symbol]
		if !found {
			log.Printf("%s is not a tradable asset: skipped", base)
			continue
		}

		rprices, err := pricesService.Symbol(symbol).Do(context.TODO())
		if err != nil {
			return nil, err
		}

		lprice, err := to_CCTB_symbol_price(rprices[0])
		if err != nil {
			return nil, err
		} else {
			lprices[lprice.Symbol] = lprice
		}
	}
	return lprices, nil
}

// GetAccount returns account inforamtion
func GetAccout() (model.RemoteAccount, error) {
	account, err := httpClient.NewGetAccountService().Do(context.TODO())
	if err != nil {
		log.Printf("%s\n", err.Error())
		return model.RemoteAccount{}, fmt.Errorf("failed to retrieve account information")
	}
	return to_CCTB_remote_account(account)
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

	order, err := ordersvc.Do(context.TODO())
	if err != nil {
		log.Printf("%s\n", err.Error())
		return fmt.Errorf("failed to place market order %s%s", base, quote)
	}
	log.Printf("symbol: %s, side: %s, qty: %f, status: %s\n", order.Symbol, order.Side, qty, order.Status)
	return nil
}

/********************** Mapping to local representation **********************/

func to_CCTB_symbol_price(rprice *binanceapi.SymbolPrice) (model.SymbolPrice, error) {
	amount, err := utils.ParseFloat32(rprice.Price)
	if err != nil {
		return model.SymbolPrice{}, err
	}

	return model.SymbolPrice{
		Symbol: rprice.Symbol,
		Price:  amount}, nil
}

func to_CCTB_remote_account(account *binanceapi.Account) (model.RemoteAccount, error) {
	balances := make([]model.RemoteBalance, 0, len(account.Balances))
	for _, rbalance := range account.Balances {
		amount, err := utils.ParseFloat32(rbalance.Free)
		if err != nil {
			return model.RemoteAccount{}, err
		}

		balances = append(balances, model.RemoteBalance{
			Asset:  rbalance.Asset,
			Amount: amount})
	}

	return model.RemoteAccount{
		MakerCommission:  account.MakerCommission,
		TakerCommission:  account.TakerCommission,
		BuyerCommission:  account.BuyerCommission,
		SellerCommission: account.SellerCommission,
		Balances:         balances}, nil
}
