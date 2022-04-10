package handler

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	abool "github.com/tevino/abool/v2"
	"github.com/valerioferretti92/crypto-trading-bot/internal/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/operations"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
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
	trading_context_init()
	go handle_mini_markets_stats()
}

var handle_mini_markets_stats = func() {
	sentinel := abool.New()
	for miniMarketsStats := range scontext.mms {
		// If the execution is not ACTIVE, no action should be applied
		if tcontext.execution.Status != model.EXE_ACTIVE {
			continue
		}

		// Get operations from mini markets stats
		operations := get_operations(miniMarketsStats)
		if len(operations) == 0 {
			continue
		}

		// Set sentinel
		ok := sentinel.SetToIf(false, true)
		if !ok {
			skip_mini_markets_stats(miniMarketsStats)
			continue
		}

		// Handle operations and defer sentinel reset
		go func(operations []model.Operation) {
			defer sentinel.UnSet()
			handle_operations(operations)
		}(operations)
	}
}

var get_operations = func(miniMarketsStats []model.MiniMarketStats) []model.Operation {
	operations := make([]model.Operation, 0)
	for _, miniMarketStats := range miniMarketsStats {
		// Getting operation from mini market stats
		operation, err := tcontext.laccount.GetOperation(miniMarketStats)
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, miniMarketStats.Asset, err.Error())
			continue
		}

		// NOOP
		if operation.IsEmpty() {
			continue
		}
		operations = append(operations, operation)
	}
	return operations
}

var skip_mini_markets_stats = func([]model.MiniMarketStats) {
	logrus.Info(logger.HANDL_SKIP_MMS_UPDATE)
}

var handle_operations = func(ops []model.Operation) {
	for _, op := range ops {
		// Price equals to zero
		if op.Amount.Equals(decimal.Zero) {
			logrus.Errorf(logger.HANDL_ERR_ZERO_REQUESTED_AMOUNT)
			continue
		}
		if op.Price.Equals(decimal.Zero) {
			logrus.Errorf(logger.HANDL_ERR_ZERO_EXP_PRICE)
			continue
		}

		// Getting remote account before operation
		raccount1, err := binance.GetAccout()
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
			continue
		}

		// Sending market order
		op, err = binance.SendSpotMarketOrder(op)
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
			continue
		}

		// Getting remote account after operation
		raccount2, err := binance.GetAccout()
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
			continue
		}

		// Computing operation results
		op.FromId = tcontext.laccount.GetAccountId()
		op, err = compute_op_results(raccount1, raccount2, op)
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
			continue
		}

		// Updating local account
		tcontext.laccount, err = tcontext.laccount.RegisterTrading(op)
		if err != nil {
			logrus.Errorf(logger.HANDL_ERR_SKIP_MMS_UPDATE, op.Base, err.Error())
			continue
		}

		// Inserting operation and updating laccount in DB
		op.ToId = tcontext.laccount.GetAccountId()
		err = operations.Create(op)
		if err != nil {
			logrus.Panicf(err.Error())
		}

		err = laccount.Create(tcontext.laccount)
		if err != nil {
			logrus.Panicf(err.Error())
		}
	}
}

func compute_op_results(old, new model.RemoteAccount, op model.Operation) (model.Operation, error) {
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

	// Checking base diff and quote diff
	baseDiff := (newBaseBalance.Sub(oldBaseBalance)).Abs().Round(8)
	quoteDiff := (newQuoteBalance.Sub(oldQuoteBalance)).Abs().Round(8)
	if baseDiff.Equals(decimal.Zero) && quoteDiff.Equals(decimal.Zero) {
		err := fmt.Errorf(logger.HANDL_ERR_ZERO_BASE_QUOTE_DIFF)
		logrus.Error(err.Error())
		op.Status = model.FAILED
		return op, err
	} else if baseDiff.Equals(decimal.Zero) {
		logrus.Warnf(logger.HANDL_ZERO_BASE_DIFF, op.OpId)
	} else if quoteDiff.Equals(decimal.Zero) {
		logrus.Warnf(logger.HANDL_ZERO_QUOTE_DIFF, op.OpId)
	}

	// Computing status
	if baseDiff.Equals(decimal.Zero) || quoteDiff.Equals(decimal.Zero) {
		op.Status = model.PARTIALLY_FILLED
	} else if op.AmountSide == model.BASE_AMOUNT && baseDiff.Equals(op.Amount) {
		op.Status = model.FILLED
	} else if op.AmountSide == model.QUOTE_AMOUNT && quoteDiff.Equals(op.Amount) {
		op.Status = model.FILLED
	} else {
		op.Status = model.PARTIALLY_FILLED
	}

	// Computing actual price and spread
	var actualPrice decimal.Decimal
	var spread decimal.Decimal
	if baseDiff.Equals(decimal.Zero) && op.Side == model.BUY {
		actualPrice = utils.MaxDecimal()
		spread = utils.MaxDecimal()
	} else if baseDiff.Equals(decimal.Zero) && op.Side == model.SELL {
		actualPrice = decimal.Zero
		spread = utils.DecimalFromString("-100")
	} else if quoteDiff.Equals(decimal.Zero) && op.Side == model.BUY {
		actualPrice = decimal.Zero
		spread = utils.DecimalFromString("-100")
	} else if quoteDiff.Equals(decimal.Zero) && op.Side == model.SELL {
		actualPrice = utils.MaxDecimal()
		spread = utils.MaxDecimal()
	} else {
		actualPrice = quoteDiff.Div(baseDiff).Round(8)
		spread = ((actualPrice.Sub(op.Price)).
			Div(op.Price)).
			Mul(decimal.NewFromInt(100)).
			Round(8)
	}

	// Setting results
	results := model.OpResults{
		ActualPrice: actualPrice,
		BaseDiff:    baseDiff,
		QuoteDiff:   quoteDiff,
		Spread:      spread}
	op.Results = results

	logrus.Infof(logger.HANDL_OPERATION_RESULTS,
		op.Results.BaseDiff, op.Results.QuoteDiff, op.Results.ActualPrice, op.Results.Spread, op.Status)
	return op, nil
}

func trading_context_init() {
	if tcontext.execution.IsEmpty() {
		execution, err := executions.GetCurrentlyActive()
		if err != nil {
			logrus.Panicf(err.Error())
		}
		tcontext.execution = execution
	}

	if tcontext.laccount == nil {
		laccount, err := laccount.GetLatestByExeId(tcontext.execution.ExeId)
		if err != nil {
			logrus.Panicf(err.Error())
		}
		tcontext.laccount = laccount
	}
}
