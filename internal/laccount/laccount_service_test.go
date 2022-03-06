package laccount

import (
	"context"
	"testing"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateOrRestore_Create_FTS(t *testing.T) {
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

	local_account_init := get_laccount_init_test(model.FIXED_THRESHOLD_STRATEGY)
	exeIds = append(exeIds, local_account_init.ExeId)

	gotten, err := CreateOrRestore(local_account_init)
	if err != nil {
		t.Fatalf("expected err == nil, gotten = %v", err)
	}
	if gotten == nil {
		t.Error("expected laccount != nil, gotten = nil")
	}

	expected := get_laccount_test_FTS()
	expected.ExeId = local_account_init.ExeId
	expected.AccountId = gotten.GetAccountId()
	expected.Timestamp = gotten.GetTimestamp()

	testutils.AssertStructEq(t, expected, gotten)
}

func TestCreateOrRestore_Restore_FTS(t *testing.T) {
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

	expected := get_laccount_test_FTS()
	exeIds = append(exeIds, expected.ExeId)
	err := insert(expected)
	if err != nil {
		t.Fatalf("expected err = nil, gotten err = %v", err)
	}

	local_account_init := get_laccount_init_test(model.FIXED_THRESHOLD_STRATEGY)
	local_account_init.ExeId = exeIds[0]
	gotten, err := CreateOrRestore(local_account_init)
	if err != nil {
		t.Fatalf("expected err == nil, gotten = %v", err)
	}
	if gotten == nil {
		t.Error("expected laccount != nil, gotten = nil")
	}

	testutils.AssertStructEq(t, expected, gotten)
}
