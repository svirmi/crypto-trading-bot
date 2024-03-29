package laccount

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/dts"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/pts"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFindLatestByExeId_None(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	got, err := find_latest_by_exeId(uuid.NewString())

	testutils.AssertNil(t, err, "err")
	testutils.AssertNil(t, got, "laccount")
}

/****************************** DTS ********************************/

func TestInsert_DTS(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	exp := get_laccount_test_DTS()
	exeIds = append(exeIds, exp.ExeId)
	err := insert(exp)
	testutils.AssertNil(t, err, "err")

	got, err := find_latest_by_exeId(exp.ExeId)

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "laccount")
}

func TestFindLatestByExeId_DTS(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	exp := get_laccount_test_DTS()
	exeIds = append(exeIds, exp.ExeId)
	err := insert(exp)
	testutils.AssertNil(t, err, "err")

	exp.Assets["DOT"] = dts.AssetStatusDTS{
		Asset:              "DOT",
		Amount:             utils.DecimalFromString("55.56"),
		Usdt:               decimal.Zero,
		LastOperationType:  dts.OP_BUY_DTS,
		LastOperationPrice: utils.DecimalFromString("18.45")}
	exp.Timestamp = time.Now().UnixMicro()
	err = insert(exp)
	testutils.AssertNil(t, err, "err")

	exeIds = append(exeIds, exp.AccountId)
	got, err := find_latest_by_exeId(exp.ExeId)

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "laccount")
}

/****************************** PTS ********************************/

func TestInsert_PTS(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	exp := get_laccount_test_PTS()
	exeIds = append(exeIds, exp.ExeId)
	err := insert(exp)
	testutils.AssertNil(t, err, "err")

	got, err := find_latest_by_exeId(exp.ExeId)

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "laccount")
}

func TestFindLatestByExeId_PTS(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	exp := get_laccount_test_PTS()
	exeIds = append(exeIds, exp.ExeId)
	err := insert(exp)
	testutils.AssertNil(t, err, "err")

	exp.Assets["DOT"] = pts.AssetStatusPTS{
		Asset:              "DOT",
		Amount:             utils.DecimalFromString("55.56"),
		LastOperationPrice: utils.DecimalFromString("18.45")}
	exp.Timestamp = time.Now().UnixMicro()
	err = insert(exp)
	testutils.AssertNil(t, err, "err")

	exeIds = append(exeIds, exp.AccountId)
	got, err := find_latest_by_exeId(exp.ExeId)

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "laccount")
}
