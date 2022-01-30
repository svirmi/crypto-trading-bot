package executions

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utilstest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestInsertOne(t *testing.T) {
	// Setting up test
	mongoClient := utilstest.GetMongoClientTest()
	old := mock_execution_collection(mongoClient)
	var exeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteOne(context.TODO(), filter, nil)

		restore_execution_collection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Building test execution
	expected := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMilli()}
	insert_one(expected)

	// Getting execution object from DB
	var gotten model.Execution
	filter := bson.D{{"exeId", expected.ExeId}}
	mongodb.GetExecutionsCol().FindOne(context.TODO(), filter).Decode(&gotten)

	// Assertions
	assert_executions(t, expected, gotten)
}

func TestFindLatestByExeId(t *testing.T) {
	// Setting up test
	mongoClient := utilstest.GetMongoClientTest()
	old := mock_execution_collection(mongoClient)
	var exeIds = []string{uuid.NewString()}
	var otherExeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter1 := bson.M{"exeId": exeIds[0]}
		filter4 := bson.M{"exeId": otherExeIds[0]}
		filter := bson.M{"$or": []bson.M{filter1, filter4}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_execution_collection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Building test execution
	var docs []interface{}
	exe1 := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMilli() + 100}
	docs = append(docs, exe1)
	exe2 := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_PAUSED,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMilli() + 200}
	docs = append(docs, exe2)
	expected := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMilli() + 300}
	docs = append(docs, expected)
	other1 := model.Execution{
		ExeId:     otherExeIds[0],
		Status:    model.EXE_TERMINATED,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMilli()}
	docs = append(docs, other1)
	mongodb.GetExecutionsCol().InsertMany(context.TODO(), docs, nil)

	// Getting execution object from DB
	gotten, err := find_latest_by_exeId(exeIds[0])
	if err != nil {
		t.Errorf(err.Error())
	}

	// Assertions
	assert_executions(t, expected, gotten)
}

func TestFindLatestByExeId_NoResults(t *testing.T) {
	// Setting up test
	mongoClient := utilstest.GetMongoClientTest()
	old := mock_execution_collection(mongoClient)

	// Restoring status after test execution
	defer func() {
		restore_execution_collection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Getting execution object from DB
	gotten, err := find_latest_by_exeId(uuid.NewString())
	if !reflect.DeepEqual(gotten, model.Execution{}) {
		t.Errorf("execution object: expected {}, gotten %v", gotten)
	}
	if err == nil {
		t.Errorf("execution repository error: expected != nil, gotten nil")
	}
}

func TestFindCurrentlyActive(t *testing.T) {
	// Setting up test
	mongoClient := utilstest.GetMongoClientTest()
	old := mock_execution_collection(mongoClient)
	var exeIds = []string{uuid.NewString()}
	var otherExeIds = []string{uuid.NewString(), uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter1 := bson.M{"exeId": exeIds[0]}
		filter3 := bson.M{"exeId": otherExeIds[0]}
		filter4 := bson.M{"exeId": otherExeIds[1]}
		filter := bson.M{"$or": []bson.M{filter1, filter3, filter4}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_execution_collection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Building test execution
	var docs []interface{}
	expected := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMilli() + 100}
	docs = append(docs, expected)
	exe1 := model.Execution{
		ExeId:     otherExeIds[0],
		Status:    model.EXE_TERMINATED,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMilli() + 200}
	docs = append(docs, exe1)
	exe2 := model.Execution{
		ExeId:     otherExeIds[1],
		Status:    model.EXE_TERMINATED,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMilli() + 300}
	docs = append(docs, exe2)
	mongodb.GetExecutionsCol().InsertMany(context.TODO(), docs, nil)

	// Getting execution object from DB
	gotten, err := find_currently_active()
	if err != nil {
		t.Errorf(err.Error())
	}

	// Assertions
	assert_executions(t, expected, gotten)
}

func TestFindCurrentlyActive_NoResults(t *testing.T) {
	// Setting up test
	mongoClient := utilstest.GetMongoClientTest()
	old := mock_execution_collection(mongoClient)

	// Restoring status after test execution
	defer func() {
		restore_execution_collection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Getting execution object from DB
	gotten, err := find_currently_active()
	if !reflect.DeepEqual(gotten, model.Execution{}) {
		t.Errorf("execution object: expected {}, gotten %v", gotten)
	}
	if err != nil {
		t.Errorf("execution repository error: expected == nil, gotten %v", err)
	}
}

func restore_execution_collection(old func() *mongo.Collection) {
	mongodb.GetExecutionsCol = old
}

func mock_execution_collection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetExecutionsCol
	mongodb.GetExecutionsCol = func() *mongo.Collection {
		return mongoClient.
			Database(utilstest.MONGODB_DATABASE_TEST).
			Collection(utilstest.EXE_COLLECTION_TEST)
	}
	return old
}

func assert_executions(t *testing.T, expected, gotten model.Execution) {
	if expected.ExeId != gotten.ExeId {
		t.Errorf("ExeId: expected %s, gotten %s", expected.ExeId, gotten.ExeId)
	}
	if expected.Status != gotten.Status {
		t.Errorf("Status: expected %s, gotten %s", expected.Status, gotten.Status)
	}
	if !reflect.DeepEqual(expected.Assets, gotten.Assets) {
		t.Errorf("Assets: expected %v, gotten %v", expected.Assets, gotten.Assets)
	}
	if expected.Timestamp != gotten.Timestamp {
		t.Errorf("Timestamp: expected %v, gotten %v", expected.Timestamp, gotten.Timestamp)
	}
}
