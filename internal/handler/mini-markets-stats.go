package handler

import (
	"log"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	abool "github.com/tevino/abool/v2"
	"github.com/valerioferretti92/crypto-trading-bot/internal/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
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
			skip_mini_markets_stats(miniMarketsStats)
			continue
		}

		go func(miniMarketsStats []model.MiniMarketStats) {
			defer sentinel.UnSet()
			handle_mini_markets_stats(miniMarketsStats)
		}(miniMarketsStats)
	}
}

var skip_mini_markets_stats = func([]model.MiniMarketStats) {
	logrus.Info(logger.HANDL_SKIP_MMS_UPDATE)
}

var handle_mini_markets_stats = func(miniMarketsStats []model.MiniMarketStats) {
	trading_context_init()

	// If the execution is PAUSED, no action should be applied
	if tcontext.execution.Status != model.EXE_ACTIVE {
		return
	}

	for _, mms := range miniMarketsStats {
		// Getting target operation
		operation, err := tcontext.laccount.GetOperation(mms)
		if err != nil && err.(model.CtbError).IsRecoverable() {
			logrus.Errorf(logger.HANDL_ERR_RECOVERABLE, err.Error())
			continue
		}
		if err != nil {
			logrus.Panicf(logger.HANDL_ERR_UNRECOVERABLE, err.Error())
		}

		// NOOP
		if operation.IsEmpty() {
			continue
		}

		// Getting remote account before operation
		raccount1, err := binance.GetAccout()
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_RECOVERABLE, err.Error())
			continue
		}

		// Sending market order
		operation, err = binance.SendSpotMarketOrder(operation)
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_RECOVERABLE, err.Error())
			continue
		}
		if operation.Status == model.FAILED {
			logrus.Errorf(logger.HANDL_ERR_MKT_ODR_FAILED, operation.OpId)
			continue
		}

		// Getting remote account after operation
		raccount2, err := binance.GetAccout()
		if err != nil {
			logrus.Panicf(logger.HANDL_ERR_UNRECOVERABLE, err.Error())
			continue
		}

		// Computing operation results
		operation = compute_op_results(raccount1, raccount2, operation)

		// Updating local account
		tcontext.laccount, err = tcontext.laccount.RegisterTrading(operation)
		if err != nil {
			log.Fatalf(err.Error())
		}

		// Inserting operation and updating laccount in DB
		operations.Create(operation)
		laccount.Create(tcontext.laccount)
	}
}

func compute_op_results(old, new model.RemoteAccount, op model.Operation) model.Operation {
	var oldBaseBalance, newBaseBalance decimal.Decimal
	var oldQuoteBalance, newQuoteBalance decimal.Decimal

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
	baseDiff := (newBaseBalance.Sub(oldBaseBalance)).Abs().Round(8)
	quoteDiff := (newQuoteBalance.Sub(oldQuoteBalance)).Abs().Round(8)
	if op.AmountSide == model.BASE_AMOUNT && baseDiff.Equals(op.Amount) {
		op.Status = model.FILLED
	} else if op.AmountSide == model.QUOTE_AMOUNT && quoteDiff.Equals(op.Amount) {
		op.Status = model.FILLED
	} else {
		op.Status = model.PARTIALLY_FILLED
	}

	// Building results
	actualPrice := quoteDiff.Div(baseDiff).Round(8)
	results := model.OpResults{
		ActualPrice: actualPrice,
		BaseDiff:    baseDiff,
		QuoteDiff:   quoteDiff,
		Spread:      ((actualPrice.Sub(op.Price)).Div(op.Price)).Mul(decimal.NewFromInt(100)).Round(8)}
	op.Results = results

	logrus.Infof(logger.HANDL_OPERATION_RESULTS,
		op.Results.BaseDiff, op.Results.QuoteDiff, op.Results.ActualPrice, op.Results.Spread, op.Status)
	return op
}

func trading_context_init() {
	if tcontext.execution.IsEmpty() {
		execution, err := executions.GetCurrentlyActive()
		if err != nil {
			logrus.Panicf(logger.HANDL_ERR_UNRECOVERABLE, err.Error())
		}
		tcontext.execution = execution
	}

	if tcontext.laccount == nil {
		laccount, err := laccount.GetLatestByExeId(tcontext.execution.ExeId)
		if err != nil {
			logrus.Panicf(logger.HANDL_ERR_UNRECOVERABLE, err.Error())
		}
		tcontext.laccount = laccount
	}
}
