package analytics

import (
	"context"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insert_many(anas []model.IAnalytics) error {
	collection := mongodb.GetAnalyticsCol()

	var payload []interface{}
	for i := range anas {
		payload = append(payload, anas[i])
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err := collection.InsertMany(context.TODO(), payload, opts)
	return err
}
