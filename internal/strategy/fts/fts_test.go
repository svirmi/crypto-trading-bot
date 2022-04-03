package fts

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/binance"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

/********************** Testing Initialize() *************************/

func TestInitialize(t *testing.T) {
	old := binance.CanSpotTrade
	defer func() {
		binance.CanSpotTrade = old
	}()

	binance.CanSpotTrade = func(symbol string) bool {
		return true
	}

	laccountInit := get_laccount_init_test()
	laccountInit.RAccount.Balances = append(laccountInit.RAccount.Balances, model.RemoteBalance{
		Asset:  "SHIBA",
		Amount: decimal.Zero})
	got, err := LocalAccountFTS{}.Initialize(get_laccount_init_test())
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_last_buy_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()

	testutils.AssertEq(t, exp, got, "fts_laccount")
}

/********************** Testing RegisterTrading() *************************/

func TestRegisterTrading_BaseAmt_BuySide(t *testing.T) {
	laccount := get_laccount_last_buy_test()
	assetStatus := laccount.Assets["BTC"]
	assetStatus.Usdt = utils.DecimalFromString("50")
	laccount.Assets["BTC"] = assetStatus

	amt := utils.DecimalFromString("0.1")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	got, err := laccount.RegisterTrading(op)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_last_buy_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	assetStatus = exp.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("11.44")
	assetStatus.Usdt = decimal.Zero
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus

	testutils.AssertEq(t, exp, got, "fts_laccount")
}

func TestRegisterTrading_BaseAmt_SellSide(t *testing.T) {
	laccount := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("11.34")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId

	got, err := laccount.RegisterTrading(op)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_last_buy_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	assetStatus := exp.Assets["BTC"]
	assetStatus.Amount = decimal.Zero
	assetStatus.Usdt = utils.DecimalFromString("5670")
	assetStatus.LastOperationType = OP_SELL_FTS
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Ignored["USDT"] = utils.DecimalFromString("155.67")

	testutils.AssertEq(t, exp, got, "fts_laccount")
}

func TestRegisterTrading_QuoteAmt_BuySide(t *testing.T) {
	laccount := get_laccount_last_buy_test()
	assetStatus := laccount.Assets["BTC"]
	assetStatus.Usdt = utils.DecimalFromString("110.18")
	laccount.Assets["BTC"] = assetStatus

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	got, err := laccount.RegisterTrading(op)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_last_buy_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	assetStatus = exp.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("11.55134")
	assetStatus.Usdt = utils.DecimalFromString("4.51")
	assetStatus.LastOperationType = OP_BUY_FTS
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Ignored["USDT"] = utils.DecimalFromString("155.67")

	testutils.AssertEq(t, exp, got, "fts_laccount")
}

func TestRegisterTrading_QuoteAmt_SellSide(t *testing.T) {
	laccount := get_laccount_last_sell_test()
	assetStatus := laccount.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("1.18")
	laccount.Assets["BTC"] = assetStatus

	amt := utils.DecimalFromString("10000")
	price := utils.DecimalFromString("43125.2")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId

	got, err := laccount.RegisterTrading(op)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_last_sell_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	assetStatus = exp.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("0.94811702")
	assetStatus.Usdt = utils.DecimalFromString("34519.999")
	assetStatus.LastOperationType = OP_SELL_FTS
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Ignored["USDT"] = utils.DecimalFromString("155.67")

	testutils.AssertEq(t, exp, got, "fts_laccount")
}

func TestRegisterTrading_WrongExeId(t *testing.T) {
	exp := get_laccount_last_buy_test()
	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	testutils.AssertPanic(t, func() {
		exp.RegisterTrading(op)
	})
}

func TestRegisterTrading_OpFailed(t *testing.T) {
	exp := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = exp.ExeId
	op.Status = model.FAILED

	got, err := exp.RegisterTrading(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "fts_laccount")
}

func TestRegisterTrading_BadQuoteCurrency(t *testing.T) {
	exp := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = exp.ExeId
	op.Quote = "ETH"

	got, err := exp.RegisterTrading(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "fts_laccount")
}

