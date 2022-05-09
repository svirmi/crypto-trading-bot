package local

import "github.com/valerioferretti92/crypto-trading-bot/internal/model"

type local_exchange struct{}

func GetExchange() model.IExchange {
	return local_exchange{}
}

func (be local_exchange) Initialize(mmsChannel chan []model.MiniMarketStats) error {
	// To be implented
	return nil
}

func (be local_exchange) CanSpotTrade(symbol string) bool {
	// To be implemented
	return true
}

func (be local_exchange) GetSpotMarketLimits(symbol string) (model.SpotMarketLimits, error) {
	// To be implemented
	return model.SpotMarketLimits{}, nil
}

func (be local_exchange) FilterTradableAssets(bases []string) []string {
	// To be implemented
	return nil
}

func (be local_exchange) GetAssetsValue(bases []string) (map[string]model.AssetPrice, error) {
	// To be implemented
	return nil, nil
}

func (be local_exchange) GetAccout() (model.RemoteAccount, error) {
	// To be implemented
	return model.RemoteAccount{}, nil
}

func (be local_exchange) SendSpotMarketOrder(op model.Operation) (model.Operation, error) {
	// To be implemented
	return model.Operation{}, nil
}

func (be local_exchange) MiniMarketsStatsServe(assets []string) error {
	// To be implemented
	return nil
}

func (be local_exchange) MiniMarketsStatsStop() {
	// To be implemented
}
