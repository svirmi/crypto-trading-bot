package strategy

import (
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/errors"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/dts"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/epts"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/pts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

func ValidateStrategyConfig(strategyType model.StrategyType, props map[string]string) errors.CtbError {
	if strategyType == model.DTS_STRATEGY {
		return dts.LocalAccountDTS{}.ValidateConfig(props)
	} else if strategyType == model.PTS_STRATEGY {
		return pts.LocalAccountPTS{}.ValidateConfig(props)
	} else if strategyType == model.EPTS_STRATEGY {
		return epts.LocalAccountEPTS{}.ValidateConfig(props)
	} else {
		err := errors.Internal(logger.STR_ERR_UNKNOWN_STRATEGY, strategyType)
		logrus.Error(err.Error())
		return err
	}
}

func InstanciateLocalAccount(strategyType model.StrategyType) (model.ILocalAccount, errors.CtbError) {
	if strategyType == model.DTS_STRATEGY {
		return dts.LocalAccountDTS{}, nil
	} else if strategyType == model.PTS_STRATEGY {
		return pts.LocalAccountPTS{}, nil
	} else if strategyType == model.EPTS_STRATEGY {
		return epts.LocalAccountEPTS{}, nil
	} else {
		err := errors.Internal(logger.STR_ERR_UNKNOWN_STRATEGY, strategyType)
		logrus.Error(err.Error())
		return nil, nil
	}
}

func DecodeLaccount(raw bson.Raw, registry *bsoncodec.Registry) (model.ILocalAccount, errors.CtbError) {
	payload := struct {
		model.LocalAccountMetadata `bson:"metadata"`
	}{}

	err := bson.Unmarshal(raw, &payload)
	if err != nil {
		return nil, errors.WrapInternal(err)
	}

	strategyType := model.StrategyType(payload.GetStrategyType())
	if strategyType == model.DTS_STRATEGY {
		laccount_dts := dts.LocalAccountDTS{}
		err := bson.UnmarshalWithRegistry(registry, raw, &laccount_dts)
		return laccount_dts, errors.WrapInternal(err)
	} else if strategyType == model.PTS_STRATEGY {
		laccount_pts := pts.LocalAccountPTS{}
		err := bson.UnmarshalWithRegistry(registry, raw, &laccount_pts)
		return laccount_pts, errors.WrapInternal(err)
	} else if strategyType == model.EPTS_STRATEGY {
		laccount_epts := epts.LocalAccountEPTS{}
		err := bson.UnmarshalWithRegistry(registry, raw, &laccount_epts)
		return laccount_epts, errors.WrapInternal(err)
	} else {
		err := errors.Internal(logger.STR_ERR_UNKNOWN_STRATEGY, strategyType)
		logrus.Error(err.Error())
		return nil, err
	}
}
