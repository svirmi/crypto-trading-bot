package binance

import (
	"context"
	"fmt"
	"log"
	"time"

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
		rprices, err := binance_get_price(pricesService.Symbol(symbol))
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
	account, err := binance_get_account(httpClient.NewGetAccountService())
	if err != nil {
		return model.RemoteAccount{}, err
	}
	return to_CCTB_remote_account(account)
}

func SendMarketOrder(op model.Operation) (model.Operation, error) {
	//Check if symbol or its inverse exists
	_, dfound := symbols[op.Base+op.Quote]
	_, ifound := symbols[op.Quote+op.Base]
	if !dfound && !ifound {
		err := fmt.Errorf("neither %s%s nor %s%s is a valid exchange symbol",
			op.Base, op.Quote, op.Quote, op.Base)
		return model.Operation{}, err
	}

	// If direct symbol does not exist, invert operation
	if ifound {
		op = op.Flip()
	}

	// Execute operation
	op.Timestamp = time.Now().UnixMicro()
	err := send_market_order(op)
	if err != nil {
		op.Status = model.FAILED
		return op, err
	}
	return op, nil
}

func send_market_order(op model.Operation) error {
	ordersvc := httpClient.NewCreateOrderService().
		Symbol(op.Base + op.Quote).
		Type(binanceapi.OrderTypeMarket)

	if op.Side == model.BUY {
		ordersvc.Side(binanceapi.SideTypeBuy)
	} else if op.Side == model.SELL {
		ordersvc.Side(binanceapi.SideTypeSell)
	} else {
		return fmt.Errorf("unknown operation side %s", op.Side)
	}

	if op.AmountSide == model.BASE_AMOUNT {
		ordersvc.Quantity(op.Amount.String())
	} else {
		ordersvc.QuoteOrderQty(op.Amount.String())
	}

	order, err := binance_create_order(ordersvc)
	if err != nil {
		return err
	}
	log.Printf("symbol: %s, side: %s, status: %s\n", order.Symbol, order.Side, order.Status)
	return nil
}

/********************** Mapping to local representation **********************/

func to_CCTB_symbol_price(rprice *binanceapi.SymbolPrice) (model.AssetPrice, error) {
	amount := utils.DecimalFromString(rprice.Price)

	return model.AssetPrice{
		Asset: utils.GetAssetFromSymbol(rprice.Symbol),
		Price: amount}, nil
}

func to_CCTB_remote_account(account *binanceapi.Account) (model.RemoteAccount, error) {
	balances := make([]model.RemoteBalance, 0, len(account.Balances))
	for _, rbalance := range account.Balances {
		amount := utils.DecimalFromString(rbalance.Free)
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

/********************** Binance calls **********************/

var binance_get_price = func(b *binanceapi.ListPricesService) ([]*binanceapi.SymbolPrice, error) {
	return b.Do(context.TODO())
}

var binance_get_account = func(b *binanceapi.GetAccountService) (*binanceapi.Account, error) {
	return b.Do(context.TODO())
}

var binance_create_order = func(b *binanceapi.CreateOrderService) (*binanceapi.CreateOrderResponse, error) {
	return b.Do(context.TODO())
}
