package executions

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsertOne(t *testing.T) {
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
	exp := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	insert_one(exp)

	// Getting execution object from DB
	var got model.Execution
	filter := bson.D{{"exeId", exp.ExeId}}
	mongodb.GetExecutionsCol().FindOne(context.TODO(), filter).Decode(&got)

	testutils.AssertEq(t, exp, got, "execution")
}

func TestFindLatestByExeId(t *testing.T) {
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
	exe1 := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro() + 100}
	docs = append(docs, exe1)
	exe2 := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_PAUSED,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro() + 200}
	docs = append(docs, exe2)
	exp := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro() + 300}
	docs = append(docs, exp)
	other1 := model.Execution{
		ExeId:     otherExeIds[0],
		Status:    model.EXE_TERMINATED,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	docs = append(docs, other1)
	mongodb.GetExecutionsCol().InsertMany(context.TODO(), docs, nil)

	// Getting execution object from DB
	got, err := find_latest_by_exeId(exeIds[0])

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "execution")
}

func TestFindLatestByExeId_NoResults(t *testing.T) {
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
	testutils.AssertNotNil(t, err, "err")
}

func TestFindCurrentlyActive(t *testing.T) {
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString()}
	var otherExeIds = []string{uuid.NewString(), uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter1 := bson.M{"exeId": exeIds[0]}
		filter3 := bson.M{"exeId": otherExeIds[0]}
		filter4 := bson.M{"exeId": otherExeIds[1]}
		filter := bson.M{"$or": []bson.M{filter1, filter3, filter4}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Building test execution
	var docs []interface{}
	exe1 := model.Execution{
		ExeId:     otherExeIds[1],
		Status:    model.EXE_TERMINATED,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro() + 100}
	docs = append(docs, exe1)
	exe2 := model.Execution{
		ExeId:     otherExeIds[0],
		Status:    model.EXE_TERMINATED,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro() + 200}
	docs = append(docs, exe2)
	exp := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro() + 300}
	docs = append(docs, exp)
	mongodb.GetExecutionsCol().InsertMany(context.TODO(), docs, nil)

	// Getting execution object from DB
	got, err := find_currently_active()

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "execution")
}

func TestFindCurrentlyActive_NoResults(t *testing.T) {
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
