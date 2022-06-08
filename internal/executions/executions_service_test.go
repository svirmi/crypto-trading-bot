package executions

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/valerioferretti92/crypto-trading-bot/internal/logger"
	"github.com/valerioferretti92/crypto-trading-bot/internal/model"
	"github.com/valerioferretti92/crypto-trading-bot/internal/mongodb"
	"github.com/valerioferretti92/crypto-trading-bot/internal/testutils"
	"github.com/valerioferretti92/crypto-trading-bot/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateOrRestore_Create(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	balances := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("5.0")},
		{Asset: "ETH", Amount: utils.DecimalFromString("10.45")}}
	raccount := model.RemoteAccount{
		MakerCommission:  0,
		TakerCommission:  1,
		BuyerCommission:  2,
		SellerCommission: 1,
		Balances:         balances}

	got, err := CreateOrRestore(raccount)
	testutils.AssertNil(t, err, "err")

	exp := model.Execution{
		ExeId:     got.ExeId,
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: got.Timestamp}

	exeIds = append(exeIds, got.ExeId)

	testutils.AssertEq(t, exp, got, "execution")
}

func TestCreateOrRestore_Create_EmptyRacc(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	raccount := model.RemoteAccount{
		MakerCommission:  0,
		TakerCommission:  1,
		BuyerCommission:  2,
		SellerCommission: 1,
		Balances:         []model.RemoteBalance{}}

	_, err := CreateOrRestore(raccount)
	testutils.AssertNotNil(t, err, "err")
}
func TestCreateOrRestore_Restore(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteOne(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	balances := []model.RemoteBalance{
		{Asset: "BTC", Amount: utils.DecimalFromString("5.0")},
		{Asset: "ETH", Amount: utils.DecimalFromString("10.45")}}
	raccount := model.RemoteAccount{
		MakerCommission:  0,
		TakerCommission:  1,
		BuyerCommission:  2,
		SellerCommission: 1,
		Balances:         balances}

	exp := model.Execution{
		ExeId:     uuid.NewString(),
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	insert_one(exp)
	exeIds = append(exeIds, exp.ExeId)

	got, err := CreateOrRestore(raccount)

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "execution")
}

func TestGetLatestByExeId(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Inserting execution v1
	exp := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	insert_one(exp)

	// Inserting execution v2
	exp.Status = model.EXE_TERMINATED
	exp.Assets = append(exp.Assets, "DOT")
	exp.Timestamp = time.Now().UnixMicro() + 500
	insert_one(exp)

	// Getting latest by exe id
	got, err := GetLastestByExeId(exeIds[0])

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, exp, got, "execution")
}

func TestGetCurrentlyActive(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString(), uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"$or", bson.A{
			bson.D{{"exeId", exeIds[0]}},
			bson.D{{"exeId", exeIds[1]}}}}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Inserting exe1 v1
	exp := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro() + 100}
	insert_one(exp)

	// Inserting exe2 v2
	exp.Status = model.EXE_TERMINATED
	exp.Assets = append(exp.Assets, "DOT")
	exp.Timestamp = time.Now().UnixMicro() + 200
	insert_one(exp)

	// Inserting exe3 v1
	exp.ExeId = exeIds[1]
	exp.Status = model.EXE_ACTIVE
	exp.Timestamp = time.Now().UnixMicro() + 300
	insert_one(exp)

	// Getting latest by exe id
	got, err := GetCurrentlyActive()

	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, got, exp, "execution")
}

func TestGetCurrentlyActive_None(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()

	// Restoring status after test execution
	defer func() {
		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Getting latest by exe id
	got, err := GetCurrentlyActive()

	testutils.AssertNil(t, err, "err")
	testutils.AssertTrue(t, got.IsEmpty(), "execution")
}

func TestStatuses(t *testing.T) {
	logger.Initialize(false, logrus.TraceLevel)
	// Setting up test
	old := mock_mongo_config()
	mongodb.Initialize()
	var exeIds = []string{uuid.NewString()}

	// Restoring status after test execution
	defer func() {
		filter := bson.D{{"exeId", exeIds[0]}}
		mongodb.GetExecutionsCol().DeleteMany(context.TODO(), filter, nil)

		restore_mongo_config(old)
		mongodb.Disconnect()
	}()

	// Inserting execution v1
	exe := model.Execution{
		ExeId:     exeIds[0],
		Status:    model.EXE_ACTIVE,
		Assets:    []string{"BTC", "ETH"},
		Timestamp: time.Now().UnixMicro()}
	insert_one(exe)

	// Updating status to TERMINATED
	got, err := Terminate(exeIds[0])
	testutils.AssertNil(t, err, "err")
	testutils.AssertEq(t, model.EXE_TERMINATED, got.Status, "execution")
}
