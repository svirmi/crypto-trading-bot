package binance

import (
	"context"
	"log"
	"time"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
)

var (
	httpClient *binanceapi.Client

	symbols map[string]binanceapi.Symbol
)

func Initialize() {
	// Web socket keep alive set up
	binanceapi.WebsocketKeepalive = false
	binanceapi.WebsocketTimeout = time.Second * 60
	// Building binance http client
	build_binance_clients()
	// Getting binance exchange symbols
	init_exchange_symbols()
}

func build_binance_clients() {
	binanceConfig := config.GetBinanceApiConfig()
	binanceapi.UseTestnet = binanceConfig.UseTestnet
	httpClient = binanceapi.NewClient(binanceConfig.ApiKey, binanceConfig.SecretKey)
}

func init_exchange_symbols() {
	res, err := httpClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	log.Println("registering trading symbols")
	symbols = make(map[string]binanceapi.Symbol)
	for _, symbol := range res.Symbols {
		symbols[symbol.Symbol] = symbol
	}
}
