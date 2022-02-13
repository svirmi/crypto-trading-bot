package laccount

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsert_FTS(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockLaccountCollection(mongoClient)
	var laccIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.accountId", laccIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteOne(context.TODO(), filter, nil)

		testutils.RestoreLaccountCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	laccount := testutils.GetLocalAccountTFSTest()
	laccIds = append(laccIds, laccount.AccountId)
	err := insert(laccount)
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}

	gotten, err := find_latest_by_exeId(laccount.ExeId)
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}
	testutils.AssertLocalAccountFTS(t, laccount, gotten.(fts.LocalAccountFTS))
}

func TestFindLatestByExeId_FTS(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockLaccountCollection(mongoClient)
	var laccIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.accountId", laccIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteMany(context.TODO(), filter, nil)

		testutils.RestoreLaccountCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	laccount := testutils.GetLocalAccountTFSTest()
	err := insert(laccount)
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}

	laccount.Assets["DOT"] = fts.AssetStatusFTS{
		Asset:              "DOT",
		Amount:             55.56,
		Usdt:               0,
		LastOperationType:  fts.OP_BUY_FTS,
		LastOperationPrice: 18.45}
	laccount.Timestamp = time.Now().UnixMicro()
	err = insert(laccount)
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}

	laccIds = append(laccIds, laccount.AccountId)
	gotten, err := find_latest_by_exeId(laccount.ExeId)
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}

	testutils.AssertLocalAccountFTS(t, laccount, gotten.(fts.LocalAccountFTS))
}

func TestFindLatestByExeId_FTS_None(t *testing.T) {
	// Setting up test
	mongoClient := testutils.GetMongoClientTest()
	old := testutils.MockLaccountCollection(mongoClient)

	// Restoring status after test execution
	defer func() {
		testutils.RestoreLaccountCollection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	gotten, err := find_latest_by_exeId(uuid.NewString())
	if err != nil {
		t.Errorf("expected err = nil, gotten = %v", err)
	}
	if gotten != nil {
		t.Errorf("expected laccount = nil, gotten = %v", gotten)
	}
}