package binance

import (
	"testing"
	"time"

	binanceapi "github.com/adshao/go-binance/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func TestFilterTradableAsset(t *testing.T) {
	old := symbols
	defer func() {
		symbols = old
	}()

	symbols = get_symbols()
	assets := []string{"BTC", "ETH", "TRX", "BNB"}

	gotten := FilterTradableAssets(assets)
	expected := []string{"BTC", "ETH"}

	testutils.AssertStructEq(t, expected, gotten)
}

func TestGetAssetsValue(t *testing.T) {
	old_get_price := binance_get_price
	old_symbols := symbols
	defer func() {
		binance_get_price = old_get_price
		symbols = old_symbols
	}()

	index := 0
	binance_get_price = func(*binanceapi.ListPricesService) ([]*binanceapi.SymbolPrice, error) {
		rprices := []*binanceapi.SymbolPrice{
			{Symbol: "BTCUSDT", Price: "35998.34"},
			{Symbol: "ETHUSDT", Price: "44978.12"},
			{Symbol: "DOTUSDT", Price: "98.12"}}
		results := []*binanceapi.SymbolPrice{rprices[index]}
		index++
		return results, nil
	}
	symbols = get_symbols()

	gotten, err := GetAssetsValue([]string{"BTC", "ETH", "DOT", "BNB", "TRX"})
	if err != nil {
		t.Errorf("err: expected nil, gotten = %v", err)
	}

	expected := make(map[string]model.AssetPrice)
	expected["BTC"] = model.AssetPrice{Asset: "BTC", Price: utils.DecimalFromString("35998.34")}
	expected["ETH"] = model.AssetPrice{Asset: "ETH", Price: utils.DecimalFromString("44978.12")}
	expected["DOT"] = model.AssetPrice{Asset: "DOT", Price: utils.DecimalFromString("98.12")}

	testutils.AssertStructEq(t, expected, gotten)
}

func TestGetAccount(t *testing.T) {
	old := binance_get_account
	defer func() {
		binance_get_account = old
	}()

	binance_get_account = func(*binanceapi.GetAccountService) (*binanceapi.Account, error) {
		return get_remote_binance_account(), nil
	}

	gotten, err := GetAccout()
	if err != nil {
		t.Errorf("err: expected = nil, gotten = %s", err)
	}
	expected := get_remote_account()

	testutils.AssertStructEq(t, expected, gotten)
}

func TestSendMarketOrder_NoSuchSymbol(t *testing.T) {
	old := symbols
	defer func() {
		symbols = old
	}()

	symbols = get_symbols()
	amt := utils.DecimalFromString("12.34")
	price := utils.DecimalFromString("56.12")
	op := get_operation_test(amt, model.BASE_AMOUNT, "TRX", "USDT", model.BUY, price)

	gotten, err := SendSpotMarketOrder(op)
	if !gotten.IsEmpty() {
		t.Errorf("op: expected empty, gotten = %v", op)
	}
	if err == nil {
		t.Errorf("err: expected != nil, gotten = nil")
	}
}

func TestSendMarketOrder_Direct_Buy_BaseAmt(t *testing.T) {
	old_symbols := symbols
	old_create_order := binance_create_order
	old_get_spot_limits := GetSpotMarketLimits
	defer func() {
		symbols = old_symbols
		binance_create_order = old_create_order
		GetSpotMarketLimits = old_get_spot_limits
	}()

	symbols = get_symbols()
	binance_create_order = func(*binanceapi.CreateOrderService) (*binanceapi.CreateOrderResponse, error) {
		return &binanceapi.CreateOrderResponse{
			Symbol: "BTCUSDT",
			Side:   binanceapi.SideTypeBuy,
			Status: binanceapi.OrderStatusTypeFilled}, nil
	}
	GetSpotMarketLimits = func(symbol string) (model.SpotMarketLimits, error) {
		return model.SpotMarketLimits{
			MinBase:  utils.DecimalFromString("0.00000001"),
			MaxBase:  utils.DecimalFromString("99999999"),
			StepBase: utils.DecimalFromString("0.00000001"),
			MinQuote: utils.DecimalFromString("0.00000001")}, nil
	}

	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("3746.34")
	expected := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	gotten, err := SendSpotMarketOrder(expected)
	expected.Timestamp = gotten.Timestamp

	if err != nil {
		t.Errorf("err: expected = nil, gotten = %v", err)
	}
	testutils.AssertStructEq(t, expected, gotten)
}

