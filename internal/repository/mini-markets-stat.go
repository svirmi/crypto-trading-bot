package repository

import (
	"context"
	"fmt"
	"log"

	binanceapi "github.com/adshao/go-binance/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpsertMiniMarketsStat(miniMarketsStat binanceapi.WsAllMiniMarketsStatEvent) error {
	models := make([]mongo.WriteModel, 0, len(miniMarketsStat))
	for _, miniMarketStat := range miniMarketsStat {
		doc, err := wsMiniMarketsStatEventToDoc(miniMarketStat)
		if err != nil {
			log.Printf("%s: skipping %s price update\n", err.Error(), miniMarketStat.Symbol)
			continue
		}

		update := bson.D{{"$set", doc}}
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.D{{"symbol", miniMarketStat.Symbol}}).
			SetUpdate(update).
			SetUpsert(true)
		models = append(models, model)
	}
	opts := options.BulkWrite().SetOrdered(false)
	res, err := miniMarketsStatCol.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return fmt.Errorf("bulk price update failed")
	}

	log.Printf("%d price updates saved into %s", res.ModifiedCount, miniMarketsStatColName)
	return nil
}

func FindSymbolByPair(symbol string) (binanceapi.WsMiniMarketsStatEvent, error) {
	var miniMarketStat binanceapi.WsMiniMarketsStatEvent
	err := miniMarketsStatCol.FindOne(context.TODO(), bson.D{{"symbol", symbol}}).
		Decode(&miniMarketStat)
	if err != nil {
		log.Printf("%s\n", err.Error())
		err = fmt.Errorf("failed to retrieve %s's lastest price", symbol)
		return binanceapi.WsMiniMarketsStatEvent{}, err
	}
	return miniMarketStat, nil
}

func wsMiniMarketsStatEventToDoc(miniMarketStat *binanceapi.WsMiniMarketsStatEvent) (doc *bson.D, err error) {
	data, err := bson.Marshal(miniMarketStat)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
