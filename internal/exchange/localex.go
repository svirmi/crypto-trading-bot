package exchange

import (
	"encoding/csv"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	crrqueue "github.com/Workiva/go-datastructures/queue"
	crrmap "github.com/golangltd/go-concurrentmap"
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
	raccount *crrmap.ConcurrentMap
	prices   map[string]*crrqueue.Queue
)

type local_exchange struct{}

func (le local_exchange) initialize(mmsch chan []model.MiniMarketStats, cllch chan model.MiniMarketStatsAck) errors.CtbError {
	return localex_initialize(mmsch, cllch)
}

func (le local_exchange) can_spot_trade(symbol string) bool {
	return localex_can_spot_trade(symbol)
}

func (le local_exchange) get_spot_market_limits(symbol string) (model.SpotMarketLimits, errors.CtbError) {
	return localex_get_spot_market_limits(symbol)
}

func (le local_exchange) filter_tradable_assets(bases []string) []string {
	return localex_filter_tradable_assets(bases)
}

func (le local_exchange) get_assets_value(bases []string) (map[string]model.AssetPrice, errors.CtbError) {
	return localex_get_assets_value(bases)
}

func (le local_exchange) get_account() (model.RemoteAccount, errors.CtbError) {
	return localex_get_account()
}

func (le local_exchange) send_spot_market_order(op model.Operation) (model.Operation, errors.CtbError) {
	return localex_send_spot_market_order(op)
}

func (le local_exchange) mini_markets_stats_serve() errors.CtbError {
	return localex_mini_markets_stats_serve()
}

func (le local_exchange) mini_markets_stats_stop() {
	localex_mini_markets_stats_stop()
}

func localex_initialize(mmsChannel chan []model.MiniMarketStats, cllChannel chan model.MiniMarketStatsAck) errors.CtbError {
	// Decoding config
	localExchangeConfig := struct {
		InitialBalances map[string]string
		PriceFilepaths  map[string]string
	}{}
	err := mapstructure.Decode(config.GetExchangeConfig(), &localExchangeConfig)
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return errors.WrapBadRequest(err)
	}

	// The wallet must be populated with assets, not symbols
	// For each asset in the wallet, but USDT, a price must be specified
	for asset := range localExchangeConfig.InitialBalances {
		if asset == "USDT" {
			continue
		}

		if strings.HasSuffix(asset, "USDT") {
			err := errors.BadRequest(logger.LOCALEX_ERR_INVALID_ASSET, asset)
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}

		symbol, err := utils.GetSymbolFromAsset(asset)
		if err != nil {
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}

		if _, found := localExchangeConfig.PriceFilepaths[symbol]; !found {
			err := errors.BadRequest(logger.LOCALEX_ERR_PRICES_NOT_PROVIDED, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}
	}

	// Prices are to be provided per symbol, the form XXXUSDT
	// USDTUSDT is not a valid symbol
	// Skip prices whose corresponding asset is not in the wallet
	for symbol := range localExchangeConfig.PriceFilepaths {
		if !utils.IsSymbolTradable(symbol) {
			err := errors.BadRequest(logger.LOCALEX_ERR_INVALID_SYMBOL, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}

		base := strings.TrimSuffix(symbol, "USDT")
		if base == "USDT" {
			err := errors.BadRequest(logger.LOCALEX_ERR_INVALID_SYMBOL, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}

		asset, err := utils.GetAssetFromSymbol(symbol)
		if err != nil {
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}

		if _, found := localExchangeConfig.InitialBalances[asset]; !found {
			logrus.WithField("comp", "localex").
				Warnf(logger.LOCALEX_SKIP_SYMBOL_PRICES, asset, symbol)
			delete(localExchangeConfig.PriceFilepaths, symbol)
		}
	}

	// Initializing mms channel
	mmsCh = mmsChannel
	cllCh = cllChannel

	// Parsing account balances
	logrus.Infof(logger.LOCALEX_INIT_RACCOUNT, len(localExchangeConfig.InitialBalances))
	raccount = crrmap.NewConcurrentMap()
	for key, value := range localExchangeConfig.InitialBalances {
		_, err := raccount.Put(key, utils.DecimalFromString(value))
		if err != nil {
			logrus.WithField("comp", "localex").Error(err.Error())
			return errors.WrapExchange(err)
		}
	}

	// Parsing asset prices
	prices = make(map[string]*crrqueue.Queue)
	for symbol, priceFilepath := range localExchangeConfig.PriceFilepaths {
		err := parse_prices_file(symbol, priceFilepath)
		if err != nil {
			return err
		}
	}
	return nil
}

func parse_prices_file(symbol, priceFilepath string) errors.CtbError {
	// Parsing priceFilepath
	priceFilepath = os.ExpandEnv(priceFilepath)

	// Opening price file
	file, err := os.Open(priceFilepath)
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return errors.WrapBadRequest(err)
	}
	defer file.Close()

	// Compiling regexps
	intre, err := regexp.Compile(utils.GetIntegerRegexp())
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return errors.WrapExchange(err)
	}
	floatre, err := regexp.Compile(utils.GetFloatRegexp())
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return errors.WrapExchange(err)
	}

	// Reading price file
	logrus.WithField("comp", "localex").
		Infof(logger.LOCALEX_PARSING_PRICE_FILE, symbol)
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return errors.WrapBadRequest(err)
	}

	// Parsing price file
	asset, ctb_err := utils.GetAssetFromSymbol(symbol)
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return ctb_err
	}

	prices[symbol] = crrqueue.New(int64(len(lines)))
	for i := len(lines) - 1; i > 0; i-- {
		mms, ctb_err := parse_mini_market_stats(lines[i], asset, intre, floatre)
		if ctb_err != nil {
			logrus.WithField("comp", "localex").
				Errorf(logger.LOCALEX_ERR_SKIP_PRICE_UPDATE, ctb_err.Error())
			continue
		}

		err = prices[symbol].Put(mms)
		if err != nil {
			logrus.WithField("comp", "localex").Error(err.Error())
			return errors.WrapExchange(err)
		}

	}

	logrus.WithField("comp", "localex").
		Infof(logger.LOCALEX_SYMBOL_PRICE_NUMBER, prices[symbol].Len(), symbol)
	return nil

}

