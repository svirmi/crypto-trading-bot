package executions

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsertOne(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Building test execution
	exp := get_execution()
	exp.ExeId = exeIds[0]
	insert_one(exp)

	// Getting execution object from DB
	var got model.Execution
	filter := bson.D{{"exeId", exp.ExeId}}
	mongodb.GetExecutionsCol().FindOne(context.TODO(), filter).Decode(&got)

	testutils.AssertEq(t, exp, got, "execution")
}

func TestFindLatestByExeId(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString()}
	var otherExeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter1 := bson.M{"exeId": exeIds[0]}
		filter4 := bson.M{"exeId": otherExeIds[0]}
		filter := bson.M{"$or": []bson.M{filter1, filter4}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Building test execution
	var docs []interface{}
	exe1 := get_execution()
	exe1.ExeId = exeIds[0]
	docs = append(docs, exe1)

	exp := exe1
	exp.Timestamp = exe1.Timestamp + 300
	docs = append(docs, exp)

	other1 := get_execution()
	other1.ExeId = otherExeIds[0]
	other1.Status = model.EXE_TERMINATED
	docs = append(docs, other1)
	mongodb.GetExecutionsCol().InsertMany(context.TODO(), docs, nil)

	// Getting execution object from DB
	got, err := find_latest_by_exeId(exeIds[0])

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "execution")
}

func TestFindLatestByExeId_NoResults(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Getting execution object from DB
	got, err := find_latest_by_exeId(uuid.NewString())

	testutils.AssertTrue(t, got.IsEmpty(), "execution")
	testutils.AssertNil(t, err, "err")
}

func TestFindCurrentlyActive(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString()}
	var otherExeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter1 := bson.M{"exeId": exeIds[0]}
		filter2 := bson.M{"exeId": otherExeIds[0]}
		filter := bson.M{"$or": []bson.M{filter1, filter2}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Building test execution
	var docs []interface{}
	exe1 := get_execution()
	exe1.ExeId = otherExeIds[0]
	exe1.Status = model.EXE_ACTIVE
	docs = append(docs, exe1)

	exe2 := exe1
	exe2.Timestamp = exe2.Timestamp + 100
	exe2.Status = model.EXE_TERMINATED
	docs = append(docs, exe2)

	exp := get_execution()
	exp.ExeId = exeIds[0]
	exp.Timestamp = exe2.Timestamp + 100
	docs = append(docs, exp)

	mongodb.GetExecutionsCol().InsertMany(context.TODO(), docs, nil)

	// Getting execution object from DB
	got, err := find_currently_active()

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "execution")
}

func TestFindCurrentlyActive_NoResults(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Getting execution object from DB
	got, err := find_currently_active()

	testutils.AssertTrue(t, got.IsEmpty(), "execution")
	testutils.AssertNil(t, err, "err")
}
