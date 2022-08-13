package exchange

import (
	"context"
	"fmt"
	"time"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

var (
	httpClient           *binanceapi.Client
	symbols              map[string]binanceapi.Symbol
	mmsDoneCh, mmsStopCh chan struct{}
)

type binance_exchange struct{}

func (be binance_exchange) initialize(mmsch chan []model.MiniMarketStats, _ chan model.MiniMarketStatsAck) errors.CtbError {
	return binancex_initialize(mmsch)
}

func (be binance_exchange) can_spot_trade(symbol string) bool {
	return binancex_can_spot_trade(symbol)
}

func (be binance_exchange) get_spot_market_limits(symbol string) (model.SpotMarketLimits, errors.CtbError) {
	return binancex_get_spot_market_limits(symbol)
}

func (be binance_exchange) filter_tradable_assets(bases []string) []string {
	return binancex_filter_tradable_assets(bases)
}

func (be binance_exchange) get_assets_value(bases []string) (map[string]model.AssetPrice, errors.CtbError) {
	return binancex_get_assets_value(bases)
}

func (be binance_exchange) get_account() (model.RemoteAccount, errors.CtbError) {
	return binancex_get_account()
}

func (be binance_exchange) send_spot_market_order(op model.Operation) (model.Operation, errors.CtbError) {
	return binancex_send_spot_market_order(op)
}

func (be binance_exchange) mini_markets_stats_serve() errors.CtbError {
	return binancex_mini_markets_stats_serve()
}

func (be binance_exchange) mini_markets_stats_stop() {
	binancex_mini_markets_stats_stop()
}

func binancex_initialize(mmsch chan []model.MiniMarketStats) errors.CtbError {
	// Decoding config
	binanceConfig := struct {
		ApiKey     string
		SecretKey  string
		UseTestnet bool
	}{}
	err := mapstructure.Decode(config.GetExchangeConfig(), &binanceConfig)
	if err != nil {
		return errors.WrapBadRequest(err)
	}

	// Web socket keep alive set up
	binanceapi.WebsocketKeepalive = true
	binanceapi.WebsocketTimeout = time.Second * 60
	mmsCh = mmsch

	// Building binance http client
	binanceapi.UseTestnet = binanceConfig.UseTestnet
	httpClient = binanceapi.NewClient(binanceConfig.ApiKey, binanceConfig.SecretKey)

	// Init exchange symbols
	res, err := httpClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return errors.WrapExchange(err)
	}

	logrus.WithField("comp", "binancex").Info(logger.BINEX_REGISTERING_SYMBOLS)
	symbols = make(map[string]binanceapi.Symbol)
	for _, symbol := range res.Symbols {
		symbols[symbol.Symbol] = symbol
	}
	return nil
}

var binancex_can_spot_trade = func(symbol string) bool {
	status, found := symbols[symbol]

	if !found {
		return false
	}
	return status.Status == string(binanceapi.SymbolStatusTypeTrading) && status.IsSpotTradingAllowed
}

var binancex_get_spot_market_limits = func(symbol string) (model.SpotMarketLimits, errors.CtbError) {
	iLotSize, err := get_spot_limit_sizes(symbol)
	if err != nil {
		return model.SpotMarketLimits{}, err
	}
	iMarketLotSize, err := get_spot_market_sizes(symbol)
	if err != nil {
		return model.SpotMarketLimits{}, err
	}
	iNotional, err := get_min_notional(symbol)
	if err != nil {
		return model.SpotMarketLimits{}, err
	}

	minBase := decimal.Max(iLotSize.MinBase, iMarketLotSize.MinBase)
	maxBase := decimal.Min(iLotSize.MaxBase, iMarketLotSize.MaxBase)
	stepBase := decimal.Max(iLotSize.StepBase, iMarketLotSize.StepBase)

	return model.SpotMarketLimits{
		MinBase:  minBase,
		MaxBase:  maxBase,
		StepBase: stepBase,
		MinQuote: iNotional}, nil
}

func get_min_notional(symbol string) (decimal.Decimal, errors.CtbError) {
	status, found := symbols[symbol]

	if !found {
		err := errors.Internal(logger.BINEX_ERR_SYMBOL_NOT_FOUND, symbol)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return decimal.Zero, err
	}

	iNotional := extract_filter(status.Filters, "MIN_NOTIONAL")
	if iNotional == nil {
		err := errors.Internal(logger.BINEX_ERR_FILTER_NOT_FOUND, "MIN_NOTIONAL", symbol)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return decimal.Zero, err
	}

	return parse_number(iNotional["minNotional"], decimal.Zero), nil
}

