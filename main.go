package main

import (
	"log"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/executions"
	"github.com/valerioferretti92/trading-bot-demo/internal/operations"
)

func main() {
	defer shutdown()

	exe, err := executions.CreateOrRestoreExecution()
	if err != nil {
		log.Fatalf(err.Error())
	}

	ops, err := operations.FindLatestOperations(exe.ExeId, exe.Symbols)
	if err != nil {
		log.Fatal(err.Error())
	}
	for i := range ops {
		log.Printf("operation: symbol=%s, type=%s, timestamp=%v",
			ops[i].Symbol, ops[i].Type, ops[i].Timestamp)
	}
}

func shutdown() {
	binance.Close()
}
