package strategy

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/dts"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/pts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

func DecodeLaccount(raw bson.Raw, registry *bsoncodec.Registry) (model.ILocalAccount, error) {
	payload := struct {
		model.LocalAccountMetadata `bson:"metadata"`
	}{}

	err := bson.Unmarshal(raw, &payload)
	if err != nil {
		return nil, err
	}

	strategyType := model.StrategyType(payload.GetStrategyType())
	if strategyType == model.DTS_STRATEGY {
		laccount_dts := dts.LocalAccountDTS{}
		err := bson.UnmarshalWithRegistry(registry, raw, &laccount_dts)
		return laccount_dts, err
	} else if strategyType == model.PTS_STRATEGY {
		laccount_pts := pts.LocalAccountPTS{}
		err := bson.UnmarshalWithRegistry(registry, raw, &laccount_pts)
		return laccount_pts, err
	} else {
		err := fmt.Errorf(logger.STR_ERR_UNKNOWN_STRATEGY, strategyType)
		logrus.Error(err.Error())
		return nil, err
	}
}