func TestRegisterTrading_AssetNotFound(t *testing.T) {
	exp := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = exp.ExeId
	op.Base = "CRO"

	got, err := exp.RegisterTrading(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "fts_laccount")
}

func TestRegisterTrading_NegativeBalanceBase(t *testing.T) {
	exp := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("1923789.12")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = exp.ExeId

	testutils.AssertPanic(t, func() {
		exp.RegisterTrading(op)
	})
}

func TestRegisterTrading_NegativeBalanceQuote(t *testing.T) {
	exp := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("1923789.12")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = exp.ExeId

	testutils.AssertPanic(t, func() {
		exp.RegisterTrading(op)
	})
}

/********************** Testing GetOperation() *************************/

func TestGetOperation_AssetNotFound(t *testing.T) {
	laccount := get_laccount_last_buy_test()
	mms := get_mms("CRO", utils.DecimalFromString("0.55"))

	op, err := laccount.GetOperation(mms)

	testutils.AssertTrue(t, op.IsEmpty(), "operation")
	testutils.AssertNotNil(t, err, "err")
}

func TestGetOperation_Noop(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_buy_test()
	mms := get_mms("BTC", utils.DecimalFromString("39560.1"))

	got, err := laccount.GetOperation(mms)

	testutils.AssertNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "operation")
}

func TestGetOperation_Sell(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_buy_test()
	amt := utils.DecimalFromString("11.34")
	price := utils.DecimalFromString("44881.330525")
	mms := get_mms("BTC", price)

	got, err := laccount.GetOperation(mms)
	testutils.AssertNil(t, err, "err")

	exp := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	exp.ExeId = got.ExeId
	exp.OpId = got.OpId
	exp.Timestamp = got.Timestamp
	exp.Status = model.PENDING
	exp.Results = model.OpResults{}
	exp.Type = model.AUTO

	testutils.AssertEq(t, exp, got, "operation")
}

func TestGetOperation_Sell_MinBaseQtyExceed(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_buy_test()
	btcSpotMarketLimits := laccount.GetSpotMarketLimits()["BTCUSDT"]
	btcSpotMarketLimits.MinBase = utils.DecimalFromString("12.00")
	btcSpotMarketLimits.MaxBase = utils.DecimalFromString("99999999")
	btcSpotMarketLimits.StepBase = utils.DecimalFromString("0.00000001")
	btcSpotMarketLimits.MinQuote = utils.DecimalFromString("0.1")
	laccount.SpotMarketLimits["BTCUSDT"] = btcSpotMarketLimits

	price := utils.DecimalFromString("44881.330525")
	mms := get_mms("BTC", price)

	got, err := laccount.GetOperation(mms)

	testutils.AssertNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "operation")
}

func TestGetOperation_StopLoss(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_buy_test()
	amt := utils.DecimalFromString("11.34")
	price := utils.DecimalFromString("31648.36")
	mms := get_mms("BTC", price)

	got, err := laccount.GetOperation(mms)
	testutils.AssertNil(t, err, "err")

	exp := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	exp.ExeId = got.ExeId
	exp.OpId = got.OpId
	exp.Timestamp = got.Timestamp
	exp.Status = model.PENDING
	exp.Results = model.OpResults{}
	exp.Type = model.AUTO

	testutils.AssertEq(t, exp, got, "operation")
}

func TestGetOperation_Buy(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_sell_test()
	amt := utils.DecimalFromString("999.99")
	price := utils.DecimalFromString("38.798975")
	mms := get_mms("DOT", price)

	got, err := laccount.GetOperation(mms)
	testutils.AssertNil(t, err, "err")

	exp := get_operation_test(amt, model.QUOTE_AMOUNT, "DOT", "USDT", model.BUY, price)
	exp.ExeId = got.ExeId
	exp.OpId = got.OpId
	exp.Timestamp = got.Timestamp
	exp.Status = model.PENDING
	exp.Results = model.OpResults{}
	exp.Type = model.AUTO

	testutils.AssertEq(t, exp, got, "operation")
}

