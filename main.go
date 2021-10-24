package main

import (
	"log"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/execution"
	"github.com/valerioferretti92/trading-bot-demo/internal/repository"
)

func main() {
	defer shutdown()

	_, err := execution.CreateOrResumeExecution()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func shutdown() {
	binance.Close()
	repository.Disconnect()
}
