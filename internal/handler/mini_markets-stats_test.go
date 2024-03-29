package handler

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func TestHandleMiniMarketsStats(t *testing.T) {
	logger.Initialize(false, true, true)
	// Saving and restoring status
	old_get_active_exe := get_latest_exe
	old_get_latest_lacc := get_latest_lacc
	old_hanlde := handle_operation
	old_skip := skip_mini_market_stats
	old_get_op := get_operation
	old_can_spot_trade := can_spot_trade
	old_get_spot_market_limits := get_spot_market_limits
	old_get_asset_statuses := get_asset_amounts
	old_store_prices_deferred := store_prices_deferred
	defer func() {
		get_latest_exe = old_get_active_exe
		get_latest_lacc = old_get_latest_lacc
		handle_operation = old_hanlde
		skip_mini_market_stats = old_skip
		get_operation = old_get_op
		can_spot_trade = old_can_spot_trade
		get_spot_market_limits = old_get_spot_market_limits
		get_asset_amounts = old_get_asset_statuses
		store_prices_deferred = old_store_prices_deferred
	}()

	// Mocking dependencies
	get_latest_exe = func() (model.Execution, errors.CtbError) {
		return model.Execution{Status: model.EXE_ACTIVE}, nil
	}

	get_latest_lacc = func(exeId string) (model.ILocalAccount, errors.CtbError) {
		return laccount_test{}, nil
	}

	handled_counter := 0
	handle_operation = func(lacc model.ILocalAccount, op model.Operation) model.ILocalAccount {
		handled_counter++
		time.Sleep(time.Millisecond * 500)
		return lacc
	}

	skipped_counter := 0
	skip_mini_market_stats = func([]model.MiniMarketStats) {
		skipped_counter++
	}

	get_operation = func(model.Execution, model.ILocalAccount, model.MiniMarketStats, model.SpotMarketLimits) (model.Operation, errors.CtbError) {
		op := model.Operation{}
		op.OpId = uuid.NewString()
		return op, nil
	}

	can_spot_trade = func(symbol string) bool {
		return true
	}

	get_spot_market_limits = func(symbol string) (model.SpotMarketLimits, errors.CtbError) {
		return model.SpotMarketLimits{}, nil
	}

	get_asset_amounts = func(model.ILocalAccount) map[string]model.AssetAmount {
		return map[string]model.AssetAmount{"BTC": {"BTC", decimal.Zero}}
	}

	store_prices_deferred = func(mmss []model.MiniMarketStats) {}

	// Producer
	end := make(chan struct{})
	mmsChannel = make(chan []model.MiniMarketStats)
	go func() {
		for i := 0; i < 6; i++ {
			mmsChannel <- get_mini_markets_stats()
			time.Sleep(time.Millisecond * 250)
		}
		close(mmsChannel)
		end <- struct{}{}
	}()

	// Consumer
	handle_mini_markets_stats()
	<-end

	testutils.AssertEq(t, 6, handled_counter+skipped_counter, "mini_market_stats_count")
}

func TestHandleMiniMarketsStats_NonActiveExe(t *testing.T) {
	logger.Initialize(false, true, true)
	// Saving and restoring status
	old_get_active_exe := get_latest_exe
	old_store_prices_deferred := store_prices_deferred
	old_get_op := get_operation
	defer func() {
		get_latest_exe = old_get_active_exe
		store_prices_deferred = old_store_prices_deferred
		get_operation = old_get_op
	}()

	// Mocking dependencies
	get_latest_exe = func() (model.Execution, errors.CtbError) {
		return model.Execution{Status: model.EXE_TERMINATED}, nil
	}

	store_prices_deferred = func(mmss []model.MiniMarketStats) {}

	get_op_counter := 0
	get_operation = func(model.Execution, model.ILocalAccount, model.MiniMarketStats, model.SpotMarketLimits) (model.Operation, errors.CtbError) {
		get_op_counter++
		return model.Operation{}, nil
	}

	// Producer
	end := make(chan struct{})
	mmsChannel = make(chan []model.MiniMarketStats)
	go func() {
		for i := 0; i < 6; i++ {
			mmsChannel <- []model.MiniMarketStats{}
			time.Sleep(time.Millisecond * 50)
		}
		close(mmsChannel)
		end <- struct{}{}
	}()

	// Consumer
	handle_mini_markets_stats()
	<-end

	testutils.AssertEq(t, 0, get_op_counter, "mini_markets_stats_count")
}