func get_spot_market_sizes(symbol string) (model.SpotMarketLimits, errors.CtbError) {
	status, found := symbols[symbol]

	if !found {
		err := errors.Internal(logger.BINEX_ERR_SYMBOL_NOT_FOUND, symbol)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return model.SpotMarketLimits{}, err
	}

	iMarketLotSize := extract_filter(status.Filters, "MARKET_LOT_SIZE")
	if iMarketLotSize == nil {
		err := errors.Internal(logger.BINEX_ERR_FILTER_NOT_FOUND, "MARKET_LOT_SIZE", symbol)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return model.SpotMarketLimits{}, err
	}

	return model.SpotMarketLimits{
		MinBase:  parse_number(iMarketLotSize["minQty"], decimal.Zero),
		MaxBase:  parse_number(iMarketLotSize["maxQty"], utils.MaxDecimal()),
		StepBase: parse_number(iMarketLotSize["stepSize"], decimal.Zero)}, nil
}

func get_spot_limit_sizes(symbol string) (model.SpotMarketLimits, errors.CtbError) {
	status, found := symbols[symbol]

	if !found {
		err := errors.Internal(logger.BINEX_ERR_SYMBOL_NOT_FOUND, symbol)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return model.SpotMarketLimits{}, err
	}

	iLotSize := extract_filter(status.Filters, "LOT_SIZE")
	if iLotSize == nil {
		err := errors.Internal(logger.BINEX_ERR_FILTER_NOT_FOUND, "LOT_SIZE", symbol)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return model.SpotMarketLimits{}, err
	}

	return model.SpotMarketLimits{
		MinBase:  parse_number(iLotSize["minQty"], decimal.Zero),
		MaxBase:  parse_number(iLotSize["maxQty"], utils.MaxDecimal()),
		StepBase: parse_number(iLotSize["stepSize"], decimal.Zero)}, nil
}

func extract_filter(filters []map[string]interface{}, filterType string) map[string]interface{} {
	for _, filter := range filters {
		if filterType == fmt.Sprintf("%v", filter["filterType"]) {
			return filter
		}
	}
	return nil
}

func parse_number(num interface{}, def decimal.Decimal) decimal.Decimal {
	if num == nil {
		return def
	}
	str := fmt.Sprintf("%s", num)
	if str == "" {
		return def
	}
	return utils.DecimalFromString(str)
}

var binancex_filter_tradable_assets = func(bases []string) []string {
	// An asset is considered to be tradable, if it can be
	// exchanged for USDT directly
	tradables := make([]string, 0)
	for _, base := range bases {
		if base == "USDT" {
			continue
		}

		symbol, err := utils.GetSymbolFromAsset(base)
		if err != nil {
			logrus.WithField("comp", "binancex").
				Error(err.Error())
			continue
		}

		_, found := symbols[symbol]
		if !found {
			logrus.WithField("comp", "binancex").
				Warnf(logger.BINEX_NON_TRADABLE_ASSET, base)
			continue
		}
		tradables = append(tradables, base)
	}
	return tradables
}

var binancex_get_assets_value = func(bases []string) (map[string]model.AssetPrice, errors.CtbError) {
	lprices := make(map[string]model.AssetPrice)
	bases = binancex_filter_tradable_assets(bases)

	pricesService := httpClient.NewListPricesService()
	for _, base := range bases {
		symbol, _ := utils.GetSymbolFromAsset(base)
		rprices, err := binance_get_price(pricesService.Symbol(symbol))
		if err != nil {
			logrus.WithField("comp", "binancex").Error(err.Error())
			return nil, err
		}

		lprice, err := to_CCTB_symbol_price(rprices[0])
		if err != nil {
			logrus.WithField("comp", "binancex").Error(err.Error())
			return nil, err
		} else {
			lprices[lprice.Asset] = lprice
		}
	}
	return lprices, nil
}

var binancex_get_account = func() (model.RemoteAccount, errors.CtbError) {
	account, err := binance_get_account(httpClient.NewGetAccountService())
	if err != nil {
		return model.RemoteAccount{}, err
	}
	return to_CCTB_remote_account(account)
}

var binancex_send_spot_market_order = func(op model.Operation) (model.Operation, errors.CtbError) {
	// Check if symbol or its inverse exists
	_, dfound := symbols[op.Base+op.Quote]
	_, ifound := symbols[op.Quote+op.Base]
	if !dfound && !ifound {
		err := errors.Internal(logger.BINEX_ERR_INVALID_SYMBOL,
			op.Base, op.Quote, op.Quote, op.Base)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return model.Operation{}, err
	}

	// If direct symbol does not exist, invert operation
	if ifound {
		op = op.Flip()
	}

	// Checking if symbol can be traded
	if !binancex_can_spot_trade(op.Base + op.Quote) {
		err := errors.Exchange(logger.BINEX_TRADING_DISABLED, op.Base+op.Quote)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return op, err
	}

	// Execute operation
	op.Timestamp = time.Now().UnixMicro()
	err := do_send_spot_market_order(op)
	if err != nil {
		op.Status = model.FAILED
		return op, err
	}
	return op, nil
}

