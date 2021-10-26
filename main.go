package main

import (
	"log"

	"github.com/valerioferretti92/trading-bot-demo/internal/binance"
	"github.com/valerioferretti92/trading-bot-demo/internal/execution"
	"github.com/valerioferretti92/trading-bot-demo/internal/repository"
)

func main() {
	defer shutdown()

	exe, err := execution.CreateOrResumeExecution()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("execution: exeId=%s", exe.ExeId)

	ops, err := repository.FindLatestOperations(exe.ExeId, exe.Symbols)
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
	repository.Disconnect()
}
