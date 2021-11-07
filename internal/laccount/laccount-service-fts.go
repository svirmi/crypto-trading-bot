package laccount

import (
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/trading-bot-demo/internal/model"
	"github.com/valerioferretti92/trading-bot-demo/internal/utils"
)

func buildLocalAccountFTS(creationRequest model.LocalAccountInit) (model.LocalAccountFTS, error) {
	var ignored = make(map[string]float32)
	var assets = make(map[string]model.AssetStatusFTS)

	for _, rbalance := range creationRequest.RAccount.Balances {
		symbol := utils.GetSymbolFromAsset(rbalance.Asset)
		price, found := creationRequest.TradableAssetsPrice[symbol]
		if !found {
			ignored[rbalance.Asset] = rbalance.Amount
			continue
		}
		assetStatus, err := buildAssetStatusFTS(rbalance, price)
		if err != nil {
			return model.LocalAccountFTS{}, err
		}
		assets[rbalance.Asset] = assetStatus
	}

	return model.LocalAccountFTS{
		LocalAccountMetadata: model.LocalAccountMetadata{
			AccountId:    uuid.NewString(),
			ExeId:        creationRequest.ExeId,
			StrategyType: model.FIXED_THRESHOLD_STRATEGY,
			Timestamp:    time.Now().UnixMilli()},
		Ignored: ignored,
		Assets:  assets}, nil
}

func buildAssetStatusFTS(rbalance model.RemoteBalance, price model.SymbolPrice) (model.AssetStatusFTS, error) {
	return model.AssetStatusFTS{
		Asset:             rbalance.Asset,
		Amount:            rbalance.Amount,
		Usdt:              0,
		LastOperationType: model.OP_BUY_FTS,
		LastOperationRate: price.Price}, nil
}