func TestGetOperation_Buy_MinQuoteQtyExceeded(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_sell_test()
	dotSpotMarketLimits := laccount.GetSpotMarketLimits()["DOTUSDT"]
	dotSpotMarketLimits.MinBase = utils.DecimalFromString("0.00000001")
	dotSpotMarketLimits.MaxBase = utils.DecimalFromString("99999999")
	dotSpotMarketLimits.StepBase = utils.DecimalFromString("0.00000001")
	dotSpotMarketLimits.MinQuote = utils.DecimalFromString("1000")
	laccount.SpotMarketLimits["DOTUSDT"] = dotSpotMarketLimits

	price := utils.DecimalFromString("38.798975")
	mms := get_mms("DOT", price)

	got, err := laccount.GetOperation(mms)

	testutils.AssertNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "operation")
}

func TestGetOperation_MissProfit(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_sell_test()
	amt := utils.DecimalFromString("999.99")
	price := utils.DecimalFromString("59.34")
	mms := get_mms("DOT", price)

	got, err := laccount.GetOperation(mms)
	testutils.AssertNil(t, err, "err")

	exp := get_operation_test(amt, model.QUOTE_AMOUNT, "DOT", "USDT", model.BUY, price)
	exp.ExeId = got.ExeId
	exp.OpId = got.OpId
	exp.Timestamp = got.Timestamp
	exp.Status = model.PENDING
	exp.Results = model.OpResults{}
	exp.Type = model.AUTO

	testutils.AssertEq(t, exp, got, "operation")
}

func TestGetOperation_ZeroPrice(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_buy_test()
	price := decimal.Zero
	mms := get_mms("BTC", price)

	got, err := laccount.GetOperation(mms)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "operation")
}

/********************** Helpers *************************/

func mock_strategy_config(bt, st, slt, mpt string) func() config.StrategyConfig {
	old := config.GetStrategyConfig
	config.GetStrategyConfig = func() config.StrategyConfig {
		return config.StrategyConfig{
			Type: string(model.FIXED_THRESHOLD_STRATEGY),
			Config: struct {
				BuyThreshold        string
				SellThreshold       string
				StopLossThreshold   string
				MissProfitThreshold string
			}{
				BuyThreshold:        bt,
				SellThreshold:       st,
				StopLossThreshold:   slt,
				MissProfitThreshold: mpt}}
	}
	return old
}

func restore_strategy_config(old func() config.StrategyConfig) {
	config.GetStrategyConfig = old
}

func get_mms(asset string, lastPrice decimal.Decimal) model.MiniMarketStats {
	return model.MiniMarketStats{
		Event:       "event",
		Time:        time.Now().UnixMicro(),
		Asset:       asset,
		LastPrice:   lastPrice,
		OpenPrice:   utils.DecimalFromString("105.56"),
		HighPrice:   utils.DecimalFromString("197.45"),
		LowPrice:    utils.DecimalFromString("105.56"),
		QuoteVolume: utils.DecimalFromString("14455678.54"),
		BaseVolume:  utils.DecimalFromString("65395234.1665")}
}

func get_operation_test(amt decimal.Decimal, amtSide model.AmountSide, base, quote string,
	side model.OpSide, price decimal.Decimal) model.Operation {

	var baseAmt, quoteAmt decimal.Decimal
	if amtSide == model.BASE_AMOUNT {
		baseAmt = amt
		quoteAmt = amt.Mul(price).Round(8)
	} else {
		quoteAmt = amt
		baseAmt = quoteAmt.Div(price).Round(8)
	}

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
		Results: model.OpResults{
			ActualPrice: price,
			BaseDiff:    baseAmt,
			QuoteDiff:   quoteAmt,
			Spread:      decimal.Zero,
		},
		Status:    model.FILLED,
		Timestamp: time.Now().UnixMicro()}
}

