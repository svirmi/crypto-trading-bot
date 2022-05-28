package operations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsert(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetOperationsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	exp := get_operation_test()
	exeIds = append(exeIds, exp.ExeId)
	err := insert(exp)
	testutils.AssertNil(t, err, "err")

	gots, err := find_by_exe_id(exp.ExeId)
	testutils.AssertNil(t, err, "err")

	testutils.AssertEq(t, 1, len(gots), "operations")
	testutils.AssertEq(t, exp, gots[0], "operations")
}

func TestInsertMany(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetOperationsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	exp2 := get_operation_test()
	exp2.ExeId = exeIds[0]
	exp2.Timestamp = time.Now().UnixMicro()
	exp1 := get_operation_test()
	exp1.Timestamp = time.Now().UnixMicro() + 100
	exp1.ExeId = exeIds[0]
	err := insert_many([]model.Operation{exp2, exp1})
	testutils.AssertNil(t, err, "err")

	gots, err := find_by_exe_id(exeIds[0])
	testutils.AssertNil(t, err, "err")

	testutils.AssertEq(t, 2, len(gots), "operations")
	testutils.AssertEq(t, []model.Operation{exp1, exp2}, gots, "operations")
}

func TestFindByExeId(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	gots, err := find_by_exe_id(uuid.NewString())

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, 0, len(gots), "opertions")
}