func TestHandleMiniMarketsStats_Noop(t *testing.T) {
	logger.Initialize(false, true, true)
	// Saving and restoring status
	old_get_active_exe := get_latest_exe
	old_get_latest_lacc := get_latest_lacc
	old_hanlde := handle_operation
	old_skip := skip_mini_market_stats
	old_get_op := get_operation
	old_can_spot_trade := can_spot_trade
	old_get_spot_market_limits := get_spot_market_limits
	old_get_asset_statuses := get_asset_amounts
	old_store_prices_deferred := store_prices_deferred
	defer func() {
		get_latest_exe = old_get_active_exe
		get_latest_lacc = old_get_latest_lacc
		get_asset_amounts = old_get_asset_statuses
		handle_operation = old_hanlde
		skip_mini_market_stats = old_skip
		get_operation = old_get_op
		can_spot_trade = old_can_spot_trade
		get_spot_market_limits = old_get_spot_market_limits
		store_prices_deferred = old_store_prices_deferred
	}()

	// Mocking dependencies
	get_latest_exe = func() (model.Execution, errors.CtbError) {
		return model.Execution{Status: model.EXE_ACTIVE}, nil
	}

	get_latest_lacc = func(exeId string) (model.ILocalAccount, errors.CtbError) {
		return laccount_test{}, nil
	}

	handled_counter := 0
	handle_operation = func(lacc model.ILocalAccount, op model.Operation) model.ILocalAccount {
		handled_counter++
		return lacc
	}

	skipped_counter := 0
	skip_mini_market_stats = func([]model.MiniMarketStats) {
		skipped_counter++
	}

	get_operation = func(model.Execution, model.ILocalAccount, model.MiniMarketStats, model.SpotMarketLimits) (model.Operation, errors.CtbError) {
		return model.Operation{}, nil
	}

	can_spot_trade = func(symbol string) bool {
		return true
	}

	get_spot_market_limits = func(symbol string) (model.SpotMarketLimits, errors.CtbError) {
		return model.SpotMarketLimits{}, nil
	}

	get_asset_amounts = func(model.ILocalAccount) map[string]model.AssetAmount {
		return map[string]model.AssetAmount{"BTC": {"BTC", decimal.Zero}}
	}

	store_prices_deferred = func(mmss []model.MiniMarketStats) {}

	// Producer
	end := make(chan struct{})
	mmsChannel = make(chan []model.MiniMarketStats)
	go func() {
		for i := 0; i < 6; i++ {
			mmsChannel <- get_mini_markets_stats()
			time.Sleep(time.Millisecond * 50)
		}
		close(mmsChannel)
		end <- struct{}{}
	}()

	// Consumer
	handle_mini_markets_stats()
	<-end

	testutils.AssertEq(t, 0, handled_counter+skipped_counter, "mini_market_stats_count")
}