func get_laccount_last_sell_test() LocalAccountFTS {
	return LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:        uuid.NewString(),
			ExeId:            uuid.NewString(),
			StrategyType:     model.FIXED_THRESHOLD_STRATEGY,
			SpotMarketLimits: get_spot_market_limits(),
			Timestamp:        time.Now().UnixMicro()},

		Ignored: map[string]decimal.Decimal{
			"USDT": utils.DecimalFromString("155.67"),
			"BUSD": utils.DecimalFromString("1232.45")},

		Assets: map[string]AssetStatusFTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             decimal.Zero,
				Usdt:               utils.DecimalFromString("24519.999"),
				LastOperationType:  OP_SELL_FTS,
				LastOperationPrice: utils.DecimalFromString("39560.45"),
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             decimal.Zero,
				Usdt:               utils.DecimalFromString("13443.12"),
				LastOperationType:  OP_SELL_FTS,
				LastOperationPrice: utils.DecimalFromString("4500.45")},
			"DOT": {
				Asset:              "DOT",
				Amount:             decimal.Zero,
				Usdt:               utils.DecimalFromString("999.99"),
				LastOperationType:  OP_SELL_FTS,
				LastOperationPrice: utils.DecimalFromString("49.45")}}}
}

func get_laccount_last_buy_test() LocalAccountFTS {
	return LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:        uuid.NewString(),
			ExeId:            uuid.NewString(),
			StrategyType:     model.FIXED_THRESHOLD_STRATEGY,
			SpotMarketLimits: get_spot_market_limits(),
			Timestamp:        time.Now().UnixMicro()},

		Ignored: map[string]decimal.Decimal{
			"USDT": utils.DecimalFromString("155.67"),
			"BUSD": utils.DecimalFromString("1232.45")},

		Assets: map[string]AssetStatusFTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             utils.DecimalFromString("11.34"),
				Usdt:               decimal.Zero,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: utils.DecimalFromString("39560.45"),
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             utils.DecimalFromString("29.12"),
				Usdt:               decimal.Zero,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: utils.DecimalFromString("4500.45")},
			"DOT": {
				Asset:              "DOT",
				Amount:             utils.DecimalFromString("13.67"),
				Usdt:               decimal.Zero,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: utils.DecimalFromString("49.45")}}}
}

func get_spot_market_limits() map[string]model.SpotMarketLimits {
	return map[string]model.SpotMarketLimits{
		"BTCUSDT": get_spot_market_limit(),
		"ETHUSDT": get_spot_market_limit(),
		"DOTUSDT": get_spot_market_limit()}
}

func get_spot_market_limit() model.SpotMarketLimits {
	return model.SpotMarketLimits{
		MinBase:  utils.DecimalFromString("0.00000001"),
		MaxBase:  utils.DecimalFromString("99999999"),
		StepBase: utils.DecimalFromString("0.00000001"),
		MinQuote: utils.DecimalFromString("0.1")}
}

func get_laccount_init_test() model.LocalAccountInit {
	return model.LocalAccountInit{
		ExeId: uuid.NewString(),
		RAccount: model.RemoteAccount{
			MakerCommission:  0,
			TakerCommission:  0,
			BuyerCommission:  0,
			SellerCommission: 0,
			Balances: []model.RemoteBalance{
				{Asset: "BTC", Amount: utils.DecimalFromString("11.34")},
				{Asset: "ETH", Amount: utils.DecimalFromString("29.12")},
				{Asset: "DOT", Amount: utils.DecimalFromString("13.67")},
				{Asset: "USDT", Amount: utils.DecimalFromString("155.67")},
				{Asset: "BUSD", Amount: utils.DecimalFromString("1232.45")}}},
		TradableAssetsPrice: map[string]model.AssetPrice{
			"BTC": {Asset: "BTC", Price: utils.DecimalFromString("39560.45")},
			"ETH": {Asset: "ETH", Price: utils.DecimalFromString("4500.45")},
			"DOT": {Asset: "DOT", Price: utils.DecimalFromString("49.45")}},
		StrategyType:     model.FIXED_THRESHOLD_STRATEGY,
		SpotMarketLimits: get_spot_market_limits()}
}
