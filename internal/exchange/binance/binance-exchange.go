package binance

import (
	"context"
	"time"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
)

var (
	httpClient         *binanceapi.Client
	symbols            map[string]binanceapi.Symbol
	mms_done, mms_stop chan struct{}
	mms                chan []model.MiniMarketStats
)

type binance_exchange struct{}

func GetExchange() model.IExchange {
	return binance_exchange{}
}

func (be binance_exchange) Initialize(mmsChannel chan []model.MiniMarketStats) error {
	// Web socket keep alive set up
	binanceapi.WebsocketKeepalive = true
	binanceapi.WebsocketTimeout = time.Second * 60
	mms = mmsChannel

	// Building binance http client
	binanceConfig := config.GetBinanceApiConfig()
	binanceapi.UseTestnet = binanceConfig.UseTestnet
	httpClient = binanceapi.NewClient(binanceConfig.ApiKey, binanceConfig.SecretKey)

	// Init exchange symbols
	res, err := httpClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return err
	}

	logrus.WithField("comp", "binance").Info(logger.BINANCE_REGISTERING_SYMBOLS)
	symbols = make(map[string]binanceapi.Symbol)
	for _, symbol := range res.Symbols {
		symbols[symbol.Symbol] = symbol
	}
	return nil
}

func (be binance_exchange) CanSpotTrade(symbol string) bool {
	return can_spot_trade(symbol)
}

func (be binance_exchange) GetSpotMarketLimits(symbol string) (model.SpotMarketLimits, error) {
	return get_spot_market_limits(symbol)
}

func (be binance_exchange) FilterTradableAssets(bases []string) []string {
	return filter_tradable_assets(bases)
}

func (be binance_exchange) GetAssetsValue(bases []string) (map[string]model.AssetPrice, error) {
	return get_assets_value(bases)
}

func (be binance_exchange) GetAccout() (model.RemoteAccount, error) {
	return get_account()
}

func (be binance_exchange) SendSpotMarketOrder(op model.Operation) (model.Operation, error) {
	return send_spot_market_order(op)
}

func (be binance_exchange) MiniMarketsStatsServe(assets []string) error {
	return mini_markets_stats_serve(assets)
}

func (be binance_exchange) MiniMarketsStatsStop() {
	mini_markets_stats_stop()
}
