package pts

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

type AssetStatusPTS struct {
	Asset              string          `bson:"asset"`              // Asset being tracked
	Amount             decimal.Decimal `bson:"amount"`             // Amount of that asset currently owned
	LastOperationPrice decimal.Decimal `bson:"lastOperationPrice"` // Asset value at the time last op was executed
}

func (a AssetStatusPTS) IsEmpty() bool {
	return reflect.DeepEqual(a, AssetStatusPTS{})
}

type LocalAccountPTS struct {
	model.LocalAccountMetadata `bson:"metadata"`
	Ignored                    map[string]decimal.Decimal `bson:"ignored"` // Ignored assets
	Assets                     map[string]AssetStatusPTS  `bson:"assets"`  // Value allocation across assets
	Usdt                       decimal.Decimal            `bson:"usdt"`    // Usdt available to invested
}

func (a LocalAccountPTS) IsEmpty() bool {
	return reflect.DeepEqual(a, LocalAccountPTS{})
}

const (
	_NP_OP         = "NO_OP"
	_PTS_BUY_DESC  = "pts_buy"
	_PTS_SELL_DESC = "pts_sell"
)

func (a LocalAccountPTS) Initialize(req model.LocalAccountInit) (model.ILocalAccount, error) {
	var ignored = make(map[string]decimal.Decimal)
	var assets = make(map[string]AssetStatusPTS)
	var usdt = decimal.Zero

	for _, rbalance := range req.RAccount.Balances {
		if rbalance.Asset == "USDT" {
			usdt = rbalance.Amount
			continue
		}

		price, found := req.TradableAssetsPrice[rbalance.Asset]
		if !found {
			logrus.WithField("comp", "pts").
				Warnf(logger.XXX_IGNORED_ASSET, rbalance.Asset)
			ignored[rbalance.Asset] = rbalance.Amount
			continue
		}
		if decimal.Zero.Equals(rbalance.Amount) {
			continue
		}
		assetStatus := AssetStatusPTS{
			Asset:              rbalance.Asset,
			Amount:             rbalance.Amount,
			LastOperationPrice: price.Price}
		assets[rbalance.Asset] = assetStatus
	}

	return LocalAccountPTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        req.ExeId,
			StrategyType: req.StrategyType,
			Timestamp:    time.Now().UnixMicro()},
		Ignored: ignored,
		Assets:  assets,
		Usdt:    usdt}, nil
}

func (a LocalAccountPTS) RegisterTrading(op model.Operation) (model.ILocalAccount, error) {
	// Check execution ids
	if op.ExeId != a.ExeId {
		logrus.WithField("comp", "pts").
			Panicf(logger.XXX_ERR_MISMATCHING_EXE_IDTS, a.ExeId, op.ExeId)
	}

	// If the result status is failed, NOP
	if op.Status == model.FAILED {
		err := fmt.Errorf(logger.XXX_ERR_FAILED_OP, op.OpId)
		logrus.WithField("comp", "pts").Error(err.Error())
		return a, err
	}

	// DTS only handle operation back and forth USDT
	if op.Quote != "USDT" {
		err := fmt.Errorf(logger.PTS_ERR_BAD_QUOTE_CURRENCY, op.Quote)
		logrus.WithField("comp", "pts").Error(err.Error())
		return a, err
	}

	// Getting asset status
	assetStatus, found := a.Assets[op.Base]
	if !found {
		err := fmt.Errorf(logger.XXX_ERR_ASSET_NOT_FOUND, op.Base)
		logrus.WithField("comp", "pts").Error(err.Error())
		return a, err
	}

	// Updating asset status
	currentAmntUsdt := a.Usdt
	baseAmount := op.Results.BaseDiff
	quoteAmount := op.Results.QuoteDiff
	if op.Side == model.BUY {
		assetStatus.Amount = assetStatus.Amount.Add(baseAmount).Round(8)
		currentAmntUsdt = currentAmntUsdt.Sub(quoteAmount).Round(8)
	} else if op.Side == model.SELL {
		assetStatus.Amount = assetStatus.Amount.Sub(baseAmount).Round(8)
		currentAmntUsdt = currentAmntUsdt.Add(quoteAmount).Round(8)
	} else {
		err := fmt.Errorf(logger.XXX_ERR_UNKNWON_OP_TYPE, op.Type)
		logrus.WithField("comp", "pts").Error(err.Error())
		return a, err
	}
	if assetStatus.Amount.LessThan(decimal.Zero) {
		logrus.WithField("comp", "pts").
			Panicf(logger.XXX_ERR_NEGATIVE_BALANCE, assetStatus.Asset, assetStatus.Amount)
	}
	if currentAmntUsdt.LessThan(decimal.Zero) {
		logrus.WithField("comp", "pts").
			Panicf(logger.XXX_ERR_NEGATIVE_BALANCE, "USDT", currentAmntUsdt)
	}
	assetStatus.LastOperationPrice = op.Results.ActualPrice

	// Returning results
	a.Assets[op.Base] = assetStatus
	a.Usdt = currentAmntUsdt
	a.Timestamp = time.Now().UnixMicro()
	a.AccountId = uuid.NewString()
	return a, nil
}

