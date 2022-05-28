package local

import (
	"encoding/csv"
	"fmt"
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
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

var (
	raccount   *crrmap.ConcurrentMap
	prices     map[string]*crrqueue.Queue
	priceDelay decimal.Decimal
	mmsCh      chan []model.MiniMarketStats
)

type local_exchange struct{}

func GetExchange() model.IExchange {
	return local_exchange{}
}

func (be local_exchange) Initialize(mmsChannel chan []model.MiniMarketStats) error {
	// Decoding config
	localExchangeConfig := struct {
		InitialBalances map[string]string
		PriceFilepaths  map[string]string
		PriceDelay      string
	}{}
	err := mapstructure.Decode(config.GetExchangeConfig(), &localExchangeConfig)
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return err
	}

	// The wallet must be populated with assets, not symbols
	// For each asset in the wallet, but USDT, a price must be specified
	for asset := range localExchangeConfig.InitialBalances {
		if asset == "USDT" {
			continue
		}

		if strings.HasSuffix(asset, "USDT") {
			err := fmt.Errorf(logger.LOCALEX_ERR_INVALID_ASSET, asset)
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}

		symbol := utils.GetSymbolFromAsset(asset)
		if _, found := localExchangeConfig.PriceFilepaths[symbol]; !found {
			err := fmt.Errorf(logger.LOCALEX_ERR_PRICES_NOT_PROVIDED, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}
	}

	// Prices are to be provided per symbol, the form XXXUSDT
	// USDTUSDT is not a valid symbol
	// Skip prices whose corresponding asset is not in the wallet
	symbolre, err := regexp.Compile(`^.+USDT$`)
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return err
	}
	for symbol := range localExchangeConfig.PriceFilepaths {
		matched := symbolre.Match([]byte(symbol))
		if !matched {
			err := fmt.Errorf(logger.LOCALEX_ERR_INVALID_SYMBOL, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}

		base := strings.TrimSuffix(symbol, "USDT")
		if base == "USDT" {
			err := fmt.Errorf(logger.LOCALEX_ERR_INVALID_SYMBOL, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}

		asset := utils.GetAssetFromSymbol(symbol)
		if _, found := localExchangeConfig.InitialBalances[asset]; !found {
			logrus.WithField("comp", "localex").
				Warnf(logger.LOCALEX_SKIP_SYMBOL_PRICES, asset, symbol)
			delete(localExchangeConfig.PriceFilepaths, symbol)
		}
	}

	// Initializing mms channel
	mmsCh = mmsChannel

	// Initializing price delay
	minPriceDelay := utils.DecimalFromString("50")
	priceDelay = utils.DecimalFromString(localExchangeConfig.PriceDelay)
	if priceDelay.LessThan(minPriceDelay) {
		err := fmt.Errorf(logger.LOCALEX_ERR_PRICE_DELAY_TOO_SMALL, priceDelay, minPriceDelay)
		logrus.Error(err.Error())
		return err
	}

	// Parsing account balances
	logrus.Infof(logger.LOCALEX_INIT_RACCOUNT, len(localExchangeConfig.InitialBalances))
	raccount = crrmap.NewConcurrentMap()
	for key, value := range localExchangeConfig.InitialBalances {
		_, err := raccount.Put(key, utils.DecimalFromString(value))
		if err != nil {
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
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

func parse_prices_file(symbol, priceFilepath string) error {
	file, err := os.Open(priceFilepath)
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return err
	}
	defer file.Close()

	// Compiling regexps
	intre, err := regexp.Compile(`^[1-9][0-9]*$`)
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return err
	}
	floatre, err := regexp.Compile(`^(([1-9][0-9]+)|[0-9])(\.[0-9]+)?$`)
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return err
	}

	// Reading price file
	logrus.WithField("comp", "localex").
		Infof(logger.LOCALEX_PARSING_PRICE_FILE, symbol)
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		logrus.WithField("comp", "localex").Error(err.Error())
		return err
	}

	// Parsing price file
	asset := utils.GetAssetFromSymbol(symbol)
	prices[symbol] = crrqueue.New(int64(len(lines)))
	for i := len(lines) - 1; i > 0; i-- {
		mms, err := parse_mini_market_stats(lines[i], asset, intre, floatre)
		if err != nil {
			logrus.WithField("comp", "localex").
				Errorf(logger.LOCALEX_ERR_SKIP_PRICE_UPDATE, err.Error())
			continue
		}

		err = prices[symbol].Put(mms)
		if err != nil {
			logrus.WithField("comp", "localex").Error(err.Error())
			return err
		}

	}

	logrus.WithField("comp", "localex").
		Infof(logger.LOCALEX_SYMBOL_PRICE_NUMBER, prices[symbol].Len(), symbol)
	return nil

}

func parse_mini_market_stats(line []string, asset string, intre, floatre *regexp.Regexp) (model.MiniMarketStats, error) {
	matched := intre.Match([]byte(line[0]))
	if !matched {
		err := fmt.Errorf(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "unix_timestamp", line[0])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	timestamp, _ := strconv.ParseInt(line[0], 10, 64)

	matched = floatre.Match([]byte(line[3]))
	if !matched {
		err := fmt.Errorf(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "open_price", line[3])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	open := utils.DecimalFromString(line[3])

	matched = floatre.Match([]byte(line[4]))
	if !matched {
		err := fmt.Errorf(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "high_price", line[4])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	high := utils.DecimalFromString(line[4])

	matched = floatre.Match([]byte(line[5]))
	if !matched {
		err := fmt.Errorf(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "low_price", line[5])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	low := utils.DecimalFromString(line[5])

	matched = floatre.Match([]byte(line[6]))
	if !matched {
		err := fmt.Errorf(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "close_price", line[6])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	close := utils.DecimalFromString(line[6])

	matched = floatre.Match([]byte(line[7]))
	if !matched {
		err := fmt.Errorf(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "base_volume", line[7])
		logrus.WithField("comp", "localex").Error(err.Error())
		return model.MiniMarketStats{}, err
	}
	baseVolume := utils.DecimalFromString(line[7])

	matched = floatre.Match([]byte(line[8]))
	if !matched {
		err := fmt.Errorf(logger.LOCALEX_ERR_FIELD_BAD_FORMAT, "quote_volume", line[8])
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

func (be local_exchange) GetAssetsValue(bases []string) (map[string]model.AssetPrice, error) {
	assetPrices := make(map[string]model.AssetPrice)

	for _, base := range bases {
		symbol := utils.GetSymbolFromAsset(base)
		values, found := prices[symbol]
		if !found {
			err := fmt.Errorf(logger.LOCALEX_ERR_UNKNOWN_SYMBOL, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return nil, err
		}

		value, err := values.Peek()
		if err != nil {
			logrus.WithField("comp", "localex").Error(err.Error())
			return nil, err
		}
		if value == nil {
			err := fmt.Errorf(logger.LOCALEX_ERR_SYMBOL_PRICE, symbol)
			logrus.WithField("comp", "localex").Error(err.Error())
			return nil, err
		}
		assetPrices[base] = model.AssetPrice{
			Asset: base,
			Price: value.(model.MiniMarketStats).LastPrice}
	}
	return assetPrices, nil
}

func (be local_exchange) SendSpotMarketOrder(op model.Operation) (model.Operation, error) {
	// Getting price list for symbol
	values, found := prices[op.Base+op.Quote]
	if !found {
		op.Status = model.FAILED
		err := fmt.Errorf(logger.LOCALEX_ERR_UNKNOWN_SYMBOL, op.Base+op.Quote)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}

	// Getting latest price for symbol
	value, err := values.Peek()
	if err != nil {
		op.Status = model.FAILED
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}
	if value == nil {
		op.Status = model.FAILED
		err := fmt.Errorf(logger.LOCALEX_ERR_SYMBOL_PRICE, op.Base+op.Quote)
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
		err := fmt.Errorf(logger.LOCALEX_ERR_UNKNOWN_AMOUNT_SIDE, op.AmountSide)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}

	// Getting base and quote available amounts
	value, err = raccount.Get(op.Base)
	if err != nil {
		op.Status = model.FAILED
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}
	if value == nil {
		value = decimal.Zero
	}
	baseAmtAvailable := value.(decimal.Decimal)
	value, err = raccount.Get(op.Quote)
	if err != nil {
		op.Status = model.FAILED
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
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
		err := fmt.Errorf(logger.LOCALEX_ERR_UNKNOWN_SIDE, op.Side)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}

	// Checking market order results
	if baseAmtAvailable.LessThan(decimal.Zero) {
		op.Status = model.FAILED
		err := fmt.Errorf(logger.LOCALEX_ERR_NEGATIVE_BASE_AMT, op.Base, baseAmtAvailable)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}
	if quoteAmtAvailable.LessThan(decimal.Zero) {
		op.Status = model.FAILED
		err := fmt.Errorf(logger.LOCALEX_ERR_NEGATIVE_QUOTE_AMT, op.Quote, quoteAmtAvailable)
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}

	// Storing market order results
	_, err = raccount.Put(op.Base, baseAmtAvailable)
	if err != nil {
		op.Status = model.FAILED
		logrus.WithField("comp", "localex").Error(err.Error())
		return op, err
	}
	_, err = raccount.Put(op.Quote, quoteAmtAvailable)
	if err != nil {
		logrus.WithField("comp", "localex").Panic(err.Error())
	}

	return op, nil
}

func (be local_exchange) MiniMarketsStatsServe(assets []string) error {
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
			mmss := make([]model.MiniMarketStats, 0, len(assets))
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
			mmsCh <- mmss
			delay := priceDelay.Abs().IntPart()
			time.Sleep(time.Duration(delay) * time.Millisecond)

			// Removing mmss from queues
			for _, symbolMmss := range prices {
				symbolMmss.Get(1)
			}
		}
	}()

	return nil
}

func (be local_exchange) MiniMarketsStatsStop() {
	logrus.WithField("comp", "localex").Info(logger.LOCALEX_PRICE_QUEUES_DEALLOCATION)
	for _, values := range prices {
		if values != nil && !values.Disposed() {
			values.Dispose()
		}
	}
}

func (be local_exchange) GetAccout() (model.RemoteAccount, error) {
	balances := make([]model.RemoteBalance, 0, raccount.Size())

	for itr := raccount.Iterator(); itr.HasNext(); {
		key, value, ok := itr.Next()
		if !ok {
			err := fmt.Errorf(logger.LOCALEX_ERR_RACCOUNT_BUILD_FAILURE)
			logrus.Error(err.Error())
			return model.RemoteAccount{}, err
		}

		rbalance := model.RemoteBalance{Asset: key.(string), Amount: value.(decimal.Decimal)}
		balances = append(balances, rbalance)
	}
	return model.RemoteAccount{Balances: balances}, nil
}

func (be local_exchange) CanSpotTrade(symbol string) bool {
	return true
}

func (be local_exchange) GetSpotMarketLimits(symbol string) (model.SpotMarketLimits, error) {
	return model.SpotMarketLimits{
		MinBase:  utils.DecimalFromString("0.00000001"),
		MinQuote: utils.DecimalFromString("0.00000001"),
		MaxBase:  utils.DecimalFromString("99999999"),
		StepBase: utils.DecimalFromString("0.00000001")}, nil
}

func (be local_exchange) FilterTradableAssets(bases []string) []string {
	tradables := make([]string, 0, len(bases))

	for _, base := range bases {
		if base != "USDT" {
			tradables = append(tradables, base)
		}
	}
	return tradables
}
