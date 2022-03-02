package fts

import (
	"testing"
	"time"

	"github.com/google/uuid"
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
	assetStatus.Usdt = 50
	laccount.Assets["BTC"] = assetStatus

	op := get_operation_test(0.1, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, 500)
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
	// TODO: substitute float32 with shopspring/decimal
	assetStatus.Amount = 11.440001
	assetStatus.Usdt = 0
	assetStatus.LastOperationPrice = 500
	expected.Assets["BTC"] = assetStatus

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_WrongExeId(t *testing.T) {
	laccount := get_laccount_test()
	assetStatus := laccount.Assets["BTC"]
	assetStatus.Usdt = 50
	laccount.Assets["BTC"] = assetStatus

	op := get_operation_test(0.1, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, 500)
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
	// TODO: substitute float32 with shopspring/decimal
	assetStatus.Amount = 11.440001
	assetStatus.Usdt = 0
	assetStatus.LastOperationPrice = 500
	expected.Assets["BTC"] = assetStatus

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_BaseAmt_SellSide(t *testing.T) {
	laccount := get_laccount_test()

	op := get_operation_test(11.34, model.BASE_AMOUNT, "BTC", "USDT", model.SELL, 500)
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
	assetStatus.Amount = 0
	assetStatus.Usdt = 5670
	assetStatus.LastOperationType = OP_SELL_FTS
	assetStatus.LastOperationPrice = 500
	expected.Assets["BTC"] = assetStatus
	expected.Ignored["USDT"] = 155.67

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_QuoteAmt_BuySide(t *testing.T) {
	laccount := get_laccount_test()
	assetStatus := laccount.Assets["BTC"]
	assetStatus.Usdt = 110.18
	laccount.Assets["BTC"] = assetStatus

	op := get_operation_test(105.67, model.QUOTE_AMOUNT, "BTC", "USDT", model.BUY, 500)
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
	assetStatus.Amount = 11.55134
	assetStatus.Usdt = 4.510002
	assetStatus.LastOperationType = OP_BUY_FTS
	assetStatus.LastOperationPrice = 500
	expected.Assets["BTC"] = assetStatus
	expected.Ignored["USDT"] = 155.67

	testutils.AssertStructEq(t, expected, gotten)
}

func TestRegisterTrading_QuoteAmt_SellSide(t *testing.T) {
	laccount := get_laccount_test()

	op := get_operation_test(105.67, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, 500)

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten = nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func TestRegisterTrading_OpFailed(t *testing.T) {
	laccount := get_laccount_test()

	op := get_operation_test(105.67, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, 500)
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

	op := get_operation_test(105.67, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, 500)
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

	op := get_operation_test(105.67, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, 500)
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

	op := get_operation_test(1923789.12, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, 500)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten == nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func TestRegisterTrading_NegativeBalanceQuote(t *testing.T) {
	laccount := get_laccount_test()

	op := get_operation_test(1923789.12, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, 500)
	op.ExeId = laccount.ExeId

	gotten, err := laccount.RegisterTrading(op)
	if err == nil {
		t.Fatalf("err: expected != nil, gotten == nil")
	}
	testutils.AssertStructEq(t, laccount, gotten)
}

func get_operation_test(amt float32, amtSide model.AmountSide, base, quote string,
	side model.OpSide, price float32) model.Operation {

	var baseAmt, quoteAmt float32
	if amtSide == model.BASE_AMOUNT {
		baseAmt = amt
		quoteAmt = amt * price
	} else {
		quoteAmt = amt
		baseAmt = quoteAmt / price
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
			Spread:      0,
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

		Ignored: map[string]float32{
			"USDT": 155.67,
			"BUSD": 1232.45},

		Assets: map[string]AssetStatusFTS{
			"BTC": {
				Asset:              "BTC",
				Amount:             11.34,
				Usdt:               0,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: 39560.45,
			},
			"ETH": {
				Asset:              "ETH",
				Amount:             29.12,
				Usdt:               0,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: 4500.45},
			"DOT": {
				Asset:              "DOT",
				Amount:             13.67,
				Usdt:               0,
				LastOperationType:  OP_BUY_FTS,
				LastOperationPrice: 49.45}}}
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
				{Asset: "BTC", Amount: 11.34},
				{Asset: "ETH", Amount: 29.12},
				{Asset: "DOT", Amount: 13.67},
				{Asset: "USDT", Amount: 155.67},
				{Asset: "BUSD", Amount: 1232.45}}},
		TradableAssetsPrice: map[string]model.AssetPrice{
			"BTC": {Asset: "BTC", Price: 39560.45},
			"ETH": {Asset: "ETH", Price: 4500.45},
			"DOT": {Asset: "DOT", Price: 49.45}},
		StrategyType: model.FIXED_THRESHOLD_STRATEGY}
}