func (a LocalAccountPTS) GetOperation(mms model.MiniMarketStats, slimts model.SpotMarketLimits) (model.Operation, error) {
	asset := mms.Asset
	assetStatus, found := a.Assets[asset]
	if !found {
		err := fmt.Errorf(logger.XXX_ERR_ASSET_NOT_FOUND, asset)
		logrus.WithField("comp", "pts").Error(err.Error())
		return model.Operation{}, err
	}

	lastOpPrice := assetStatus.LastOperationPrice
	currentAmnt := assetStatus.Amount
	currentAmntUsdt := a.Usdt
	currentPrice := mms.LastPrice
	if currentPrice.Equals(decimal.Zero) {
		err := fmt.Errorf(logger.XXX_ERR_ZERO_EXP_PRICE, asset)
		logrus.WithField("comp", "pts").Errorf(err.Error())
		return model.Operation{}, err
	}

	config := get_pts_config(config.GetStrategyConfig())
	sellPrice := utils.IncrementByPercentage(lastOpPrice, config.SellPercentage)
	buyPrice := utils.IncrementByPercentage(lastOpPrice, utils.SignChangeDecimal(config.BuyPercentage))
	sellAmnt := utils.PercentageOf(currentAmnt, config.SellAmountPercentage)
	buyAmnt := utils.PercentageOf(currentAmntUsdt, config.BuyAmountPercentage)

	var op model.Operation
	if currentPrice.GreaterThanOrEqual(sellPrice) {
		op = build_sell_op(a.ExeId, asset, _PTS_SELL_DESC, sellAmnt, currentPrice)
	} else if currentPrice.LessThanOrEqual(buyPrice) {
		op = build_buy_op(a.ExeId, asset, _PTS_BUY_DESC, buyAmnt, currentPrice)
	} else {
		logrus.WithField("comp", "pts").
			Tracef(logger.PTS_TRADE, _NP_OP, asset, lastOpPrice, currentPrice)
		return model.Operation{}, nil
	}

	err := check_spot_market_limits(op, slimts)
	if err != nil {
		return model.Operation{}, err
	}

	logrus.WithField("comp", "pts").
		Infof(logger.PTS_TRADE, op.Cause, asset, lastOpPrice, currentPrice)
	return op, nil
}

func check_spot_market_limits(op model.Operation, slimits model.SpotMarketLimits) error {
	if op.AmountSide == model.QUOTE_AMOUNT && op.Amount.LessThan(slimits.MinQuote) {
		err := fmt.Errorf(logger.XXX_BELOW_QUOTE_LIMIT,
			op.Base+op.Quote, op.Side, op.Amount, op.AmountSide, slimits.MinQuote)
		logrus.WithField("comp", "dts").Error(err.Error())
		return err
	}
	if op.AmountSide == model.BASE_AMOUNT && op.Amount.LessThan(slimits.MinBase) {
		err := fmt.Errorf(logger.XXX_BELOW_BASE_LIMIT,
			op.Base+op.Quote, op.Side, op.Amount, op.AmountSide, slimits.MinBase)
		logrus.WithField("comp", "dts").Error(err.Error())
		return err
	}
	// No checks on MaxBase as big orders should be broken down into
	// smaller orders by the exchange package
	return nil
}

func build_buy_op(exeId, asset, cause string, amnt, targetPrice decimal.Decimal) model.Operation {
	return model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      exeId,
		Type:       model.AUTO,
		Base:       asset,
		Quote:      "USDT",
		Side:       model.BUY,
		Amount:     amnt,
		AmountSide: model.QUOTE_AMOUNT,
		Price:      targetPrice,
		Cause:      cause,
		Status:     model.PENDING}
}

func build_sell_op(exeId, asset, cause string, amnt, targetPrice decimal.Decimal) model.Operation {
	return model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      exeId,
		Type:       model.AUTO,
		Base:       asset,
		Quote:      "USDT",
		Side:       model.SELL,
		Amount:     amnt,
		AmountSide: model.BASE_AMOUNT,
		Price:      targetPrice,
		Cause:      cause,
		Status:     model.PENDING}
}
