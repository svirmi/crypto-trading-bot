package dts

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

/********************** Testing Initialize() *************************/

func TestInitialize(t *testing.T) {
	logger.Initialize(false, true, true)
	laccountInit := get_laccount_init_test()

	// Addind a zero balance crypto that should be filtered out
	balances := laccountInit.RAccount.Balances
	balances = append(balances, model.RemoteBalance{Asset: "SHIBA", Amount: decimal.Zero})
	laccountInit.RAccount.Balances = balances
	got, err := LocalAccountDTS{}.Initialize(laccountInit)
	testutils.AssertNil(t, err, "err")

	// Testing the initialize function
	exp := get_laccount_last_buy_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()

	testutils.AssertEq(t, exp, got, "dts_laccount")
}

/********************** Testing RegisterTrading() *************************/

func TestRegisterTrading_BaseAmt_BuySide(t *testing.T) {
	logger.Initialize(false, true, true)
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

	testutils.AssertEq(t, exp, got, "dts_laccount")
}

func TestRegisterTrading_BaseAmt_SellSide(t *testing.T) {
	logger.Initialize(false, true, true)
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
	assetStatus.LastOperationType = OP_SELL_DTS
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Ignored["USDT"] = utils.DecimalFromString("155.67")

	testutils.AssertEq(t, exp, got, "dts_laccount")
}

func TestRegisterTrading_QuoteAmt_BuySide(t *testing.T) {
	logger.Initialize(false, true, true)
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
	assetStatus.LastOperationType = OP_BUY_DTS
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Ignored["USDT"] = utils.DecimalFromString("155.67")

	testutils.AssertEq(t, exp, got, "dts_laccount")
}

func TestRegisterTrading_QuoteAmt_SellSide(t *testing.T) {
	logger.Initialize(false, true, true)
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
	assetStatus.LastOperationType = OP_SELL_DTS
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Ignored["USDT"] = utils.DecimalFromString("155.67")

	testutils.AssertEq(t, exp, got, "dts_laccount")
}

func TestRegisterTrading_WrongExeId(t *testing.T) {
	logger.Initialize(false, true, true)
	exp := get_laccount_last_buy_test()
	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	testutils.AssertPanic(t, func() {
		exp.RegisterTrading(op)
	})
}

func TestRegisterTrading_OpFailed(t *testing.T) {
	logger.Initialize(false, true, true)
	exp := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = exp.ExeId
	op.Status = model.FAILED

	got, err := exp.RegisterTrading(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "dts_laccount")
}

func TestRegisterTrading_BadQuoteCurrency(t *testing.T) {
	logger.Initialize(false, true, true)
	exp := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "ETH", model.SELL, price)
	op.ExeId = exp.ExeId

	got, err := exp.RegisterTrading(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "dts_laccount")
}

func TestRegisterTrading_AssetNotFound(t *testing.T) {
	logger.Initialize(false, true, true)
	exp := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "CRO", "USDT", model.SELL, price)
	op.ExeId = exp.ExeId

	got, err := exp.RegisterTrading(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "dts_laccount")
}

func TestRegisterTrading_NegativeBalanceBase(t *testing.T) {
	logger.Initialize(false, true, true)
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
	logger.Initialize(false, true, true)
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
	logger.Initialize(false, true, true)
	laccount := get_laccount_last_buy_test()
	mms := get_mms("CRO", utils.DecimalFromString("0.55"))

	props := get_props("13.45", "13.45", "20", "20")
	op, err := laccount.GetOperation(props, mms, get_spot_market_limit())

	testutils.AssertTrue(t, op.IsEmpty(), "operation")
	testutils.AssertNotNil(t, err, "err")
}

func TestGetOperation_Noop(t *testing.T) {
	logger.Initialize(false, true, true)

	laccount := get_laccount_last_buy_test()
	mms := get_mms("BTC", utils.DecimalFromString("39560.1"))

	props := get_props("13.45", "13.45", "20", "20")
	got, err := laccount.GetOperation(props, mms, get_spot_market_limit())

	testutils.AssertNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "operation")
}

