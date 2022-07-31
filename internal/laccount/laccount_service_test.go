package laccount

import (
	"context"
	"testing"

	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"go.mongodb.org/mongo-driver/bson"
)

/**************************** DTS ******************************/

func TestCreate_DTS(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old_mongo_conf := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old_mongo_conf)
		mongodb.Disconnect()
	}()

	local_account_init := get_laccount_init_test(model.DTS_STRATEGY)
	exeIds = append(exeIds, local_account_init.ExeId)

	got, err := Create(local_account_init)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_test_DTS()
	exp.ExeId = local_account_init.ExeId
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()

	testutils.AssertEq(t, exp, got, "laccount")
}

func TestCreate_DTS_EmptyRAcc(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old_mongo_conf := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old_mongo_conf)
		mongodb.Disconnect()
	}()

	local_account_init := get_laccount_init_test(model.DTS_STRATEGY)
	local_account_init.RAccount.Balances = make([]model.RemoteBalance, 0)

	_, err := Create(local_account_init)
	testutils.AssertNotNil(t, err, "err")
}

func TestCreate_DTS_AlreadyExists(t *testing.T) {
	logger.Initialize(false, true, true)
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

	exp := get_laccount_test_DTS()
	exeIds = append(exeIds, exp.ExeId)
	err := insert(exp)
	testutils.AssertNil(t, err, "err")

	local_account_init := get_laccount_init_test(model.DTS_STRATEGY)
	local_account_init.ExeId = exeIds[0]
	got, err := Create(local_account_init)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertNil(t, got, "laccount")
}

/**************************** PTS ******************************/

func TestCreate_PTS(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old_mongo_conf := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"metadata.exeId", exeIds[0]}}
		mongodb.GetLocalAccountsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old_mongo_conf)
		mongodb.Disconnect()
	}()

	local_account_init := get_laccount_init_test(model.PTS_STRATEGY)
	exeIds = append(exeIds, local_account_init.ExeId)

	got, err := Create(local_account_init)
	testutils.AssertNil(t, err, "err")

	exp := get_laccount_test_PTS()
	exp.ExeId = local_account_init.ExeId
	exp.AccountId = got.GetAccountId()
	exp.Timestamp = got.GetTimestamp()

	testutils.AssertEq(t, exp, got, "laccount")
}

func TestCreate_PTS_EmptyRAcc(t *testing.T) {
	logger.Initialize(false, true, true)
	// Setting up test
	old_mongo_conf := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old_mongo_conf)
		mongodb.Disconnect()
	}()

	local_account_init := get_laccount_init_test(model.PTS_STRATEGY)
	local_account_init.RAccount.Balances = make([]model.RemoteBalance, 0)

	_, err := Create(local_account_init)
	testutils.AssertNotNil(t, err, "err")
}

func TestCreate_PTS_AlreadyExists(t *testing.T) {
	logger.Initialize(false, true, true)
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

	exp := get_laccount_test_PTS()
	exeIds = append(exeIds, exp.ExeId)
	err := insert(exp)
	testutils.AssertNil(t, err, "err")

	local_account_init := get_laccount_init_test(model.PTS_STRATEGY)
	local_account_init.ExeId = exeIds[0]
	got, err := Create(local_account_init)

	testutils.AssertNotNil(t, err, "err")
	testutils.AssertNil(t, got, "laccount")
}
