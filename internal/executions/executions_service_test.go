package executions

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateOrRestore_Create(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	got, err := CreateOrRestore(get_execution_init())
	testutils.AssertNil(t, err, "err")
	exeIds = append(exeIds, got.ExeId)

	exp := get_execution()
	exp.ExeId = got.ExeId
	exp.Timestamp = got.Timestamp

	testutils.AssertEq(t, exp, got, "execution")
}

func TestCreateOrRestore_Create_EmptyRacc(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	exeReq := get_execution_init()
	exeReq.Raccount.Balances = []model.RemoteBalance{}

	_, err := CreateOrRestore(exeReq)
	testutils.AssertNotNil(t, err, "err")
}

func TestCreateOrRestore_Restore(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	exp := get_execution()
	insert_one(exp)
	exeIds = append(exeIds, exp.ExeId)

	got, err := CreateOrRestore(get_execution_init())

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "execution")
}

func TestGetLatestByExeId(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Inserting execution v1
	exp := get_execution()
	exp.ExeId = exeIds[0]
	insert_one(exp)

	// Inserting execution v2
	exp.Status = model.EXE_TERMINATED
	exp.Assets = append(exp.Assets, "DOT")
	exp.Timestamp = time.Now().UnixMicro() + 500
	insert_one(exp)

	// Getting latest by exe id
	got, err := GetLastestByExeId(exeIds[0])

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "execution")
}

func TestGetCurrentlyActive(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString(), uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"$or", bson.A{
			bson.D{{"exeId", exeIds[0]}},
			bson.D{{"exeId", exeIds[1]}}}}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Inserting exe1 v1
	exp := get_execution()
	exp.ExeId = exeIds[0]
	insert_one(exp)

	// Inserting exe1 v2
	exp.Status = model.EXE_TERMINATED
	exp.Assets = append(exp.Assets, "DOT")
	exp.Timestamp = time.Now().UnixMicro() + 200
	insert_one(exp)

	// Inserting exe3 v1
	exp.ExeId = exeIds[1]
	exp.Status = model.EXE_ACTIVE
	exp.Timestamp = time.Now().UnixMicro() + 300
	insert_one(exp)

	// Getting latest by exe id
	got, err := GetCurrentlyActive()

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, got, exp, "execution")
}

func TestGetCurrentlyActive_None(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Getting latest by exe id
	got, err := GetCurrentlyActive()

	testutils.AssertNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "execution")
}

func TestStatuses(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Inserting execution v1
	exe := get_execution()
	exe.ExeId = exeIds[0]
	insert_one(exe)

	// Updating status to TERMINATED
	got, err := Terminate(exeIds[0])
	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, model.EXE_TERMINATED, got.Status, "execution")
}
