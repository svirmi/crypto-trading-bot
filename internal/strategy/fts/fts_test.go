package fts

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
)

func TestInitialize(t *testing.T) {
	gotten, err := LocalAccountFTS{}.Initialize(get_laccount_init_test())
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_BaseAmt_BuySide(t *testing.T) {
	laccount := get_laccount_test()
	assetStatus := laccount.Assets["BTC"]
	assetStatus.Usdt = decimal.NewFromFloat32(50)
	laccount.Assets["BTC"] = assetStatus

	amt := decimal.NewFromFloat32(0.1)
	price := decimal.NewFromFloat32(500)
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()
	assetStatus = expected.Assets["BTC"]
	assetStatus.Amount = decimal.NewFromFloat32(11.44)
	assetStatus.Usdt = decimal.Zero
	assetStatus.LastOperationPrice = price
	expected.Assets["BTC"] = assetStatus

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_WrongExeId(t *testing.T) {
	laccount := get_laccount_test()
	assetStatus := laccount.Assets["BTC"]
	assetStatus.Usdt = decimal.NewFromFloat32(50)
	laccount.Assets["BTC"] = assetStatus

	amt := decimal.NewFromFloat32(0.1)
	price := decimal.NewFromFloat32(500)
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()
	assetStatus = expected.Assets["BTC"]
	assetStatus.Amount = decimal.NewFromFloat32(11.44)
	assetStatus.Usdt = decimal.Zero
	assetStatus.LastOperationPrice = price
	expected.Assets["BTC"] = assetStatus

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_BaseAmt_SellSide(t *testing.T) {
	laccount := get_laccount_test()

	amt := decimal.NewFromFloat32(11.34)
	price := decimal.NewFromFloat32(500)
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()
	assetStatus := expected.Assets["BTC"]
	assetStatus.Amount = decimal.Zero
	assetStatus.Usdt = decimal.NewFromFloat32(5670)
	assetStatus.LastOperationType = OP_SELL_FTS
	assetStatus.LastOperationPrice = price
	expected.Assets["BTC"] = assetStatus
	expected.Ignored["USDT"] = decimal.NewFromFloat32(155.67)

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_QuoteAmt_BuySide(t *testing.T) {
	laccount := get_laccount_test()
	assetStatus := laccount.Assets["BTC"]
	assetStatus.Usdt = decimal.NewFromFloat32(110.18)
	laccount.Assets["BTC"] = assetStatus

	amt := decimal.NewFromFloat32(105.67)
	price := decimal.NewFromFloat32(500)
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err != nil {
		t.Fatalf("err: expected = nil, gotten = %v", err)
	}

	expected := get_laccount_test()
	expected.ExeId = gotten.GetExeId()
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()
	assetStatus = expected.Assets["BTC"]
	assetStatus.Amount = decimal.NewFromFloat32(11.55134)
	assetStatus.Usdt = decimal.NewFromFloat32(4.51)
	assetStatus.LastOperationType = OP_BUY_FTS
	assetStatus.LastOperationPrice = price
	expected.Assets["BTC"] = assetStatus
	expected.Ignored["USDT"] = decimal.NewFromFloat32(155.67)

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_QuoteAmt_SellSide(t *testing.T) {
	laccount := get_laccount_test()

	amt := decimal.NewFromFloat32(105.67)
	price := decimal.NewFromFloat32(500)
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten = nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func TestRegisterTrading_OpFailed(t *testing.T) {
	laccount := get_laccount_test()

	amt := decimal.NewFromFloat32(105.67)
	price := decimal.NewFromFloat32(500)
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
	laccount := get_laccount_test()

	amt := decimal.NewFromFloat32(105.67)
	price := decimal.NewFromFloat32(500)
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
	laccount := get_laccount_test()

	amt := decimal.NewFromFloat32(105.67)
	price := decimal.NewFromFloat32(500)
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
	laccount := get_laccount_test()

	amt := decimal.NewFromFloat32(1923789.12)
	price := decimal.NewFromFloat32(500)
	op := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten == nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func TestRegisterTrading_NegativeBalanceQuote(t *testing.T) {
	laccount := get_laccount_test()

	amt := decimal.NewFromFloat32(1923789.12)
	price := decimal.NewFromFloat32(500)
	op := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten == nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func get_operation_test(amt decimal.Decimal, amtSide model.AmountSide, base, quote string,
	side model.OpSide, price decimal.Decimal) model.Operation {

	var baseAmt, quoteAmt decimal.Decimal
	if amtSide == model.BASE_AMOUNT {
		baseAmt = amt
		quoteAmt = amt.Mul(price)
	} else {
		quoteAmt = amt
		baseAmt = quoteAmt.Div(price)
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

func get_laccount_test() LocalAccountFTS {
	return LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        uuid.NewString(),
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMicro()},

		Ignored: map[string]decimal.Decimal{
			"USDT": decimal.NewFromFloat32(155.67),
			"BUSD": decimal.NewFromFloat32(1232.45)},

		Assets: map[string]AssetStatusFTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             decimal.NewFromFloat32(11.34),
				Usdt:               decimal.Zero,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: decimal.NewFromFloat32(39560.45),
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             decimal.NewFromFloat32(29.12),
				Usdt:               decimal.Zero,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: decimal.NewFromFloat32(4500.45)},
			"DOT": {
				Asset:              "DOT",
				Amount:             decimal.NewFromFloat32(13.67),
				Usdt:               decimal.Zero,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: decimal.NewFromFloat32(49.45)}}}
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
				{Asset: "BTC", Amount: decimal.NewFromFloat32(11.34)},
				{Asset: "ETH", Amount: decimal.NewFromFloat32(29.12)},
				{Asset: "DOT", Amount: decimal.NewFromFloat32(13.67)},
				{Asset: "USDT", Amount: decimal.NewFromFloat32(155.67)},
				{Asset: "BUSD", Amount: decimal.NewFromFloat32(1232.45)}}},
		TradableAssetsPrice: map[string]model.AssetPrice{
			"BTC": {Asset: "BTC", Price: decimal.NewFromFloat32(39560.45)},
			"ETH": {Asset: "ETH", Price: decimal.NewFromFloat32(4500.45)},
			"DOT": {Asset: "DOT", Price: decimal.NewFromFloat32(49.45)}},
		StrategyType: model.FIXED_THRESHOLD_STRATEGY}
}
