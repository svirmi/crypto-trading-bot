package handler

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func TestReadingMarketStatsCh(t *testing.T) {
	// Saving and restoring status
	old_hanlde := handle_mini_markets_stats
	old_skip := skip_mini_markets_stats
	defer func() {
		handle_mini_markets_stats = old_hanlde
		skip_mini_markets_stats = old_skip
	}()

	// Mocking dependencies
	handled := 0
	handle_mini_markets_stats = func([]model.MiniMarketStats) {
		handled++
		time.Sleep(time.Millisecond * 500)
	}
	skipped := 0
	skip_mini_markets_stats = func([]model.MiniMarketStats) {
		skipped++
	}

	end := make(chan struct{})

	// Producer
	scontext.mms = make(chan []model.MiniMarketStats)
	go func() {
		for i := 0; i < 6; i++ {
			scontext.mms <- []model.MiniMarketStats{}
			time.Sleep(time.Millisecond * 250)
		}
		close(scontext.mms)
		end <- struct{}{}
	}()

	// Consumer
	read_mini_markets_stats_ch()
	<-end

	testutils.AssertEq(t, 6, handled+skipped, "mini_market_stats_read_ch")
}

func TestComputeOpResults_Filled_NoSpread_Buy_BaseAmt(t *testing.T) {
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

	got := compute_op_results(r1, r2, op)
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

	got := compute_op_results(r1, r2, op)
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

	got := compute_op_results(r1, r2, op)
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

	got := compute_op_results(r1, r2, op)
	exp := op
	exp.Status = model.FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.DecimalFromString("32000"),
		BaseDiff:    utils.DecimalFromString("0.1"),
		QuoteDiff:   utils.DecimalFromString("3200"),
		Spread:      utils.DecimalFromString("-2.69758775")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_Filled_Spread_Sell_QuoteAmt(t *testing.T) {
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

	got := compute_op_results(r1, r2, op)
	exp := op
	exp.Status = model.FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.DecimalFromString("32999.98944"),
		BaseDiff:    utils.DecimalFromString("0.00757576"),
		QuoteDiff:   utils.DecimalFromString("250"),
		Spread:      utils.DecimalFromString("1.53842905")}

	testutils.AssertEq(t, exp, got, "operation_results")
}

func TestComputeOpResults_PartiallyFilled_Spread_Buy_BaseAmt(t *testing.T) {
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

	got := compute_op_results(r1, r2, op)
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

	got := compute_op_results(r1, r2, op)
	exp := op
	exp.Status = model.PARTIALLY_FILLED
	exp.Results = model.OpResults{
		ActualPrice: utils.DecimalFromString("31225"),
		BaseDiff:    utils.DecimalFromString("0.008"),
		QuoteDiff:   utils.DecimalFromString("249.8"),
		Spread:      utils.DecimalFromString("-3.92307692")}

	testutils.AssertEq(t, exp, got, "operation_results")
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
