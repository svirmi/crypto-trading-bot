package operations

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

func TestInsert(t *testing.T) {
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

	expected := get_operation_test()
	exeIds = append(exeIds, expected.ExeId)
	err := insert(expected)
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}

	gottens, err := find_by_exe_id(expected.ExeId)
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}

	if len(gottens) != 1 {
		t.Fatalf("expected len(results) = 1, gotten len(results) = %d", len(gottens))
		return
	}
	testutils.AssertStructEq(t, expected, gottens[0])
}

func TestInsertMany(t *testing.T) {
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

	expected1 := get_operation_test()
	expected1.ExeId = exeIds[0]
	expected1.Timestamp = time.Now().UnixMicro()
	expected2 := get_operation_test()
	expected2.Timestamp = time.Now().UnixMicro() + 100
	expected2.ExeId = exeIds[0]
	err := insert_many([]model.Operation{expected1, expected2})
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}

	gottens, err := find_by_exe_id(exeIds[0])
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}

	if len(gottens) != 2 {
		t.Fatalf("expected len(results) = 2, gotten len(results) = %d", len(gottens))
		return
	}
	testutils.AssertStructEq(t, expected2, gottens[0])
	testutils.AssertStructEq(t, expected1, gottens[1])
}

func TestFindByExeId(t *testing.T) {
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	gottens, err := find_by_exe_id(uuid.NewString())
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}

	if len(gottens) != 0 {
		t.Fatalf("expected len(results) = , gotten len(results) = %d", len(gottens))
	}
}
