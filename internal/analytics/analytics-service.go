package analytics

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/config"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/operations"
	"github.com/valerioferretti92/crypto-trading-bot/internal/prices"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
)

func StoreAnalytics(exeId string) {
	// Collecting execution data
	exes, err := executions.GetByExeId(exeId)
	if err != nil {
		logrus.Errorf(logger.ANAL_ERR_FAILED_TO_GENERATE, err.Error())
		return
	}
	if len(exes) < 2 {
		logrus.Errorf(logger.ANAL_ERR_CURRENTLY_ACTIVE, exeId)
		return
	}
	exe1 := exes[0]
	exe2 := exes[len(exes)-1]
	if exe1.Status != model.EXE_ACTIVE {
		logrus.Errorf(logger.ANAL_ERR_BAD_EXE_STATUS,
			exeId, model.EXE_ACTIVE, exe1.Status)
		return
	}
	if exe2.Status != model.EXE_TERMINATED {
		logrus.Errorf(logger.ANAL_ERR_BAD_EXE_STATUS,
			exeId, model.EXE_TERMINATED, exe2.Status)
		return
	}
	assets := exe1.Assets
	timestamp1 := exe1.Timestamp
	timestamp2 := exe2.Timestamp
	symbols, _ := utils.GetSymbolsFromAssets(assets)

	// Collecting sprices
	sprices, err := prices.GetByTimestamp(symbols, timestamp1, timestamp2)
	if err != nil {
		logrus.Errorf(logger.ANAL_ERR_FAILED_TO_GENERATE, err.Error())
		return
	}
	if len(sprices) == 0 {
		logrus.Errorf(logger.ANAL_ERR_NO_PRICES, timestamp1, timestamp2)
		return
	}

	// Collecting wallet data
	laccs, err := laccount.GetByExeId(exeId)
	if err != nil {
		logrus.Errorf(logger.ANAL_ERR_FAILED_TO_GENERATE, err.Error())
		return
	}
	if len(laccs) == 0 {
		logrus.Errorf(logger.ANAL_ERR_NO_LACCS, exeId)
		return
	}

	// Collection operations
	ops, err := operations.GetByExeId(exeId)
	if err != nil {
		err := fmt.Errorf(logger.ANAL_ERR_FAILED_TO_GENERATE, err.Error())
		logrus.Error(err.Error())
		return
	}

	// Declaring slice to temporarily holds analystics in memory
	anals := make([]model.IAnalytics, 0, len(sprices)+len(ops)+1)

	// Building execution analytics
	logrus.Infof(logger.ANAL_BUILDING_EXE, exeId)
	strategyConfig := config.GetStrategyConfig()
	anals = append(anals, build_exe_analytics(exe1, strategyConfig))

	// Building wallet analytics
	logrus.Infof(logger.ANAL_BUILDING_WALLETS, exeId, len(laccs), len(sprices))
	wanal := init_wallet_analytics(exeId, assets)
	var priceIdx, laccIdx int = 0, 0
	for {
		laccount := laccs[laccIdx]
		spricesByTs := sprices[priceIdx]
		wanal, err = build_wallet_analytics(laccount, spricesByTs, wanal)
		if err == nil {
			anals = append(anals, wanal)
		} else {
			logrus.Errorf(logger.ANAL_ERR_SKIP_ANALYTICS, err.Error())
		}

		if priceIdx == len(sprices)-1 && laccIdx == len(laccs)-1 {
			break
		}

		if priceIdx < len(sprices)-1 {
			priceIdx = priceIdx + 1
		}

		if laccIdx == len(laccs)-1 {
			continue
		}
		if laccs[laccIdx+1].GetTimestamp() < sprices[priceIdx].Timestamp {
			laccIdx = laccIdx + 1
		}
		if priceIdx == len(sprices)-1 && laccIdx < len(laccs)-1 {
			laccIdx = laccIdx + 1
		}
	}

	// Build operation analytics
	logrus.Infof(logger.ANAL_BUILDING_OPS, exeId, len(ops))
	for _, op := range ops {
		anals = append(anals, build_op_analytics(op))
	}

	// Storing analytics
	logrus.Infof(logger.ANAL_STORE_ANALYTICS, len(anals))
	insert_many(anals)
}