func do_send_spot_market_order(op model.Operation) errors.CtbError {
	// Get spot market limits
	limits, err := binancex_get_spot_market_limits(op.Base + op.Quote)
	if err != nil {
		return err
	}

	// Check market order lower bounds
	if op.AmountSide == model.QUOTE_AMOUNT && op.Amount.LessThan(limits.MinQuote) {
		err = errors.Internal(logger.BINEX_BELOW_QUOTE_LIMIT,
			op.Base+op.Quote, op.Side, op.Amount, op.AmountSide, limits.MinQuote.String())
		logrus.WithField("comp", "binancex").Error(err.Error())
		return err
	}
	if op.AmountSide == model.BASE_AMOUNT && op.Amount.LessThan(limits.MinBase) {
		err := errors.Internal(logger.BINEX_BELOW_BASE_LIMIT,
			op.Base+op.Quote, op.Side, op.Amount, op.AmountSide, limits.MinBase.String())
		logrus.WithField("comp", "binancex").Error(err.Error())
		return err
	}

	// Get market order upper bound
	var max decimal.Decimal
	var min decimal.Decimal
	if op.AmountSide == model.BASE_AMOUNT {
		max = limits.MaxBase
		min = limits.MinBase
	} else {
		max = limits.MaxBase.Mul(op.Price)
		min = limits.MinQuote
	}

	// Regular order
	if op.Amount.LessThanOrEqual(max) {
		return do_do_send_spot_market_order(op)
	}

	// Iceberg order
	failed := true
	intdiv := int(op.Amount.Div(max).IntPart())
	reminder := op.Amount.Sub(decimal.NewFromInt(int64(intdiv)).Mul(max))
	logrus.WithField("comp", "binancex").Infof(logger.BINEX_ICEBERG_ORDER,
		op.Base+op.Quote, op.Side, op.AmountSide, intdiv, max, reminder)

	for i := 1; i < intdiv; i++ {
		op.Amount = max
		err := do_do_send_spot_market_order(op)
		failed = failed && err != nil
	}

	if reminder.Equals(decimal.Zero) {
		op.Amount = max
		err := do_do_send_spot_market_order(op)
		failed = failed && err != nil
	} else if reminder.GreaterThanOrEqual(min) {
		op.Amount = max
		err := do_do_send_spot_market_order(op)
		failed = failed && err != nil
		op.Amount = reminder
		err = do_do_send_spot_market_order(op)
		failed = failed && err != nil
	} else {
		op.Amount = max.Div(utils.DecimalFromString("2"))
		err := do_do_send_spot_market_order(op)
		failed = failed && err != nil
		op.Amount = max.Div(utils.DecimalFromString("2")).Add(reminder)
		err = do_do_send_spot_market_order(op)
		failed = failed && err != nil
	}

	if failed {
		amount := decimal.NewFromInt(int64(intdiv)).Mul(max).Add(reminder)
		err := errors.Exchange(logger.BINEX_ERR_ICEBERG_ORDER_FAILED,
			op.Base+op.Quote, op.Side, amount, op.AmountSide)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return err
	}
	return nil
}

var do_do_send_spot_market_order = func(op model.Operation) errors.CtbError {
	ordersvc := httpClient.NewCreateOrderService().
		Symbol(op.Base + op.Quote).
		Type(binanceapi.OrderTypeMarket)

	if op.Side == model.BUY {
		ordersvc.Side(binanceapi.SideTypeBuy)
	} else if op.Side == model.SELL {
		ordersvc.Side(binanceapi.SideTypeSell)
	} else {
		err := errors.Internal(logger.BINEX_ERR_UNKNOWN_SIDE, op.Side)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return err
	}

	if op.AmountSide == model.BASE_AMOUNT {
		ordersvc.Quantity(op.Amount.String())
	} else {
		ordersvc.QuoteOrderQty(op.Amount.String())
	}

	order, err := binance_create_order(ordersvc)
	if err != nil {
		logrus.WithField("comp", "binancex").Error(err.Error())
		return err
	}
	logrus.WithField("comp", "binancex").
		Infof(logger.BINEX_MKT_ORDER_RESULT,
			order.Symbol,
			order.OrigQuantity,
			order.ExecutedQuantity,
			order.Status,
			order.Side)
	return nil
}

