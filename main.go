package main

import (
	"log"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/executions"
	"github.com/valerioferretti92/trading-bot-demo/internal/laccount"
	"github.com/valerioferretti92/trading-bot-demo/internal/operations"
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

	_, err = laccount.CreateOrRestore(exe.ExeId, raccount)
	if err != nil {
		log.Fatalf(err.Error())
	}

	ops, err := operations.FindLatestOperations(exe.ExeId, exe.Symbols)
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, op := range ops {
		log.Printf("operation: symbol=%s, type=%s, timestamp=%v",
			op.Base+op.Quote, op.Type, op.Timestamp)
	}
}

func shutdown() {
	binance.Close()
}