func parse_mini_market_stats(line []string, asset string, intre, floatre *regexp.Regexp) (model.MiniMarketStats, errors.CtbError) {
	matched := intre.Match([]byte(line[0]))
	if !matched {
		err := errors.BadRequest(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "unix_timestamp", line[0])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	timestamp, _ := strconv.ParseInt(line[0], 10, 64)

	matched = floatre.Match([]byte(line[3]))
	if !matched {
		err := errors.BadRequest(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "open_price", line[3])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	open := utils.DecimalFromString(line[3])

	matched = floatre.Match([]byte(line[4]))
	if !matched {
		err := errors.BadRequest(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "high_price", line[4])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	high := utils.DecimalFromString(line[4])

	matched = floatre.Match([]byte(line[5]))
	if !matched {
		err := errors.BadRequest(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "low_price", line[5])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	low := utils.DecimalFromString(line[5])

	matched = floatre.Match([]byte(line[6]))
	if !matched {
		err := errors.BadRequest(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "close_price", line[6])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	close := utils.DecimalFromString(line[6])

	matched = floatre.Match([]byte(line[7]))
	if !matched {
		err := errors.BadRequest(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "base_volume", line[7])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	baseVolume := utils.DecimalFromString(line[7])

	matched = floatre.Match([]byte(line[8]))
	if !matched {
		err := errors.BadRequest(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "quote_volume", line[8])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	quoteVolume := utils.DecimalFromString(line[8])

	return model.MiniMarketStats{
		Event:       "mini_market_stats_update",
		Time:        timestamp,
		Asset:       asset,
		LastPrice:   close,
		HighPrice:   high,
		LowPrice:    low,
		OpenPrice:   open,
		BaseVolume:  baseVolume,
		QuoteVolume: quoteVolume}, nil
}

func localex_can_spot_trade(symbol string) bool {
	return true
}

func localex_get_spot_market_limits(symbol string) (model.SpotMarketLimits, errors.CtbError) {
	return model.SpotMarketLimits{
		MinBase:  utils.DecimalFromString("0.00000001"),
		MinQuote: utils.DecimalFromString("0.00000001"),
		MaxBase:  utils.DecimalFromString("99999999"),
		StepBase: utils.DecimalFromString("0.00000001")}, nil
}

func localex_filter_tradable_assets(bases []string) []string {
	tradables := make([]string, 0, len(bases))

	for _, base := range bases {
		if base != "USDT" {
			tradables = append(tradables, base)
		}
	}
	return tradables
}

func localex_get_assets_value(bases []string) (map[string]model.AssetPrice, errors.CtbError) {
	assetPrices := make(map[string]model.AssetPrice)

	for _, base := range bases {
		symbol, ctb_err := utils.GetSymbolFromAsset(base)
		if ctb_err != nil {
			logrus.WithField("comp", "localex").Error(ctb_err.Error())
			return nil, ctb_err
		}

		values, found := prices[symbol]
		if !found {
			err := errors.Internal(logger.LOCALEX_ERR_UNKNOWN_SYMBOL, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return nil, err
		}

		value, err := values.Peek()
		if err != nil {
			logrus.WithField("comp", "localex").Error(err.Error())
			return nil, errors.WrapExchange(err)
		}
		if value == nil {
			err := errors.Exchange(logger.LOCALEX_ERR_SYMBOL_PRICE, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return nil, err
		}
		assetPrices[base] = model.AssetPrice{
			Asset: base,
			Price: value.(model.MiniMarketStats).LastPrice}
	}
	return assetPrices, nil
}

func localex_get_account() (model.RemoteAccount, errors.CtbError) {
	balances := make([]model.RemoteBalance, 0, raccount.Size())

	for itr := raccount.Iterator(); itr.HasNext(); {
		key, value, ok := itr.Next()
		if !ok {
			err := errors.Internal(logger.LOCALEX_ERR_RACCOUNT_BUILD_FAILURE)
			logrus.Error(err.Error())
			return model.RemoteAccount{}, err
		}

		rbalance := model.RemoteBalance{Asset: key.(string), Amount: value.(decimal.Decimal)}
		balances = append(balances, rbalance)
	}
	return model.RemoteAccount{Balances: balances}, nil
}

func localex_send_spot_market_order(op model.Operation) (model.Operation, errors.CtbError) {
	// Getting price list for symbol
	values, found := prices[op.Base+op.Quote]
	if !found {
		op.Status = model.FAILED
		err := errors.Internal(logger.LOCALEX_ERR_UNKNOWN_SYMBOL, op.Base+op.Quote)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}

	// Getting latest price for symbol
	value, err := values.Peek()
	if err != nil {
		op.Status = model.FAILED
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, errors.WrapExchange(err)
	}
	if value == nil {
		op.Status = model.FAILED
		err := errors.Exchange(logger.LOCALEX_ERR_SYMBOL_PRICE, op.Base+op.Quote)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}
	price := value.(model.MiniMarketStats).LastPrice

	// Computing opposite side amount
	var computedAmt decimal.Decimal
	if op.AmountSide == model.BASE_AMOUNT {
		computedAmt = op.Amount.Mul(price).Round(8)
	} else if op.AmountSide == model.QUOTE_AMOUNT {
		computedAmt = op.Amount.Mul(utils.DecimalFromString("1").Div(price).Round(8)).Round(8)
	} else {
		op.Status = model.FAILED
		err := errors.Internal(logger.LOCALEX_ERR_UNKNOWN_AMOUNT_SIDE, op.AmountSide)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}

	// Getting base and quote available amounts
	value, err = raccount.Get(op.Base)
	if err != nil {
		op.Status = model.FAILED
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, errors.WrapExchange(err)
	}
	if value == nil {
		value = decimal.Zero
	}
	baseAmtAvailable := value.(decimal.Decimal)
	value, err = raccount.Get(op.Quote)
	if err != nil {
		op.Status = model.FAILED
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, errors.WrapExchange(err)
	}
	if value == nil {
		value = decimal.Zero
	}
	quoteAmtAvailable := value.(decimal.Decimal)

	// Executing market order
	op.Timestamp = time.Now().UnixMicro()
	if op.Side == model.SELL && op.AmountSide == model.BASE_AMOUNT {
		baseAmtAvailable = baseAmtAvailable.Sub(op.Amount).Round(8)
		quoteAmtAvailable = quoteAmtAvailable.Add(computedAmt).Round(8)
	} else if op.Side == model.SELL && op.AmountSide == model.QUOTE_AMOUNT {
		baseAmtAvailable = baseAmtAvailable.Sub(computedAmt).Round(8)
		quoteAmtAvailable = quoteAmtAvailable.Add(op.Amount).Round(8)
	} else if op.Side == model.BUY && op.AmountSide == model.BASE_AMOUNT {
		baseAmtAvailable = baseAmtAvailable.Add(op.Amount).Round(8)
		quoteAmtAvailable = quoteAmtAvailable.Sub(computedAmt).Round(8)
	} else if op.Side == model.BUY && op.AmountSide == model.QUOTE_AMOUNT {
		baseAmtAvailable = baseAmtAvailable.Add(computedAmt).Round(8)
		quoteAmtAvailable = quoteAmtAvailable.Sub(op.Amount).Round(8)
	} else {
		op.Status = model.FAILED
		err := errors.Internal(logger.LOCALEX_ERR_UNKNOWN_SIDE, op.Side)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}

	// Checking market order results
	if baseAmtAvailable.LessThan(decimal.Zero) {
		op.Status = model.FAILED
		err := errors.Internal(logger.LOCALEX_ERR_NEGATIVE_BASE_AMT, op.Base, baseAmtAvailable)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}
	if quoteAmtAvailable.LessThan(decimal.Zero) {
		op.Status = model.FAILED
		err := errors.Internal(logger.LOCALEX_ERR_NEGATIVE_QUOTE_AMT, op.Quote, quoteAmtAvailable)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}

	// Storing market order results
	_, err = raccount.Put(op.Base, baseAmtAvailable)
	if err != nil {
		op.Status = model.FAILED
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, errors.WrapExchange(err)
	}
	_, err = raccount.Put(op.Quote, quoteAmtAvailable)
	if err != nil {
		logrus.WithField("comp", "localex").Panic(err.Error())
	}

	return op, nil
}

func localex_mini_markets_stats_serve() errors.CtbError {
	go func() {
		defer func() {
			logrus.WithField("comp", "localex").Infof(logger.LOCALEX_DONE)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()

		for {
			// Termination condition check
			allDisposed := true
			allEmpty := true
			for _, symbolMmss := range prices {
				disposed := symbolMmss.Disposed()
				empty := symbolMmss.Empty()
				allDisposed = allDisposed && disposed
				allEmpty := allEmpty && empty

				if allDisposed || allEmpty {
					return
				}
			}

			// Getting mmss from queues
			mmss := make([]model.MiniMarketStats, 0, raccount.Size())
			for symbol, symbolMmss := range prices {
				if symbolMmss.Disposed() {
					continue
				}
				if symbolMmss.Len() == 0 {
					continue
				}

				value, err := symbolMmss.Peek()
				if err != nil {
					logrus.WithField("comp", "localex").Panic(err.Error())
				}
				if value == nil {
					logrus.WithField("comp", "localex").
						Panicf(logger.LOCALEX_ERR_FAILT_TO_GET_MMS, symbol)
				}

				mmss = append(mmss, value.(model.MiniMarketStats))
			}

			// Serving mmss through the channel
			logrus.WithField("comp", "localex").
				Tracef(logger.BINEX_MMSS_TO_CHANNEL, utils.ToAssets(mmss))
			mmsCh <- mmss
			wait_mms_acks(len(mmss))

			// Removing mmss from queues
			for _, symbolMmss := range prices {
				symbolMmss.Get(1)
			}
		}
	}()
	return nil
}

func wait_mms_acks(size int) {
	var sum int = 0

	for mmsAck := range cllCh {
		sum = sum + mmsAck.Count
		if sum == size {
			break
		}
	}
}

func localex_mini_markets_stats_stop() {
	logrus.WithField("comp", "localex").Info(logger.LOCALEX_PRICE_QUEUES_DEALLOCATION)
	for _, values := range prices {
		if values != nil && !values.Disposed() {
			values.Dispose()
		}
	}
}
