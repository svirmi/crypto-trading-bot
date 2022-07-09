package pts

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

	// Testing that LocalAccountPTS is properly initialized and
	// zero balance cryptos are filtered out
	balances := laccountInit.RAccount.Balances
	balances = append(balances, model.RemoteBalance{Asset: "SHIBA", Amount: decimal.Zero})
	laccountInit.RAccount.Balances = balances
	got, err := LocalAccountPTS{}.Initialize(laccountInit)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	testutils.AssertEq(t, exp, got, "dts_laccount")

	// Testing that LocalAccountPTS is properly initialized and
	// the usdt field is set to zero if no USDT owned
	laccountInit = get_laccount_init_test()
	balances = laccountInit.RAccount.Balances
	balances = balances[:len(balances)-1]
	laccountInit.RAccount.Balances = balances
	got, err = LocalAccountPTS{}.Initialize(laccountInit)
	testutils.AssertNil(t, err, "err")

	exp = get_laccount_test()
	exp.Usdt = decimal.Zero
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	testutils.AssertEq(t, exp, got, "pts_laccount")
}

/********************** Testing RegisterTrading() *************************/

func TestRegisterTrading_BaseAmt_BuySide(t *testing.T) {
	logger.Initialize(false, true, true)
	laccount := get_laccount_test()

	amt := utils.DecimalFromString("0.1")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	got, err := laccount.RegisterTrading(op)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	assetStatus := exp.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("11.44")
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Usdt = utils.DecimalFromString("105.67")

	testutils.AssertEq(t, exp, got, "pts_laccount")
}

func TestRegisterTrading_BaseAmt_SellSide(t *testing.T) {
	logger.Initialize(false, true, true)
	laccount := get_laccount_test()

	amt := utils.DecimalFromString("11.34")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId

	got, err := laccount.RegisterTrading(op)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	assetStatus := exp.Assets["BTC"]
	assetStatus.Amount = decimal.Zero
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Usdt = utils.DecimalFromString("5825.67")

	testutils.AssertEq(t, exp, got, "pts_laccount")
}

func TestRegisterTrading_QuoteAmt_BuySide(t *testing.T) {
	logger.Initialize(false, true, true)
	laccount := get_laccount_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	got, err := laccount.RegisterTrading(op)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	assetStatus := exp.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("11.55134")
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Usdt = utils.DecimalFromString("50")

	testutils.AssertEq(t, exp, got, "pts_laccount")
}

func TestRegisterTrading_QuoteAmt_SellSide(t *testing.T) {
	logger.Initialize(false, true, true)
	laccount := get_laccount_test()

	amt := utils.DecimalFromString("10000")
	price := utils.DecimalFromString("43125.2")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId

	got, err := laccount.RegisterTrading(op)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_test()
	exp.ExeId = got.GetExeId()
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()
	assetStatus := exp.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("11.10811702")
	assetStatus.LastOperationPrice = price
	exp.Assets["BTC"] = assetStatus
	exp.Usdt = utils.DecimalFromString("10155.67")

	testutils.AssertEq(t, exp, got, "pts_laccount")
}

func TestRegisterTrading_WrongExeId(t *testing.T) {
	logger.Initialize(false, true, true)
	exp := get_laccount_test()
	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	testutils.AssertPanic(t, func() {
		exp.RegisterTrading(op)
	})
}

func TestRegisterTrading_OpFailed(t *testing.T) {
	logger.Initialize(false, true, true)
	exp := get_laccount_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = exp.ExeId
	op.Status = model.FAILED

	got, err := exp.RegisterTrading(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "pts_laccount")
}

func TestRegisterTrading_BadQuoteCurrency(t *testing.T) {
	logger.Initialize(false, true, true)
	exp := get_laccount_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "ETH", model.SELL, price)
	op.ExeId = exp.ExeId

	got, err := exp.RegisterTrading(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "pts_laccount")
}

func TestRegisterTrading_AssetNotFound(t *testing.T) {
	logger.Initialize(false, true, true)
	exp := get_laccount_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "CRO", "USDT", model.SELL, price)
	op.ExeId = exp.ExeId

	got, err := exp.RegisterTrading(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "pts_laccount")
}

func TestRegisterTrading_NegativeBalanceBase(t *testing.T) {
	logger.Initialize(false, true, true)
	exp := get_laccount_test()

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
	exp := get_laccount_test()

	amt := utils.DecimalFromString("1923789.12")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = exp.ExeId

	testutils.AssertPanic(t, func() {
		exp.RegisterTrading(op)
	})
}

/********************** Testing GetOperation() *************************/

/********************** GetAssetStatuses() *************************/

func TestGetAssetAmounts(t *testing.T) {
	laccount := get_laccount_test()

	exp := map[string]model.AssetAmount{
		"USDT": {"USDT", utils.DecimalFromString("155.67")},
		"BUSD": {"BUSD", utils.DecimalFromString("1232.45")},
		"BTC":  {"BTC", utils.DecimalFromString("11.34")},
		"ETH":  {"ETH", utils.DecimalFromString("29.12")},
		"DOT":  {"DOT", utils.DecimalFromString("13.67")}}

	got := laccount.GetAssetAmounts()
	testutils.AssertEq(t, exp, got, "asset_statuses")
}

/********************** Helpers *************************/

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

func get_laccount_test() LocalAccountPTS {
	return LocalAccountPTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.PTS_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

		Ignored: map[string]decimal.Decimal{
			"BUSD": utils.DecimalFromString("1232.45")},

		Assets: map[string]AssetStatusPTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             utils.DecimalFromString("11.34"),
				LastOperationPrice: utils.DecimalFromString("39560.45"),
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             utils.DecimalFromString("29.12"),
				LastOperationPrice: utils.DecimalFromString("4500.45")},
			"DOT": {
				Asset:              "DOT",
				Amount:             utils.DecimalFromString("13.67"),
				LastOperationPrice: utils.DecimalFromString("49.45")}},
		Usdt: utils.DecimalFromString("155.67")}
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
				{Asset: "BUSD", Amount: utils.DecimalFromString("1232.45")},
				{Asset: "USDT", Amount: utils.DecimalFromString("155.67")}}},
		TradableAssetsPrice: map[string]model.AssetPrice{
			"BTC": {Asset: "BTC", Price: utils.DecimalFromString("39560.45")},
			"ETH": {Asset: "ETH", Price: utils.DecimalFromString("4500.45")},
			"DOT": {Asset: "DOT", Price: utils.DecimalFromString("49.45")}},
		StrategyType: model.PTS_STRATEGY}
}
