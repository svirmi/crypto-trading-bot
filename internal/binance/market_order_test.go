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

	got := FilterTradableAssets(assets)
	exp := []string{"BTC", "ETH"}

	testutils.AssertEq(t, exp, got, "tradable_assets")
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

	got, err := GetAssetsValue([]string{"BTC", "ETH", "DOT", "BNB", "TRX"})
	testutils.AssertNil(t, err, "err")

	exp := make(map[string]model.AssetPrice)
	exp["BTC"] = model.AssetPrice{Asset: "BTC", Price: utils.DecimalFromString("35998.34")}
	exp["ETH"] = model.AssetPrice{Asset: "ETH", Price: utils.DecimalFromString("44978.12")}
	exp["DOT"] = model.AssetPrice{Asset: "DOT", Price: utils.DecimalFromString("98.12")}

	testutils.AssertEq(t, exp, got, "asset_value")
}

func TestGetAccount(t *testing.T) {
	old := binance_get_account
	defer func() {
		binance_get_account = old
	}()

	binance_get_account = func(*binanceapi.GetAccountService) (*binanceapi.Account, error) {
		raccount := get_remote_binance_account()
		raccount.Balances = append(raccount.Balances, binanceapi.Balance{
			Asset:  "SHIBA",
			Free:   "0",
			Locked: "100"})
		return get_remote_binance_account(), nil
	}

	got, err := GetAccout()
	testutils.AssertNil(t, err, "err")
	exp := get_remote_account()

	testutils.AssertEq(t, exp, got, "remote_account")
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

	got, err := SendSpotMarketOrder(op)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "operation")
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
	exp := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	got, err := SendSpotMarketOrder(exp)
	exp.Timestamp = got.Timestamp

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "operation")
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
	exp := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	_, err := SendSpotMarketOrder(exp)

	testutils.AssertNotNil(t, err, "err")
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
	exp := get_operation_test(amt, model.BASE_AMOUNT, "SHIBA", "USDT", model.BUY, price)
	_, err := SendSpotMarketOrder(exp)

	testutils.AssertNotNil(t, err, "err")
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
	exp := get_operation_test(amt, model.BASE_AMOUNT, "BTC", "USDT", model.BUY, price)
	_, err := SendSpotMarketOrder(exp)

	testutils.AssertNotNil(t, err, "err")
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
	exp := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	got, err := SendSpotMarketOrder(exp)
	exp.Timestamp = got.Timestamp

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "operation")
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
	exp := get_operation_test(amt, model.QUOTE_AMOUNT, "BTC", "USDT", model.SELL, price)
	_, err := SendSpotMarketOrder(exp)

	testutils.AssertNotNil(t, err, "err")
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
	exp := get_operation_test(amt, model.BASE_AMOUNT, "USDT", "BTC", model.BUY, price)
	got, err := SendSpotMarketOrder(exp)
	exp.Timestamp = got.Timestamp
	exp.Base = "BTC"
	exp.Quote = "USDT"
	exp.Side = model.SELL
	exp.AmountSide = model.QUOTE_AMOUNT

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "operation")
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
	exp := get_operation_test(amt, model.QUOTE_AMOUNT, "USDT", "BTC", model.SELL, price)
	got, err := SendSpotMarketOrder(exp)
	exp.Timestamp = got.Timestamp
	exp.Base = "BTC"
	exp.Quote = "USDT"
	exp.Side = model.BUY
	exp.AmountSide = model.BASE_AMOUNT

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "operation")
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
