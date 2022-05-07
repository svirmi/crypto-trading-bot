package binance

import (
	"context"
	"fmt"
	"time"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

var filter_tradable_assets = func(bases []string) []string {
	// An asset is considered to be tradable, if it can be
	// exchanged for USDT directly
	tradables := make([]string, 0)
	for _, base := range bases {
		_, found := symbols[utils.GetSymbolFromAsset(base)]
		if !found {
			logrus.WithField("comp", "binance").
				Warnf(logger.BINANCE_STABLECOIN_ASSET, base)
			continue
		}
		tradables = append(tradables, base)
	}
	return tradables
}

var get_assets_value = func(bases []string) (map[string]model.AssetPrice, error) {
	lprices := make(map[string]model.AssetPrice)
	bases = filter_tradable_assets(bases)

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
var get_account = func() (model.RemoteAccount, error) {
	account, err := binance_get_account(httpClient.NewGetAccountService())
	if err != nil {
		return model.RemoteAccount{}, err
	}
	return to_CCTB_remote_account(account)
}

var send_spot_market_order = func(op model.Operation) (model.Operation, error) {
	// Check if symbol or its inverse exists
	_, dfound := symbols[op.Base+op.Quote]
	_, ifound := symbols[op.Quote+op.Base]
	if !dfound && !ifound {
		err := fmt.Errorf(logger.BINANCE_ERR_INVALID_SYMBOL, op.Base, op.Quote, op.Quote, op.Base)
		logrus.WithField("comp", "binance").Error(err.Error())
		return model.Operation{}, err
	}

	// Check spot market order limits
	err := check_spot_market_order(op)
	if err != nil {
		return op, err
	}

	// If direct symbol does not exist, invert operation
	if ifound {
		op = op.Flip()
	}

	// Checking if symbol can be traded
	if !can_spot_trade(op.Base + op.Quote) {
		err := fmt.Errorf(logger.BINANCE_TRADING_DISABLED, op.Base+op.Quote)
		logrus.WithField("comp", "binance").Error(err.Error())
		return op, err
	}

	// Execute operation
	op.Timestamp = time.Now().UnixMicro()
	err = do_send_spot_market_order(op)
	if err != nil {
		op.Status = model.FAILED
		return op, err
	}
	return op, nil
}

func check_spot_market_order(op model.Operation) error {
	limit, err := get_spot_market_limits(op.Base + op.Quote)
	if err != nil {
		return err
	}

	if op.AmountSide == model.QUOTE_AMOUNT && op.Amount.LessThan(limit.MinQuote) {
		err = fmt.Errorf(logger.BINANCE_BELOW_QUOTE_LIMIT,
			op.Base+op.Quote, op.Side, op.Amount, op.AmountSide, limit.MinQuote.String())
		logrus.WithField("comp", "binance").Error(err.Error())
		return err
	}
	if op.AmountSide == model.BASE_AMOUNT && op.Amount.LessThan(limit.MinBase) {
		err := fmt.Errorf(logger.BINANCE_BELOW_BASE_LIMIT,
			op.Base+op.Quote, op.Side, op.Amount, op.AmountSide, limit.MinBase.String())
		logrus.WithField("comp", "binance").Error(err.Error())
		return err
	}
	if op.AmountSide == model.BASE_AMOUNT && op.Amount.GreaterThan(limit.MaxBase) {
		err := fmt.Errorf(logger.BINANCE_ABOVE_BASE_LIMIT,
			op.Base+op.Quote, op.Side, op.Amount, op.AmountSide, limit.MaxBase.String())
		logrus.WithField("comp", "binance").Error(err.Error())
		return err
	}
	return nil
}

func do_send_spot_market_order(op model.Operation) error {
	ordersvc := httpClient.NewCreateOrderService().
		Symbol(op.Base + op.Quote).
		Type(binanceapi.OrderTypeMarket)

	if op.Side == model.BUY {
		ordersvc.Side(binanceapi.SideTypeBuy)
	} else if op.Side == model.SELL {
		ordersvc.Side(binanceapi.SideTypeSell)
	} else {
		err := fmt.Errorf(logger.BINANCE_ERR_UNKNOWN_SIDE, op.Side)
		logrus.WithField("comp", "binance").Error(err.Error())
		return err
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
	logrus.WithField("comp", "binance").
		Infof(logger.BINANCE_MKT_ORDER_RESULT,
			order.Symbol,
			order.OrigQuantity,
			order.ExecutedQuantity,
			order.Status,
			order.Side)
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
		if amount.Equals(decimal.Zero) {
			logrus.WithField("comp", "binance").
				Warnf(logger.BINANACE_ZERO_AMOUNT_ASSET, rbalance.Asset)
			continue
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
