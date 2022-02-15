package handler

import (
	"log"
	"math"

	abool "github.com/tevino/abool/v2"
	"github.com/valerioferretti92/crypto-trading-bot/internal/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/operations"
)

type trading_context struct {
	laccount  model.ILocalAccount
	execution model.Execution
}

var tcontext trading_context

type stream_context struct {
	mms chan []model.MiniMarketStats
}

var scontext stream_context

func InitTradingContext(laccount model.ILocalAccount, execution model.Execution) {
	tcontext.laccount = laccount
	tcontext.execution = execution
}

func InvalidateTradingContext() {
	tcontext.execution = model.Execution{}
	tcontext.laccount = nil
}

func InitMmsChannel(mms chan []model.MiniMarketStats) {
	scontext.mms = mms
}

func HandleMiniMarketsStats() {
	go read_mini_markets_stats_ch()
}

func read_mini_markets_stats_ch() {
	sentinel := abool.New()

	for miniMarketsStats := range scontext.mms {
		ok := sentinel.SetToIf(false, true)
		if !ok {
			log.Printf("skipping mini markets stats update...")
			continue
		}

		go func(miniMarketsStats []model.MiniMarketStats) {
			defer sentinel.UnSet()
			handle_mini_markets_stats(miniMarketsStats)
		}(miniMarketsStats)
	}
}

func handle_mini_markets_stats(miniMarketsStats []model.MiniMarketStats) {
	trading_context_init()

	// If the execution is PAUSED, no action should be applied
	if tcontext.execution.Status != model.EXE_ACTIVE {
		return
	}

	for _, mms := range miniMarketsStats {
		// Getting target operation
		operation, err := tcontext.laccount.GetOperation(mms)
		if err != nil {
			log.Printf("%s", err.Error())
			continue
		}

		// NOOP
		if operation.IsEmpty() {
			continue
		}

		// Getting remote account before operation
		raccount1, err := binance.GetAccout()
		if err != nil {
			log.Fatalf("failed to get remote account")
			continue
		}

		// Sending market order
		operation, err = binance.SendMarketOrder(operation)
		if err != nil {
			log.Printf("%s", err.Error())
			continue
		}
		if operation.Status == model.FAILED {
			log.Printf("failed to place market order %v", operation)
			continue
		}

		// Getting remote account after operation
		raccount2, err := binance.GetAccout()
		if err != nil {
			log.Fatalf("failed to get remote account")
		}

		// Computing operation results
		operation = compute_op_results(raccount1, raccount2, operation)

		// Updating local account
		tcontext.laccount, err = tcontext.laccount.RegisterTrading(operation)
		if err != nil {
			log.Fatalf(err.Error())
		}

		// Inserting operation in DB
		operations.Create(operation)
	}
}

func compute_op_results(old, new model.RemoteAccount, op model.Operation) model.Operation {
	var oldBaseBalance, newBaseBalance float32
	var oldQuoteBalance, newQuoteBalance float32

	// Getting old abase and quote balances
	for _, balance := range old.Balances {
		if balance.Asset == op.Base {
			oldBaseBalance = balance.Amount
		}
		if balance.Asset == op.Quote {
			oldQuoteBalance = balance.Amount
		}
	}

	// Getting new base and quote balances
	for _, balance := range new.Balances {
		if balance.Asset == op.Base {
			newBaseBalance = balance.Amount
		}
		if balance.Asset == op.Quote {
			newQuoteBalance = balance.Amount
		}
	}

	// Setting operation results
	baseDiff := float32(math.Abs(float64(newBaseBalance) - float64(oldBaseBalance)))
	quoteDiff := float32(math.Abs(float64(newQuoteBalance) - float64(oldQuoteBalance)))
	if op.AmountSide == model.BASE_AMOUNT && baseDiff == op.Amount {
		op.Status = model.FILLED
	} else if op.AmountSide == model.QUOTE_AMOUNT && quoteDiff == op.Amount {
		op.Status = model.FILLED
	} else {
		op.Status = model.PARTIALLY_FILLED
	}

	// Building results
	actualPrice := quoteDiff / baseDiff
	results := model.OpResults{
		ActualPrice: actualPrice,
		BaseAmount:  baseDiff,
		QuoteAmount: quoteDiff,
		Spread:      ((actualPrice - op.Price) / op.Price) * 100}
	op.Results = results
	log_operation_results(op, baseDiff, quoteDiff)

	return op
}

func log_operation_results(op model.Operation, baseDiff, quoteDiff float32) {
	log.Printf("operation %s: %s", op.OpId, op.Status)
	log.Printf("side: %s", op.Side)
	log.Printf("amount side: %s", op.AmountSide)
	log.Printf("amount: %f", op.Amount)
	log.Printf("baseDiff: %f", baseDiff)
	log.Printf("quoteDiff: %f", quoteDiff)
	log.Printf("price: %f", op.Price)
	log.Printf("actualPrice: %f", op.Results.ActualPrice)
	log.Printf("spread: %f", op.Results.Spread)
}

func trading_context_init() {
	if tcontext.execution.IsEmpty() {
		execution, err := executions.GetCurrentlyActive()
		if err != nil {
			log.Fatalf("failed to retrieve active execution")
		}
		tcontext.execution = execution
	}

	if tcontext.laccount == nil {
		laccount, err := laccount.GetLatestByExeId(tcontext.execution.ExeId)
		if err != nil {
			log.Fatalf("failed to retrieve local account")
		}
		tcontext.laccount = laccount
	}
}