func TestSendMarketOrder_Direct_Buy_BaseAmt_BelowMinBase(t *testing.T) {
	old_symbols := symbols
	old_get_spot_limits := GetSpotMarketLimits
	defer func() {
		symbols = old_symbols
		GetSpotMarketLimits = old_get_spot_limits
	}()

	symbols = get_symbols()
	GetSpotMarketLimits = func(symbol string) (model.SpotMarketLimits, error) {
		return model.SpotMarketLimits{
			MinBase:  utils.DecimalFromString("100.1"),
			MaxBase:  utils.DecimalFromString("99999999"),
			StepBase: utils.DecimalFromString("0.00000001"),
			MinQuote: utils.DecimalFromString("0.00000001")}, nil
	}

	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("3746.34")
	expected := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	_, err := SendSpotMarketOrder(expected)

	if err == nil {
		t.Errorf("err: expected != nil, gotten = nil")
	}
}

func TestSendMarketOrder_Direct_Buy_BaseAmt_SpotDisabled(t *testing.T) {
	old_symbols := symbols
	old_get_spot_limits := GetSpotMarketLimits
	defer func() {
		symbols = old_symbols
		GetSpotMarketLimits = old_get_spot_limits
	}()

	symbols = get_symbols()
	GetSpotMarketLimits = func(symbol string) (model.SpotMarketLimits, error) {
		return model.SpotMarketLimits{
			MinBase:  utils.DecimalFromString("0.00000001"),
			MaxBase:  utils.DecimalFromString("99999999"),
			StepBase: utils.DecimalFromString("0.00000001"),
			MinQuote: utils.DecimalFromString("0.00000001")}, nil
	}

	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("3746.34")
	expected := get_operation_test(amt, model.BASE_AMOUNT, "SHIBA", "USDT", model.BUY, price)
	_, err := SendSpotMarketOrder(expected)

	if err == nil {
		t.Errorf("err: expected != nil, gotten = nil")
	}
}

func TestSendMarketOrder_Direct_Buy_BaseAmt_AboveMaxBase(t *testing.T) {
	old_symbols := symbols
	old_get_spot_limits := GetSpotMarketLimits
	defer func() {
		symbols = old_symbols
		GetSpotMarketLimits = old_get_spot_limits
	}()

	symbols = get_symbols()
	GetSpotMarketLimits = func(symbol string) (model.SpotMarketLimits, error) {
		return model.SpotMarketLimits{
			MinBase:  utils.DecimalFromString("0.00000001"),
			MaxBase:  utils.DecimalFromString("99.9"),
			StepBase: utils.DecimalFromString("0.00000001"),
			MinQuote: utils.DecimalFromString("0.00000001")}, nil
	}

	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("3746.34")
	expected := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	_, err := SendSpotMarketOrder(expected)

	if err == nil {
		t.Errorf("err: expected != nil, gotten = nil")
	}
}

func TestSendMarketOrder_Direct_Sell_QuoteAmt(t *testing.T) {
	old_symbols := symbols
	old_create_order := binance_create_order
	old_get_spot_limits := GetSpotMarketLimits
	defer func() {
		symbols = old_symbols
		binance_create_order = old_create_order
		GetSpotMarketLimits = old_get_spot_limits
	}()

	symbols = get_symbols()
	binance_create_order = func(*binanceapi.CreateOrderService) (*binanceapi.CreateOrderResponse, error) {
		return &binanceapi.CreateOrderResponse{
			Symbol: "BTCUSDT",
			Side:   binanceapi.SideTypeBuy,
			Status: binanceapi.OrderStatusTypeFilled}, nil
	}
	GetSpotMarketLimits = func(symbol string) (model.SpotMarketLimits, error) {
		return model.SpotMarketLimits{
			MinBase:  utils.DecimalFromString("0.00000001"),
			MaxBase:  utils.DecimalFromString("99999999"),
			StepBase: utils.DecimalFromString("0.00000001"),
			MinQuote: utils.DecimalFromString("0.00000001")}, nil
	}

	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("3746.34")
	expected := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	gotten, err := SendSpotMarketOrder(expected)
	expected.Timestamp = gotten.Timestamp

	if err != nil {
		t.Errorf("err: expected = nil, gotten = %v", err)
	}
	testutils.AssertStructEq(t, expected, gotten)
}

func TestSendMarketOrder_Direct_Sell_QuoteAmt_BelowMinQuote(t *testing.T) {
	old_symbols := symbols
	old_get_spot_limits := GetSpotMarketLimits
	defer func() {
		symbols = old_symbols
		GetSpotMarketLimits = old_get_spot_limits
	}()

	symbols = get_symbols()
	GetSpotMarketLimits = func(symbol string) (model.SpotMarketLimits, error) {
		return model.SpotMarketLimits{
			MinBase:  utils.DecimalFromString("0.00000001"),
			MaxBase:  utils.DecimalFromString("99999999"),
			StepBase: utils.DecimalFromString("0.00000001"),
			MinQuote: utils.DecimalFromString("100.1")}, nil
	}

	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("3746.34")
	expected := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	_, err := SendSpotMarketOrder(expected)

	if err == nil {
		t.Errorf("err: expected != nil, gotten = nil")
	}
}

