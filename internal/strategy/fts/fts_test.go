package fts

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

/********************** Testing Initialize() *************************/

func TestInitialize(t *testing.T) {
	gotten, err := LocalAccountFTS{}.Initialize(get_laccount_init_test())
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_last_buy_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()

	testutils.AssertStructEq(t, expected, gotten)
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

	gotten, err := laccount.RegisterTrading(op)
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_last_buy_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()
	assetStatus = expected.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("11.44")
	assetStatus.Usdt = decimal.Zero
	assetStatus.LastOperationPrice = price
	expected.Assets["BTC"] = assetStatus

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_WrongExeId(t *testing.T) {
	laccount := get_laccount_last_buy_test()
	assetStatus := laccount.Assets["BTC"]
	assetStatus.Usdt = utils.DecimalFromString("50")
	laccount.Assets["BTC"] = assetStatus

	amt := utils.DecimalFromString("0.1")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_last_buy_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()
	assetStatus = expected.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("11.44")
	assetStatus.Usdt = decimal.Zero
	assetStatus.LastOperationPrice = price
	expected.Assets["BTC"] = assetStatus

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_BaseAmt_SellSide(t *testing.T) {
	laccount := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("11.34")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_last_buy_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()
	assetStatus := expected.Assets["BTC"]
	assetStatus.Amount = decimal.Zero
	assetStatus.Usdt = utils.DecimalFromString("5670")
	assetStatus.LastOperationType = OP_SELL_FTS
	assetStatus.LastOperationPrice = price
	expected.Assets["BTC"] = assetStatus
	expected.Ignored["USDT"] = utils.DecimalFromString("155.67")

	testutils.AssertStructEq(t, expected, gotten)
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

	gotten, err := laccount.RegisterTrading(op)
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_last_buy_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()
	assetStatus = expected.Assets["BTC"]
	assetStatus.Amount = utils.DecimalFromString("11.55134")
	assetStatus.Usdt = utils.DecimalFromString("4.51")
	assetStatus.LastOperationType = OP_BUY_FTS
	assetStatus.LastOperationPrice = price
	expected.Assets["BTC"] = assetStatus
	expected.Ignored["USDT"] = utils.DecimalFromString("155.67")

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_QuoteAmt_SellSide(t *testing.T) {
	laccount := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten = nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func TestRegisterTrading_OpFailed(t *testing.T) {
	laccount := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId
	op.Status = model.FAILED

	gotten, err := laccount.RegisterTrading(op)
	if err != nil {
		t.Fatalf("err: expected == nil, gotten = %s", err)
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func TestRegisterTrading_BadQuoteCurrency(t *testing.T) {
	laccount := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId
	op.Quote = "ETH"

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten == nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func TestRegisterTrading_AssetNotFound(t *testing.T) {
	laccount := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("105.67")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId
	op.Base = "CRO"

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten == nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func TestRegisterTrading_NegativeBalanceBase(t *testing.T) {
	laccount := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("1923789.12")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten == nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func TestRegisterTrading_NegativeBalanceQuote(t *testing.T) {
	laccount := get_laccount_last_buy_test()

	amt := utils.DecimalFromString("1923789.12")
	price := utils.DecimalFromString("500")
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten == nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

/********************** Testing GetOperation() *************************/

func TestGetOperation_AssetNotFound(t *testing.T) {
	laccount := get_laccount_last_buy_test()
	mms := get_mms("CRO", utils.DecimalFromString("0.55"))

	op, err := laccount.GetOperation(mms)

	if !op.IsEmpty() {
		t.Errorf("op: expected empty, gotten %v", op)
	}
	if err == nil {
		t.Errorf("err: expected != nil, gotten nil")
	}
}

func TestGetOperation_Noop(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_buy_test()
	mms := get_mms("BTC", utils.DecimalFromString("39560.1"))

	gotten, err := laccount.GetOperation(mms)

	if err != nil {
		t.Errorf("err: expected == nil, gotten = %v", err)
	}
	if !gotten.IsEmpty() {
		t.Errorf("op: expected empty, gotten %v", gotten)
	}
}

func TestGetOperation_Sell(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_buy_test()
	amt := utils.DecimalFromString("11.34")
	price := utils.DecimalFromString("44881.330525")
	mms := get_mms("BTC", price)

	gotten, err := laccount.GetOperation(mms)
	if err != nil {
		t.Errorf("err: expected == nil, gotten = %v", err)
	}

	expected := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	expected.ExeId = gotten.ExeId
	expected.OpId = gotten.OpId
	expected.Timestamp = gotten.Timestamp
	expected.Status = model.PENDING
	expected.Results = model.OpResults{}
	expected.Type = model.AUTO

	testutils.AssertStructEq(t, expected, gotten)
}

func TestGetOperation_StopLoss(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_buy_test()
	amt := utils.DecimalFromString("11.34")
	price := utils.DecimalFromString("31648.36")
	mms := get_mms("BTC", price)

	gotten, err := laccount.GetOperation(mms)
	if err != nil {
		t.Errorf("err: expected == nil, gotten = %v", err)
	}

	expected := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	expected.ExeId = gotten.ExeId
	expected.OpId = gotten.OpId
	expected.Timestamp = gotten.Timestamp
	expected.Status = model.PENDING
	expected.Results = model.OpResults{}
	expected.Type = model.AUTO

	testutils.AssertStructEq(t, expected, gotten)
}

func TestGetOperation_Buy(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_sell_test()
	amt := utils.DecimalFromString("999.99")
	price := utils.DecimalFromString("38.798975")
	mms := get_mms("DOT", price)

	gotten, err := laccount.GetOperation(mms)
	if err != nil {
		t.Errorf("err: expected == nil, gotten = %v", err)
	}

	expected := get_operation_test(amt, model.QUOTE_AMOUNT, "DOT", "USDT", model.BUY, price)
	expected.ExeId = gotten.ExeId
	expected.OpId = gotten.OpId
	expected.Timestamp = gotten.Timestamp
	expected.Status = model.PENDING
	expected.Results = model.OpResults{}
	expected.Type = model.AUTO

	testutils.AssertStructEq(t, expected, gotten)
}

func TestGetOperation_MissProfit(t *testing.T) {
	old := mock_strategy_config("13.45", "13.45", "20", "20")
	defer restore_strategy_config(old)

	laccount := get_laccount_last_sell_test()
	amt := utils.DecimalFromString("999.99")
	price := utils.DecimalFromString("59.34")
	mms := get_mms("DOT", price)

	gotten, err := laccount.GetOperation(mms)
	if err != nil {
		t.Errorf("err: expected == nil, gotten = %v", err)
	}

	expected := get_operation_test(amt, model.QUOTE_AMOUNT, "DOT", "USDT", model.BUY, price)
	expected.ExeId = gotten.ExeId
	expected.OpId = gotten.OpId
	expected.Timestamp = gotten.Timestamp
	expected.Status = model.PENDING
	expected.Results = model.OpResults{}
	expected.Type = model.AUTO

	testutils.AssertStructEq(t, expected, gotten)
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
			BaseAmount:  baseAmt,
			QuoteAmount: quoteAmt,
			Spread:      decimal.Zero,
		},
		Status:    model.FILLED,
		Timestamp: time.Now().UnixMicro()}
}

func get_laccount_last_sell_test() LocalAccountFTS {
	return LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

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
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

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
		StrategyType: model.FIXED_THRESHOLD_STRATEGY}
}
