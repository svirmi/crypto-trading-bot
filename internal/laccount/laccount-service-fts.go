package laccount

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
)

func buildLocalAccountFTS(exeId string, raccount model.RemoteAccount) (model.LocalAccountFTS, error) {
	var frozenUsdt float32 = 0
	var assets = make(map[string]model.AssetStatusFTS)

	for _, rbalance := range raccount.Balances {
		amount, err := parseFloat32(rbalance.Amount)
		if err != nil {
			return model.LocalAccountFTS{}, err
		}
		if rbalance.Asset == "USDT" {
			frozenUsdt = amount
			continue
		}
		assetStatus, err := buildAssetStatusFTS(rbalance)
		if err != nil {
			return model.LocalAccountFTS{}, err
		}
		assets[rbalance.Asset] = assetStatus
	}

	return model.LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        exeId,
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMilli()},
		FrozenUsdt: frozenUsdt,
		Assets:     assets}, nil
}

func buildAssetStatusFTS(rbalance model.RemoteBalance) (model.AssetStatusFTS, error) {
	amount, err := parseFloat32(rbalance.Amount)
	if err != nil {
		return model.AssetStatusFTS{}, err
	}

	return model.AssetStatusFTS{
		Asset:             rbalance.Asset,
		Amount:            amount,
		Usdt:              0,
		LastOperationType: model.OP_BUY_FTS,
		LastOperationRate: 0}, nil
}

func parseFloat32(payload string) (float32, error) {
	value, err := strconv.ParseFloat(payload, 32)
	if err != nil {
		return 0, err
	}
	return float32(value), nil
}
