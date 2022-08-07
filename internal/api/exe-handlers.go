package api

import (
	"net/http"

	"github.com/valerioferretti92/crypto-trading-bot/internal/exchange"
	"github.com/valerioferretti92/crypto-trading-bot/internal/executions"
	"github.com/valerioferretti92/crypto-trading-bot/internal/laccount"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy"
)

func create_execution(req exe_create_req_dto) (int, exe_res_dto, error_dto) {
	strategyType, err := model.ParseStr(req.StrategyType)
	if err != nil {
		return bad_request(error_dto{err.Error()})
	}

	err = strategy.ValidateStrategyConfig(strategyType, req.StrategyConfig)
	if err != nil {
		return bad_request(error_dto{err.Error()})
	}

	racc, err := exchange.GetAccount()
	if err != nil {
		return internal_server_error(error_dto{err.Error()})
	}

	exeReq := model.ExecutionInit{
		Raccount:     racc,
		StrategyType: strategyType,
		Props:        req.StrategyConfig}
	exe, err := executions.Create(exeReq)
	if err != nil {
		return bad_request(error_dto{err.Error()})
	}

	tradableAssets := exchange.FilterTradableAssets(exe.Assets)
	assetPrices, err := exchange.GetAssetsValue(tradableAssets)
	if err != nil {
		return internal_server_error(error_dto{err.Error()})
	}

	laccReq := model.LocalAccountInit{
		ExeId:               exe.ExeId,
		RAccount:            racc,
		StrategyType:        strategyType,
		TradableAssetsPrice: assetPrices}
	_, err = laccount.Create(laccReq)
	if err != nil {
		return internal_server_error(error_dto{err.Error()})
	}
	return created(to_exe_res_dto(exe))
}

func update_execution(exeId string, req exe_update_req_dto) (int, exe_res_dto, error_dto) {
	update := model.Execution{
		ExeId:  exeId,
		Status: model.ExeStatus(req.Status),
	}

	exe, err := executions.Update(update)
	if err != nil {
		return bad_request(error_dto{err.Error()})
	}

	exeRes := to_exe_res_dto(exe)
	return ok(exeRes)
}

func ok(e exe_res_dto) (int, exe_res_dto, error_dto) {
	return http.StatusOK, e, error_dto{}
}

func created(e exe_res_dto) (int, exe_res_dto, error_dto) {
	return http.StatusCreated, e, error_dto{}
}

func bad_request(e error_dto) (int, exe_res_dto, error_dto) {
	return http.StatusBadRequest, exe_res_dto{}, e
}

func internal_server_error(e error_dto) (int, exe_res_dto, error_dto) {
	return http.StatusInternalServerError, exe_res_dto{}, e
}
