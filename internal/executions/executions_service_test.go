package executions

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateOrRestore_Create(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockExecutionCollection(mongoClient)
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteOne(context.TODO(), filter, nil)

		testutils.RestoreExecutionCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	balances := []model.RemoteBalance{
		{Asset: "BTC", Amount: 5.0},
		{Asset: "ETH", Amount: 10.45}}
	raccount := model.RemoteAccount{
		MakerCommission:  0,
		TakerCommission:  1,
		BuyerCommission:  2,
		SellerCommission: 1,
		Balances:         balances}

	gotten, err := CreateOrRestore(raccount)
	exeIds = append(exeIds, gotten.ExeId)

	if err != nil {
		t.Errorf("exepected nil, gotten %v", err)
	}
	if gotten.ExeId == "" {
		t.Errorf("expected exeId != \"\", gotten \"\"")
	}
	if !reflect.DeepEqual([]string{"BTC", "ETH"}, gotten.Assets) {
		t.Errorf("expected assets = [BTC, ETH], gotten = %v", gotten.Assets)
	}
	if gotten.Status != model.EXE_ACTIVE {
		t.Errorf("expected status = ACTIVE, gotten %v", gotten.Status)
	}
	if gotten.Timestamp == 0 {
		t.Errorf("expected timestamp != 0, gotten 0")
	}
}

func TestCreateOrRestore_Restore(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockExecutionCollection(mongoClient)
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteOne(context.TODO(), filter, nil)

		testutils.RestoreExecutionCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	balances := []model.RemoteBalance{
		{Asset: "BTC", Amount: 5.0},
		{Asset: "ETH", Amount: 10.45}}
	raccount := model.RemoteAccount{
		MakerCommission:  0,
		TakerCommission:  1,
		BuyerCommission:  2,
		SellerCommission: 1,
		Balances:         balances}

	exe := model.Execution{
		ExeId:     uuid.NewString(),
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	insert_one(exe)
	exeIds = append(exeIds, exe.ExeId)

	gotten, err := CreateOrRestore(raccount)
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}

	testutils.AssertExecutions(t, gotten, exe)
}

func TestGetLatestByExeId(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockExecutionCollection(mongoClient)
	var exeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		testutils.RestoreExecutionCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Inserting execution v1
	exe := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	insert_one(exe)

	// Inserting execution v2
	exe.Status = model.EXE_TERMINATED
	exe.Assets = append(exe.Assets, "DOT")
	exe.Timestamp = time.Now().UnixMicro() + 500
	insert_one(exe)

	// Getting latest by exe id
	gotten, err := GetLatestByExeId(exeIds[0])
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}
	testutils.AssertExecutions(t, exe, gotten)
}

func TestGetCurrentlyActive(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockExecutionCollection(mongoClient)
	var exeIds = []string{uuid.NewString(), uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"$or", bson.A{
			bson.D{{"exeId", exeIds[0]}},
			bson.D{{"exeId", exeIds[1]}}}}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		testutils.RestoreExecutionCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Inserting exe1 v1
	exe := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	insert_one(exe)

	// Inserting exe2 v2
	exe.Status = model.EXE_TERMINATED
	exe.Assets = append(exe.Assets, "DOT")
	exe.Timestamp = time.Now().UnixMicro() + 500
	insert_one(exe)

	// Inserting exe3 v1
	exe.ExeId = exeIds[1]
	exe.Status = model.EXE_PAUSED
	insert_one(exe)

	// Getting latest by exe id
	gotten, err := GetCurrentlyActive()
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}
	testutils.AssertExecutions(t, exe, gotten)
}

func TestGetCurrentlyActive_None(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockExecutionCollection(mongoClient)

	// Restoring status after test execution
	defer func() {
		testutils.RestoreExecutionCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Getting latest by exe id
	gotten, err := GetCurrentlyActive()
	if err != nil {
		t.Errorf("expected err = nil, gptten = %v", err)
	}
	if !gotten.IsEmpty() {
		t.Errorf("expected exe = model.Execution{}, gotten = %v", gotten)
	}
}

func TestGetCurrentlyActive_Many(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockExecutionCollection(mongoClient)
	var exeIds = []string{uuid.NewString(), uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"$or", bson.A{
			bson.D{{"exeId", exeIds[0]}},
			bson.D{{"exeId", exeIds[1]}}}}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		testutils.RestoreExecutionCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Inserting exe1 v1
	exe := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	insert_one(exe)

	// Inserting exe2 v2
	exe.Status = model.EXE_ACTIVE
	exe.Assets = append(exe.Assets, "DOT")
	exe.Timestamp = time.Now().UnixMicro() + 500
	insert_one(exe)

	// Inserting exe3 v1
	exe.ExeId = exeIds[1]
	exe.Status = model.EXE_PAUSED
	insert_one(exe)

	// Getting latest by exe id
	gotten, err := GetCurrentlyActive()
	if err == nil {
		t.Errorf("expected err != nil, gotten = nil")
	}
	if !gotten.IsEmpty() {
		t.Errorf("expected exe = model.Execution{}, gotten = %v", gotten)
	}
}

func TestStatuses(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockExecutionCollection(mongoClient)
	var exeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		testutils.RestoreExecutionCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	// Inserting execution v1
	exe := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	insert_one(exe)

	// Updating status to PAUSED
	gotten, err := Pause(exeIds[0])
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}
	if gotten.Status != model.EXE_PAUSED {
		t.Errorf("expected exe.Status = PAUSED, gotten = %v", gotten.Status)
	}

	// Updating status to ACTIVE
	gotten, err = Resume(exeIds[0])
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}
	if gotten.Status != model.EXE_ACTIVE {
		t.Errorf("expected exe.Status = ACTIVE, gotten = %v", gotten.Status)
	}

	// Updating status to ACTIVE again
	gotten, err = Resume(exeIds[0])
	if err == nil {
		t.Error("expected err != nil, gotten = nil", err)
	}
	if !gotten.IsEmpty() {
		t.Errorf("expected exe = model.Execution{}, gotten = %v", gotten)
	}

	// Updating status to TERMINATED
	gotten, err = Terminate(exeIds[0])
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}
	if gotten.Status != model.EXE_TERMINATED {
		t.Errorf("expected exe.Status = TERMINATED, gotten = %v", gotten.Status)
	}
}