func TestGetOperation_Sell(t *testing.T) {
	logger.Initialize(false, true, true)

	laccount := get_laccount_last_buy_test()
	amt := utils.DecimalFromString("11.34")
	price := utils.DecimalFromString("44881.330525")
	mms := get_mms("BTC", price)

	props := get_props("13.45", "13.45", "20", "20")
	got, err := laccount.GetOperation(props, mms, get_spot_market_limit())
	testutils.AssertNil(t, err, "err")

	exp := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	exp.ExeId = got.ExeId
	exp.OpId = got.OpId
	exp.Timestamp = got.Timestamp
	exp.Status = model.PENDING
	exp.Results = model.OpResults{}
	exp.Type = model.AUTO
	exp.Cause = _DTS_SELL_DESC

	testutils.AssertEq(t, exp, got, "operation")
}

func TestGetOperation_Sell_MinBaseQtyExceed(t *testing.T) {
	logger.Initialize(false, true, true)

	laccount := get_laccount_last_buy_test()
	btcSpotMarketLimits := get_spot_market_limit()
	btcSpotMarketLimits.MinBase = utils.DecimalFromString("12.00")
	btcSpotMarketLimits.MaxBase = utils.DecimalFromString("99999999")
	btcSpotMarketLimits.StepBase = utils.DecimalFromString("0.00000001")
	btcSpotMarketLimits.MinQuote = utils.DecimalFromString("0.1")

	price := utils.DecimalFromString("44881.330525")
	mms := get_mms("BTC", price)

	props := get_props("13.45", "13.45", "20", "20")
	got, err := laccount.GetOperation(props, mms, btcSpotMarketLimits)
	testutils.AssertNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "operation")
}

func TestGetOperation_StopLoss(t *testing.T) {
	logger.Initialize(false, true, true)

	laccount := get_laccount_last_buy_test()
	amt := utils.DecimalFromString("11.34")
	price := utils.DecimalFromString("31648.36")
	mms := get_mms("BTC", price)

	props := get_props("13.45", "13.45", "20", "20")
	got, err := laccount.GetOperation(props, mms, get_spot_market_limit())
	testutils.AssertNil(t, err, "err")

	exp := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	exp.ExeId = got.ExeId
	exp.OpId = got.OpId
	exp.Timestamp = got.Timestamp
	exp.Status = model.PENDING
	exp.Results = model.OpResults{}
	exp.Type = model.AUTO
	exp.Cause = _DTS_STOP_LOSS_DESC

	testutils.AssertEq(t, exp, got, "operation")
}

func TestGetOperation_Buy(t *testing.T) {
	logger.Initialize(false, true, true)

	laccount := get_laccount_last_sell_test()
	amt := utils.DecimalFromString("999.99")
	price := utils.DecimalFromString("38.798975")
	mms := get_mms("DOT", price)

	props := get_props("13.45", "13.45", "20", "20")
	got, err := laccount.GetOperation(props, mms, get_spot_market_limit())
	testutils.AssertNil(t, err, "err")

	exp := get_operation_test(amt, model.QUOTE_AMOUNT, "DOT", "USDT", model.BUY, price)
	exp.ExeId = got.ExeId
	exp.OpId = got.OpId
	exp.Timestamp = got.Timestamp
	exp.Status = model.PENDING
	exp.Results = model.OpResults{}
	exp.Type = model.AUTO
	exp.Cause = _DTS_BUY_DESC

	testutils.AssertEq(t, exp, got, "operation")
}

func TestGetOperation_Buy_MinQuoteQtyExceeded(t *testing.T) {
	logger.Initialize(false, true, true)

	laccount := get_laccount_last_sell_test()
	dotSpotMarketLimits := get_spot_market_limit()
	dotSpotMarketLimits.MinBase = utils.DecimalFromString("0.00000001")
	dotSpotMarketLimits.MaxBase = utils.DecimalFromString("99999999")
	dotSpotMarketLimits.StepBase = utils.DecimalFromString("0.00000001")
	dotSpotMarketLimits.MinQuote = utils.DecimalFromString("1000")

	price := utils.DecimalFromString("38.798975")
	mms := get_mms("DOT", price)

	props := get_props("13.45", "13.45", "20", "20")
	got, err := laccount.GetOperation(props, mms, dotSpotMarketLimits)

	testutils.AssertNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "operation")
}

