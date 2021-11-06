package main

import (
	"log"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/executions"
	"github.com/valerioferretti92/trading-bot-demo/internal/laccount"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

func main() {
	defer shutdown()

	raccount, err := binance.GetAccout()
	if err != nil {
		log.Fatalf(err.Error())
	}

	exe, err := executions.CreateOrRestore(raccount)
	if err != nil {
		log.Fatalf(err.Error())
	}

	_, err = laccount.CreateOrRestore(exe.ExeId, raccount, model.FIXED_THRESHOLD_STRATEGY)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func shutdown() {
	binance.Close()
}
