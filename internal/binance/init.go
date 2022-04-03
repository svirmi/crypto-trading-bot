package binance

import (
	"context"
	"time"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
)

var (
	httpClient *binanceapi.Client
	symbols    map[string]binanceapi.Symbol
)

func Initialize() error {
	// Web socket keep alive set up
	binanceapi.WebsocketKeepalive = false
	binanceapi.WebsocketTimeout = time.Second * 60

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
