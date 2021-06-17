package repository

import (
	"context"
	"log"

	binanceapi "github.com/adshao/go-binance/v2"
)

func Insert(miniMarketsStat binanceapi.WsAllMiniMarketsStatEvent) {
	objs := toInterfaces(miniMarketsStat)

	result, err := miniMarketsStatCol.InsertMany(context.TODO(), objs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Inserted %d items into %s", len(result.InsertedIDs), miniMarketsStatColName)
}

func toInterfaces(miniMarketsStat binanceapi.WsAllMiniMarketsStatEvent) []interface{} {
	objs := make([]interface{}, 0, len(miniMarketsStat))
	for _, obj := range miniMarketsStat {
		objs = append(objs, *obj)
	}
	return objs
}
