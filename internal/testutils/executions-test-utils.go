package testutils

import (
	"reflect"
	"testing"

	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

var _EXE_COLLECTION_TEST string = "executions"

func RestoreExecutionCollection(old func() *mongo.Collection) {
	mongodb.GetExecutionsCol = old
}

func MockExecutionCollection(mongoClient *mongo.Client) func() *mongo.Collection {
	old := mongodb.GetExecutionsCol
	mongodb.GetExecutionsCol = func() *mongo.Collection {
		return mongoClient.
			Database(_MONGODB_DATABASE_TEST).
			Collection(_EXE_COLLECTION_TEST)
	}
	return old
}

func AssertExecutions(t *testing.T, expected, gotten model.Execution) {
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
