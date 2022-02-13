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
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockLaccountCollection(mongoClient)
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteOne(context.TODO(), filter, nil)

		testutils.RestoreLaccountCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	local_account_init := testutils.GetLocalAccountInitTest(model.FIXED_THRESHOLD_STRATEGY)

	exeIds = append(exeIds, local_account_init.ExeId)
	gotten, err := CreateOrRestore(local_account_init)
	if err != nil {
		t.Errorf("expected err == nil, gotten = %v", err)
	}
	if gotten == nil {
		t.Error("expected laccount != nil, gotten = nil")
	}

	testutils.AssertInitLocalAccount(t, local_account_init, gotten)
}

func TestCreateOrRestore_Restore_FTS(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockLaccountCollection(mongoClient)
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteOne(context.TODO(), filter, nil)

		testutils.RestoreLaccountCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	laccount := testutils.GetLocalAccountTest_FTS()
	exeIds = append(exeIds, laccount.ExeId)
	err := insert(laccount)
	if err != nil {
		t.Errorf("expected err = nil, gotten err = %v", err)
	}

	local_account_init := testutils.GetLocalAccountInitTest(model.FIXED_THRESHOLD_STRATEGY)
	local_account_init.ExeId = exeIds[0]
	gotten, err := CreateOrRestore(local_account_init)
	if err != nil {
		t.Errorf("expected err == nil, gotten = %v", err)
	}
	if gotten == nil {
		t.Error("expected laccount != nil, gotten = nil")
	}

	testutils.AssertInitLocalAccount(t, local_account_init, gotten)
}
