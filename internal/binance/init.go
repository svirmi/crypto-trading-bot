package binance

import (
	"context"
	"fmt"
	"log"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/valerioferretti92/trading-bot-demo/internal/config"
)

var (
	httpClient *binanceapi.Client

	symbols map[string]binanceapi.Symbol
)

func init() {
	// Building binance http client
	buildBinanceClients()
	// Getting binance exchange symbols
	initExchangeSymbols()
}

func buildBinanceClients() {
	binanceConfig := config.AppConfig.BinanceApi
	binanceapi.UseTestnet = binanceConfig.UseTestnet
	httpClient = binanceapi.NewClient(binanceConfig.ApiKey, binanceConfig.SecretKey)
}

func initExchangeSymbols() {
	res, err := httpClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("Registering trading symbols")
	symbols = make(map[string]binanceapi.Symbol)
	for _, symbol := range res.Symbols {
		symbols[symbol.Symbol] = symbol
	}
}
