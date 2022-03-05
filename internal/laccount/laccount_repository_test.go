package laccount

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/strategy/fts"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsert_FTS(t *testing.T) {
	// Setting up test
	mongoClient := mongodb.GetMongoClientTest()
	old := mock_laccount_collection(mongoClient)
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteOne(context.TODO(), filter, nil)

		restore_laccount_collection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	laccount := get_laccount_test_FTS()
	exeIds = append(exeIds, laccount.ExeId)
	err := insert(laccount)
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}

	gotten, err := find_latest_by_exeId(laccount.ExeId)
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}
	testutils.AssertStructEq(t, laccount, gotten.(fts.LocalAccountFTS))
}

func TestFindLatestByExeId_FTS(t *testing.T) {
	// Setting up test
	mongoClient := mongodb.GetMongoClientTest()
	old := mock_laccount_collection(mongoClient)
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteMany(context.TODO(), filter, nil)

		restore_laccount_collection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	laccount := get_laccount_test_FTS()
	exeIds = append(exeIds, laccount.ExeId)
	err := insert(laccount)
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}

	laccount.Assets["DOT"] = fts.AssetStatusFTS{
		Asset:              "DOT",
		Amount:             decimal.NewFromFloat32(55.56),
		Usdt:               decimal.Zero,
		LastOperationType:  fts.OP_BUY_FTS,
		LastOperationPrice: decimal.NewFromFloat32(18.45)}
	laccount.Timestamp = time.Now().UnixMicro()
	err = insert(laccount)
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}

	exeIds = append(exeIds, laccount.AccountId)
	gotten, err := find_latest_by_exeId(laccount.ExeId)
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}

	testutils.AssertStructEq(t, laccount, gotten.(fts.LocalAccountFTS))
}

func TestFindLatestByExeId_FTS_None(t *testing.T) {
	// Setting up test
	mongoClient := mongodb.GetMongoClientTest()
	old := mock_laccount_collection(mongoClient)

	// Restoring status after test execution
	defer func() {
		restore_laccount_collection(old)
		mongoClient.Disconnect(context.TODO())
	}()

	gotten, err := find_latest_by_exeId(uuid.NewString())
	if err != nil {
		t.Fatalf("expected err = nil, gotten = %v", err)
	}
	if gotten != nil {
		t.Fatalf("expected laccount = nil, gotten = %v", gotten)
	}
}