func init_wallet_analytics(exeId string, assets []string) model.WalletAnalytics {
	wanal := model.WalletAnalytics{
		ExeId:         exeId,
		AnalyticsType: model.WALLET_ANALYTICS,
		Timestamp:     0,
		AssetStatuses: make(map[string]model.AssetStatus, len(assets)),
		WalletValue:   decimal.Zero}

	for _, asset := range assets {
		wanal.AssetStatuses[asset] = model.AssetStatus{
			Asset:  asset,
			Price:  decimal.Zero,
			Amount: decimal.Zero}
	}
	return wanal
}

func build_wallet_analytics(lacc model.ILocalAccount, spricesByTs model.SymbolPriceByTimestamp, wanal model.WalletAnalytics) (model.WalletAnalytics, error) {
	// Checking execution id
	if wanal.ExeId != lacc.GetExeId() {
		err := fmt.Errorf(logger.ANAL_ERR_MISMATCHING_EXE_IDS, wanal.ExeId, lacc.GetExeId())
		logrus.Errorf(err.Error())
		return wanal, err
	}

	// Allocating asset analytics map
	old := wanal.AssetStatuses
	wanal.AssetStatuses = make(map[string]model.AssetStatus, len(old))
	for k, v := range old {
		wanal.AssetStatuses[k] = v
	}

	// Applying price updates to asset analytics
	assetStatuses := lacc.GetAssetAmounts()
	for _, symbolPrice := range spricesByTs.SymbolPrices {
		asset, err := utils.GetAssetFromSymbol(symbolPrice.Symbol)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}

		assetStatus, found := assetStatuses[asset]
		if !found || assetStatus.IsEmpty() {
			err := fmt.Errorf(logger.ANAL_ERR_ASSET_NOT_FOUND, asset)
			logrus.Errorf(err.Error())
			return wanal, err
		}
		wanal.AssetStatuses[asset] = model.AssetStatus{
			Asset:  asset,
			Price:  symbolPrice.Price,
			Amount: assetStatus.Amount}
	}

	// Updating USDT analytics
	var usdt decimal.Decimal = decimal.Zero
	usdtStatus, found := assetStatuses["USDT"]
	if found && !usdtStatus.IsEmpty() {
		usdt = usdtStatus.Amount
	}
	wanal.AssetStatuses["USDT"] = model.AssetStatus{
		Asset:  "USDT",
		Price:  utils.DecimalFromString("1"),
		Amount: usdt}

	// Updating wallet value
	var walletValue decimal.Decimal = decimal.Zero
	for _, assetAanalytics := range wanal.AssetStatuses {
		value := assetAanalytics.Price.Mul(assetAanalytics.Amount)
		walletValue = walletValue.Add(value)
	}
	wanal.WalletValue = walletValue

	// Updating timestamp
	wanal.Timestamp = lacc.GetTimestamp()

	return wanal, nil
}

func build_op_analytics(op model.Operation) model.OpAnalytics {
	return model.OpAnalytics{
		ExeId:         op.ExeId,
		AnalyticsType: model.OP_ANALYTICS,
		Timestamp:     op.Timestamp,
		Base:          op.Base,
		Quote:         op.Quote,
		Amount:        op.Amount,
		Side:          op.Side,
		AmountSide:    op.AmountSide,
		Price:         op.Price}
}

func build_exe_analytics(exe model.Execution, strategyConfig config.StrategyConfig) model.ExeAnalytics {
	stype := model.StrategyType(strategyConfig.Type)

	sconfig := strategyConfig.Config
	var props map[string]string
	err := mapstructure.Decode(sconfig, &props)
	if err != nil {
		logrus.Error(err.Error())
	}

	return model.ExeAnalytics{
		ExeId:         exe.ExeId,
		AnalyticsType: model.EXE_ANALYTICS,
		Timestamp:     exe.Timestamp,
		Assets:        exe.Assets,
		Status:        model.EXE_TERMINATED,
		StrategyType:  stype,
		Props:         props}
}
