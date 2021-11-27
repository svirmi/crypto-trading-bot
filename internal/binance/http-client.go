package binance

import (
	"context"
	"fmt"
	"log"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func FilterTradableAssets(bases []string) []string {
	// An asset is considered to be tradable, if it can be
	// exchanged for USDT directly
	tradables := make([]string, 0)
	for _, base := range bases {
		_, found := symbols[utils.GetSymbolFromAsset(base)]
		if !found {
			log.Printf("%s is not a tradable asset", base)
			continue
		}
		tradables = append(tradables, base)
	}
	return tradables
}

func GetAssetsValue(bases []string) (map[string]model.AssetPrice, error) {
	lprices := make(map[string]model.AssetPrice)
	bases = FilterTradableAssets(bases)

	pricesService := httpClient.NewListPricesService()
	for _, base := range bases {
		symbol := utils.GetSymbolFromAsset(base)
		rprices, err := pricesService.Symbol(symbol).Do(context.TODO())
		if err != nil {
			return nil, err
		}

		lprice, err := to_CCTB_symbol_price(rprices[0])
		if err != nil {
			return nil, err
		} else {
			lprices[lprice.Asset] = lprice
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

func to_CCTB_symbol_price(rprice *binanceapi.SymbolPrice) (model.AssetPrice, error) {
	amount, err := utils.ParseFloat32(rprice.Price)
	if err != nil {
		return model.AssetPrice{}, err
	}

	return model.AssetPrice{
		Asset: utils.GetAssetFromSymbol(rprice.Symbol),
		Price: amount}, nil
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