func TestComputeOpResults_Filled_NoSpread_Buy_BaseAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("0.1")
	price := utils.DecimalFromString("32887.16")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("9.02")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("6934.384")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.FILLED
	exp.Results = model.OpResults{
		ActualPrice: exp.Price,
		BaseDiff:    utils.DecimalFromString("0.1"),
		QuoteDiff:   utils.DecimalFromString("3288.716"),
		Spread:      decimal.Zero}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_Filled_NoSpread_Sell_QuoteAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("250.00")
	price := utils.DecimalFromString("32000.0")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.9121875")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10473.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.FILLED
	exp.Results = model.OpResults{
		ActualPrice: exp.Price,
		BaseDiff:    utils.DecimalFromString("0.0078125"),
		QuoteDiff:   utils.DecimalFromString("250"),
		Spread:      decimal.Zero}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_Filled_PositiveSpread_Buy_BaseAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("0.1")
	price := utils.DecimalFromString("32887.16")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("9.02")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("6923.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.DecimalFromString("33000"),
		BaseDiff:    utils.DecimalFromString("0.1"),
		QuoteDiff:   utils.DecimalFromString("3300"),
		Spread:      utils.DecimalFromString("0.34311263")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_Filled_NegativeSpread_Buy_BaseAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("0.1")
	price := utils.DecimalFromString("32887.16")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("9.02")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("7023.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.DecimalFromString("32000"),
		BaseDiff:    utils.DecimalFromString("0.1"),
		QuoteDiff:   utils.DecimalFromString("3200"),
		Spread:      utils.DecimalFromString("-2.69758775")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_Filled_PositiveSpread_Sell_QuoteAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("250.00")
	price := utils.DecimalFromString("32500.00")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.91242424")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10473.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.DecimalFromString("32999.98944"),
		BaseDiff:    utils.DecimalFromString("0.00757576"),
		QuoteDiff:   utils.DecimalFromString("250"),
		Spread:      utils.DecimalFromString("1.53842905")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_Filled_Negative_Spread_Sell_QuoteAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("250.00")
	price := utils.DecimalFromString("32500.00")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("9")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10473.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.DecimalFromString("3125"),
		BaseDiff:    utils.DecimalFromString("0.08"),
		QuoteDiff:   utils.DecimalFromString("250"),
		Spread:      utils.DecimalFromString("-90.38461538")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_PartiallyFilled_Spread_Buy_BaseAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("0.1")
	price := utils.DecimalFromString("32887.16")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("9.019")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("6923.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.PARTIALLY_FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.DecimalFromString("33333.33333333"),
		BaseDiff:    utils.DecimalFromString("0.099"),
		QuoteDiff:   utils.DecimalFromString("3300"),
		Spread:      utils.DecimalFromString("1.35667943")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_PartiallyFilled_Spread_Sell_QuoteAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("250.00")
	price := utils.DecimalFromString("32500.00")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.912")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("10472.9")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.PARTIALLY_FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.DecimalFromString("31225"),
		BaseDiff:    utils.DecimalFromString("0.008"),
		QuoteDiff:   utils.DecimalFromString("249.8"),
		Spread:      utils.DecimalFromString("-3.92307692")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_NonExecuted_Buy_BaseAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("1")
	price := utils.DecimalFromString("32500.00")
	exp := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)

	b := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("100223.1")}}
	r := get_remote_account(b)

	exp.Status = model.FAILED
	got, err := compute_op_results(r, r, exp)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_NonExecuted_Sell_QuoteAmt(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("10000.5")
	price := utils.DecimalFromString("32500.00")
	exp := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	b := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("100223.1")}}
	r := get_remote_account(b)

	exp.Status = model.FAILED
	got, err := compute_op_results(r, r, exp)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_ZeroBaseDiff_Buy(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("0.75")
	price := utils.DecimalFromString("32500.00")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("100223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("75848.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.PARTIALLY_FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.MaxDecimal(),
		BaseDiff:    decimal.Zero,
		QuoteDiff:   utils.DecimalFromString("24375"),
		Spread:      utils.MaxDecimal()}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_ZeroBaseDiff_Sell(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("0.75")
	price := utils.DecimalFromString("32500.00")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("100223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("124598.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.PARTIALLY_FILLED
	exp.Results = model.OpResults{
		ActualPrice: decimal.Zero,
		BaseDiff:    decimal.Zero,
		QuoteDiff:   utils.DecimalFromString("24375"),
		Spread:      utils.DecimalFromString("-100")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_ZeroQuoteDiff_Sell(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("32500.00")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("100223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.91692308")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("100223.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.PARTIALLY_FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.MaxDecimal(),
		BaseDiff:    utils.DecimalFromString("0.00307692"),
		QuoteDiff:   decimal.Zero,
		Spread:      utils.MaxDecimal()}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_ZeroQuoteDiff_Buy(t *testing.T) {
	logger.Initialize(false, true, true)
	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("32500.00")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.BUY, price)

	b1 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("100223.1")}}
	r1 := get_remote_account(b1)
	b2 := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("8.92307692")},
		{Asset: "ETH", Amount: utils.DecimalFromString("18.92")},
		{Asset: "USDT", Amount: utils.DecimalFromString("100223.1")}}
	r2 := get_remote_account(b2)

	got, err := compute_op_results(r1, r2, op)
	testutils.AssertNil(t, err, "err")

	exp := op
	exp.Status = model.PARTIALLY_FILLED
	exp.Results = model.OpResults{
		ActualPrice: decimal.Zero,
		BaseDiff:    utils.DecimalFromString("0.00307692"),
		QuoteDiff:   decimal.Zero,
		Spread:      utils.DecimalFromString("-100")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func get_mini_markets_stats() []model.MiniMarketStats {
	return []model.MiniMarketStats{
		{Asset: "BTC", LastPrice: utils.DecimalFromString("36781.12")},
		{Asset: "NON_EXISTING", LastPrice: utils.DecimalFromString("0")}}
}

func get_operation_test(amt decimal.Decimal, amtSide model.AmountSide, base, quote string,
	side model.OpSide, price decimal.Decimal) model.Operation {

	return model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      uuid.NewString(),
		Type:       model.AUTO,
		Base:       base,
		Quote:      quote,
		Side:       side,
		Amount:     amt,
		AmountSide: amtSide,
		Price:      price,
		Status:     model.PENDING,
		Timestamp:  time.Now().UnixMicro()}
}

func get_remote_account(balances []model.RemoteBalance) model.RemoteAccount {
	return model.RemoteAccount{
		MakerCommission:  1,
		TakerCommission:  1,
		BuyerCommission:  1,
		SellerCommission: 1,
		Balances:         balances}
}

type laccount_test struct{}

func (a laccount_test) GetAccountId() string {
	return ""
}

func (a laccount_test) GetExeId() string {
	return ""
}

func (a laccount_test) GetStrategyType() model.StrategyType {
	return ""
}

func (a laccount_test) GetTimestamp() int64 {
	return 0
}

func (a laccount_test) Initialize(model.LocalAccountInit) (model.ILocalAccount, errors.CtbError) {
	return nil, nil
}

func (a laccount_test) RegisterTrading(model.Operation) (model.ILocalAccount, errors.CtbError) {
	return nil, nil
}

func (a laccount_test) GetOperation(map[string]string, model.MiniMarketStats, model.SpotMarketLimits) (model.Operation, errors.CtbError) {
	return model.Operation{}, nil
}

func (a laccount_test) GetAssetAmounts() map[string]model.AssetAmount {
	return nil
}

func (a laccount_test) ValidateConfig(map[string]string) errors.CtbError {
	return nil
}