var binancex_mini_markets_stats_serve = func() errors.CtbError {
	if mmsCh == nil {
		err := errors.Internal(logger.BINEX_ERR_NIL_MMS_CH)
		logrus.WithField("comp", "binancex").Error(err.Error())
		return err
	}

	errorHandler := func(err error) {
		logrus.WithField("comp", "binancex").
			Errorf(logger.BINEX_ERR_FAILED_TO_HANLDE_MMS, err.Error())
	}

	callback := func(rMiniMarketsStats binanceapi.WsAllMiniMarketsStatEvent) {
		mmss := make([]model.MiniMarketStats, 0, len(rMiniMarketsStats))
		for _, rMiniMarketStats := range rMiniMarketsStats {
			if !utils.IsSymbolTradable(rMiniMarketStats.Symbol) {
				continue
			}

			mms, err := to_mini_market_stats(*rMiniMarketStats)
			if err != nil {
				logrus.Errorf(logger.BINEX_ERR_SKIPPING_MMS, err.Error())
				continue
			}

			mmss = append(mmss, mms)
		}

		// Return if no mini markets stats left after filtering
		if len(mmss) == 0 {
			return
		} else {
			logrus.WithField("comp", "binancex").
				Tracef(logger.BINEX_MMSS_TO_CHANNEL, utils.ToAssets(mmss))
		}

		// Send a mini markets stats through the channel
		select {
		case mmsCh <- mmss:
		default:
			logrus.WithField("comp", "binancex").
				Warnf(logger.BINEX_DROP_MMS_UPDATE, len(mmss))
		}
	}

	// Opening web socket and intialising control structure
	done, stop, err := binanceapi.WsAllMiniMarketsStatServe(callback, errorHandler)
	if err != nil {
		logrus.WithField("comp", "binancex").Error(err.Error())
		return errors.WrapExchange(err)
	} else {
		mmsDoneCh = done
		mmsStopCh = stop
	}
	return nil
}

var binancex_mini_markets_stats_stop = func() {
	if mmsStopCh == nil || mmsDoneCh == nil {
		return
	}

	logrus.WithField("comp", "binancex").Info(logger.BINEX_CLOSING_MMS)
	mmsStopCh <- struct{}{}
	<-mmsDoneCh

	if mmsCh != nil {
		close(mmsCh)
	}
}

/********************** Mapping to local representation **********************/

func to_CCTB_symbol_price(rprice *binanceapi.SymbolPrice) (model.AssetPrice, errors.CtbError) {
	amount := utils.DecimalFromString(rprice.Price)
	asset, err := utils.GetAssetFromSymbol(rprice.Symbol)
	if err != nil {
		return model.AssetPrice{}, err
	}

	return model.AssetPrice{Asset: asset, Price: amount}, nil
}

func to_CCTB_remote_account(account *binanceapi.Account) (model.RemoteAccount, errors.CtbError) {
	balances := make([]model.RemoteBalance, 0, len(account.Balances))
	for _, rbalance := range account.Balances {
		amount := utils.DecimalFromString(rbalance.Free)
		if amount.Equals(decimal.Zero) {
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

func to_mini_market_stats(rMiniMarketStat binanceapi.WsMiniMarketsStatEvent) (model.MiniMarketStats, errors.CtbError) {
	lastPrice := utils.DecimalFromString(rMiniMarketStat.LastPrice)
	openPrice := utils.DecimalFromString(rMiniMarketStat.OpenPrice)
	lowPrice := utils.DecimalFromString(rMiniMarketStat.LowPrice)
	highPrice := utils.DecimalFromString(rMiniMarketStat.HighPrice)
	baseVolume := utils.DecimalFromString(rMiniMarketStat.BaseVolume)
	quoteVolume := utils.DecimalFromString(rMiniMarketStat.QuoteVolume)
	asset, err := utils.GetAssetFromSymbol(rMiniMarketStat.Symbol)
	if err != nil {
		return model.MiniMarketStats{}, err
	}

	return model.MiniMarketStats{
		Event:       rMiniMarketStat.Event,
		Time:        rMiniMarketStat.Time,
		Asset:       asset,
		LastPrice:   lastPrice,
		OpenPrice:   openPrice,
		LowPrice:    lowPrice,
		HighPrice:   highPrice,
		BaseVolume:  baseVolume,
		QuoteVolume: quoteVolume}, nil
}

/********************** Binance calls **********************/

var binance_get_price = func(b *binanceapi.ListPricesService) ([]*binanceapi.SymbolPrice, errors.CtbError) {
	p, err := b.Do(context.TODO())
	if err != nil {
		return nil, errors.WrapExchange(err)
	}
	return p, nil
}

var binance_get_account = func(b *binanceapi.GetAccountService) (*binanceapi.Account, errors.CtbError) {
	a, err := b.Do(context.TODO())
	if err != nil {
		return nil, errors.WrapExchange(err)
	}
	return a, nil
}

var binance_create_order = func(b *binanceapi.CreateOrderService) (*binanceapi.CreateOrderResponse, errors.CtbError) {
	o, err := b.Do(context.TODO())
	if err != nil {
		return nil, errors.WrapExchange(err)
	}
	return o, nil
}