func TestSendMarketOrder_Indirect_Buy_BaseAmt(t *testing.T) {
	old_symbols := symbols
	old_create_order := binance_create_order
	old_get_spot_limits := GetSpotMarketLimits
	defer func() {
		symbols = old_symbols
		binance_create_order = old_create_order
		GetSpotMarketLimits = old_get_spot_limits
	}()

	symbols = get_symbols()
	binance_create_order = func(*binanceapi.CreateOrderService) (*binanceapi.CreateOrderResponse, error) {
		return &binanceapi.CreateOrderResponse{
			Symbol: "BTCUSDT",
			Side:   binanceapi.SideTypeBuy,
			Status: binanceapi.OrderStatusTypeFilled}, nil
	}
	GetSpotMarketLimits = func(symbol string) (model.SpotMarketLimits, error) {
		return model.SpotMarketLimits{
			MinBase:  utils.DecimalFromString("0.00000001"),
			MaxBase:  utils.DecimalFromString("99999999"),
			StepBase: utils.DecimalFromString("0.00000001"),
			MinQuote: utils.DecimalFromString("0.00000001")}, nil
	}

	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("3746.34")
	expected := get_operation_test(amt, model.BASE_AMOUNT, "USDT", "BTC", model.BUY, price)
	gotten, err := SendSpotMarketOrder(expected)
	expected.Timestamp = gotten.Timestamp
	expected.Base = "BTC"
	expected.Quote = "USDT"
	expected.Side = model.SELL
	expected.AmountSide = model.QUOTE_AMOUNT

	if err != nil {
		t.Errorf("err: expected = nil, gotten = %v", err)
	}
	testutils.AssertStructEq(t, expected, gotten)
}

func TestSendMarketOrder_Indirect_Sell_QuoteAmt(t *testing.T) {
	old_symbols := symbols
	old_create_order := binance_create_order
	old_get_spot_limits := GetSpotMarketLimits
	defer func() {
		symbols = old_symbols
		binance_create_order = old_create_order
		GetSpotMarketLimits = old_get_spot_limits
	}()

	symbols = get_symbols()
	binance_create_order = func(*binanceapi.CreateOrderService) (*binanceapi.CreateOrderResponse, error) {
		return &binanceapi.CreateOrderResponse{
			Symbol: "BTCUSDT",
			Side:   binanceapi.SideTypeBuy,
			Status: binanceapi.OrderStatusTypeFilled}, nil
	}
	GetSpotMarketLimits = func(symbol string) (model.SpotMarketLimits, error) {
		return model.SpotMarketLimits{
			MinBase:  utils.DecimalFromString("0.00000001"),
			MaxBase:  utils.DecimalFromString("99999999"),
			StepBase: utils.DecimalFromString("0.00000001"),
			MinQuote: utils.DecimalFromString("0.00000001")}, nil
	}

	amt := utils.DecimalFromString("100")
	price := utils.DecimalFromString("3746.34")
	expected := get_operation_test(amt, model.QUOTE_AMOUNT, "USDT", "BTC", model.SELL, price)
	gotten, err := SendSpotMarketOrder(expected)
	expected.Timestamp = gotten.Timestamp
	expected.Base = "BTC"
	expected.Quote = "USDT"
	expected.Side = model.BUY
	expected.AmountSide = model.BASE_AMOUNT

	if err != nil {
		t.Errorf("err: expected = nil, gotten = %v", err)
	}
	testutils.AssertStructEq(t, expected, gotten)
}

/************************* Helpers ***************************/

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

func get_remote_binance_account() *binanceapi.Account {
	return &binanceapi.Account{
		MakerCommission:  1000,
		TakerCommission:  2000,
		BuyerCommission:  3000,
		SellerCommission: 0,
		CanTrade:         true,
		CanWithdraw:      true,
		CanDeposit:       false,
		UpdateTime:       100,
		AccountType:      "acctype",
		Balances: []binanceapi.Balance{
			{Asset: "BTC", Free: "12.13", Locked: "11.11"},
			{Asset: "ETH", Free: "0.12", Locked: "122.56"},
			{Asset: "DOT", Free: "12.13", Locked: "1"},
			{Asset: "USDT", Free: "10900", Locked: "0"}},
	}
}

func get_remote_account() model.RemoteAccount {
	return model.RemoteAccount{
		MakerCommission:  1000,
		TakerCommission:  2000,
		BuyerCommission:  3000,
		SellerCommission: 0,
		Balances: []model.RemoteBalance{
			{Asset: "BTC", Amount: utils.DecimalFromString("12.13")},
			{Asset: "ETH", Amount: utils.DecimalFromString("0.12")},
			{Asset: "DOT", Amount: utils.DecimalFromString("12.13")},
			{Asset: "USDT", Amount: utils.DecimalFromString("10900")}}}
}
