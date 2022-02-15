package testutils

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

var _OP_COLLECTION_TEST string = "operations"

func RestoreOperationCollection(old func() *mongo.Collection) {
	mongodb.GetOperationsCol = old
}

func MockOperationCollection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetOperationsCol
	mongodb.GetOperationsCol = func() *mongo.Collection {
		return mongoClient.
			Database(_MONGODB_DATABASE_TEST).
			Collection(_OP_COLLECTION_TEST)
	}
	return old
}

func AssertOperations(t *testing.T, expected, gotten model.Operation) {
	if expected.OpId != gotten.OpId {
		t.Fatalf("OpId: expected = %s, gotten = %s", expected.OpId, gotten.OpId)
	}
	if expected.ExeId != gotten.ExeId {
		t.Fatalf("ExeId: expected = %s, gotten = %s", expected.ExeId, gotten.ExeId)
	}
	if expected.Type != gotten.Type {
		t.Fatalf("Type: expected = %v, gotten = %v", expected.Type, gotten.Type)
	}
	if expected.Base != gotten.Base {
		t.Fatalf("Base: expected = %s, gotten = %s", expected.Base, gotten.Base)
	}
	if expected.Quote != gotten.Quote {
		t.Fatalf("Quote: expected = %s, gotten = %s", expected.Quote, gotten.Quote)
	}
	if expected.Side != gotten.Side {
		t.Fatalf("Side: expected = %v, gotten = %v", expected.Side, gotten.Side)
	}
	if expected.Amount != gotten.Amount {
		t.Fatalf("Amount: expected = %f, gotten = %f", expected.Amount, gotten.Amount)
	}
	if expected.AmountSide != gotten.AmountSide {
		t.Fatalf("AmountSide: expected = %v, gotten = %v", expected.AmountSide, gotten.AmountSide)
	}
	if expected.Price != gotten.Price {
		t.Fatalf("Price: expected = %f, gotten = %f", expected.Price, gotten.Price)
	}
	if expected.Status != gotten.Status {
		t.Fatalf("Status: expected = %v, gotten = %v", expected.Status, gotten.Status)
	}
	if expected.Timestamp != gotten.Timestamp {
		t.Fatalf("Timestamp: expected = %d, gotten = %d", expected.Timestamp, gotten.Timestamp)
	}
	AssertOpResult(t, expected.Results, gotten.Results)
}

func AssertOpResult(t *testing.T, expected, gotten model.OpResults) {
	if expected.IsEmpty() && gotten.IsEmpty() {
		return
	}
	if expected.IsEmpty() && !gotten.IsEmpty() {
		t.Error("OpResult: expected empty, gotten initialized")
	}
	if !expected.IsEmpty() && gotten.IsEmpty() {
		t.Error("OpResult: expected initialized, gotten empty")
	}
	if expected.ActualPrice != gotten.ActualPrice {
		t.Fatalf("ActualPrice: expected = %f, gotten = %f", expected.ActualPrice, gotten.ActualPrice)
	}
	if expected.BaseAmount != gotten.BaseAmount {
		t.Fatalf("BaseAmount: expected = %f, gotten = %f", expected.BaseAmount, gotten.BaseAmount)
	}
	if expected.QuoteAmount != gotten.QuoteAmount {
		t.Fatalf("QuoteAmount: expected = %f, gotten = %f", expected.QuoteAmount, gotten.QuoteAmount)
	}
	if expected.Spread != gotten.Spread {
		t.Fatalf("Spread: expected = %f, gotten = %f", expected.Spread, gotten.Spread)
	}
}

func GetOperationTest() model.Operation {
	return model.Operation{
		OpId:       uuid.NewString(),
		ExeId:      uuid.NewString(),
		Type:       model.AUTO,
		Base:       "BTC",
		Quote:      "USDT",
		Side:       model.BUY,
		Amount:     153.78,
		AmountSide: model.BASE_AMOUNT,
		Price:      133.23,
		Results: model.OpResults{
			ActualPrice: 133.58,
			BaseAmount:  153.78,
			QuoteAmount: 11224.56,
			Spread:      12.1,
		},
		Status:    model.FILLED,
		Timestamp: time.Now().UnixMicro()}
}