func TestGetOperation_MissProfit(t *testing.T) {
	logger.Initialize(false, true, true)

	laccount := get_laccount_last_sell_test()
	amt := utils.DecimalFromString("999.99")
	price := utils.DecimalFromString("59.34")
	mms := get_mms("DOT", price)

	props := get_props("13.45", "13.45", "20", "20")
	got, err := laccount.GetOperation(props, mms, get_spot_market_limit())
	testutils.AssertNil(t, err, "err")

	exp := get_operation_test(amt, model.QUOTE_AMOUNT, "DOT", "USDT", model.BUY, price)
	exp.ExeId = got.ExeId
	exp.OpId = got.OpId
	exp.Timestamp = got.Timestamp
	exp.Status = model.PENDING
	exp.Results = model.OpResults{}
	exp.Type = model.AUTO
	exp.Cause = _DTS_MISS_PROFIT_DESC

	testutils.AssertEq(t, exp, got, "operation")
}

func TestGetOperation_ZeroPrice(t *testing.T) {
	logger.Initialize(false, true, true)

	laccount := get_laccount_last_buy_test()
	price := decimal.Zero
	mms := get_mms("BTC", price)

	props := get_props("13.45", "13.45", "20", "20")
	got, err := laccount.GetOperation(props, mms, get_spot_market_limit())

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "operation")
}

/********************** GetAssetStatuses() *************************/

func TestGetAssetAmounts(t *testing.T) {
	laccount := get_laccount_last_buy_test()
	btc := laccount.Assets["BTC"]
	eth := laccount.Assets["ETH"]
	btc.Usdt = utils.DecimalFromString("0.0035")
	eth.Usdt = utils.DecimalFromString("0.5")
	laccount.Assets["BTC"] = btc
	laccount.Assets["ETH"] = eth

	exp := map[string]model.AssetAmount{
		"USDT": {"USDT", utils.DecimalFromString("156.1735")},
		"BUSD": {"BUSD", utils.DecimalFromString("1232.45")},
		"BTC":  {"BTC", utils.DecimalFromString("11.34")},
		"ETH":  {"ETH", utils.DecimalFromString("29.12")},
		"DOT":  {"DOT", utils.DecimalFromString("13.67")}}

	got := laccount.GetAssetAmounts()
	testutils.AssertEq(t, exp, got, "asset_statuses")
}

/********************** Helpers *************************/

func get_props(bt, st, slt, mpt string) map[string]string {
	return map[string]string{
		_BUY_THRESHOLD:         bt,
		_SELL_THRESHOLD:        st,
		_MISS_PROFIT_THRESHOLD: mpt,
		_STOP_LOSS_THRESHOLD:   slt}
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

func get_laccount_last_sell_test() LocalAccountDTS {
	return LocalAccountDTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.DTS_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

		Ignored: map[string]decimal.Decimal{
			"USDT": utils.DecimalFromString("155.67"),
			"BUSD": utils.DecimalFromString("1232.45")},

		Assets: map[string]AssetStatusDTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             decimal.Zero,
				Usdt:               utils.DecimalFromString("24519.999"),
				LastOperationType:  OP_SELL_DTS,
				LastOperationPrice: utils.DecimalFromString("39560.45"),
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             decimal.Zero,
				Usdt:               utils.DecimalFromString("13443.12"),
				LastOperationType:  OP_SELL_DTS,
				LastOperationPrice: utils.DecimalFromString("4500.45")},
			"DOT": {
				Asset:              "DOT",
				Amount:             decimal.Zero,
				Usdt:               utils.DecimalFromString("999.99"),
				LastOperationType:  OP_SELL_DTS,
				LastOperationPrice: utils.DecimalFromString("49.45")}}}
}

func get_laccount_last_buy_test() LocalAccountDTS {
	return LocalAccountDTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.DTS_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

		Ignored: map[string]decimal.Decimal{
			"USDT": utils.DecimalFromString("155.67"),
			"BUSD": utils.DecimalFromString("1232.45")},

		Assets: map[string]AssetStatusDTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             utils.DecimalFromString("11.34"),
				Usdt:               decimal.Zero,
				LastOperationType:  OP_BUY_DTS,
				LastOperationPrice: utils.DecimalFromString("39560.45"),
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             utils.DecimalFromString("29.12"),
				Usdt:               decimal.Zero,
				LastOperationType:  OP_BUY_DTS,
				LastOperationPrice: utils.DecimalFromString("4500.45")},
			"DOT": {
				Asset:              "DOT",
				Amount:             utils.DecimalFromString("13.67"),
				Usdt:               decimal.Zero,
				LastOperationType:  OP_BUY_DTS,
				LastOperationPrice: utils.DecimalFromString("49.45")}}}
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
		StrategyType: model.DTS_STRATEGY}
}
