package api

import (
	"net/http"

	"github.com/valerioferretti92/crypto-trading-bot/internal/exchange"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy"
)

func create_execution(req exe_create_req_dto) (ctb_response_dto, ctb_error_dto) {
	strategyType, err := model.ParseStr(req.StrategyType)
	if err != nil {
		return ctb_response_dto{}, to_ctb_error_dto(err)
	}

	err = strategy.ValidateStrategyConfig(strategyType, req.StrategyConfig)
	if err != nil {
		return ctb_response_dto{}, to_ctb_error_dto(err)
	}

	racc, err := exchange.GetAccount()
	if err != nil {
		return ctb_response_dto{}, to_ctb_error_dto(err)
	}

	exeReq := model.ExecutionInit{
		Raccount:     racc,
		StrategyType: strategyType,
		Props:        req.StrategyConfig}
	exe, err := executions.Create(exeReq)
	if err != nil {
		return ctb_response_dto{}, to_ctb_error_dto(err)
	}

	tradableAssets := exchange.FilterTradableAssets(exe.Assets)
	assetPrices, err := exchange.GetAssetsValue(tradableAssets)
	if err != nil {
		return ctb_response_dto{}, to_ctb_error_dto(err)
	}

	laccReq := model.LocalAccountInit{
		ExeId:               exe.ExeId,
		RAccount:            racc,
		StrategyType:        strategyType,
		TradableAssetsPrice: assetPrices}
	_, err = laccount.Create(laccReq)
	if err != nil {
		return ctb_response_dto{}, to_ctb_error_dto(err)
	}
	return exe_to_ctb_response_dto(http.StatusCreated, exe), ctb_error_dto{}
}

func update_execution(exeId string, req exe_update_req_dto) (ctb_response_dto, ctb_error_dto) {
	update := model.Execution{
		ExeId:  exeId,
		Status: model.ExeStatus(req.Status),
	}

	exe, err := executions.Update(update)
	if err != nil {
		return ctb_response_dto{}, to_ctb_error_dto(err)
	}

	return exe_to_ctb_response_dto(http.StatusOK, exe), ctb_error_dto{}
}
