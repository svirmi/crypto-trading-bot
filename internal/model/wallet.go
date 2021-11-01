package model

import "reflect"

// Representation of remote Binance wallet
type RemoteAccount struct {
	MakerCommission  int64
	TakerCommission  int64
	BuyerCommission  int64
	SellerCommission int64
	Balances         []RemoteBalance
}

func (a RemoteAccount) IsEmpty() bool {
	return reflect.DeepEqual(a, RemoteAccount{})
}

type RemoteBalance struct {
	Asset  string
	Amount string
}

func (b RemoteBalance) IsEmpty() bool {
	return reflect.DeepEqual(b, RemoteBalance{})
}

// Representation of local wallet
type LocalAccount struct {
	AccountId string                  `bson:"accountId"` // Local account object id
	ExeId     string                  `bson:"exeId"`     // Execution id this local wallet is bound to
	Balances  map[string]LocalBalance `bson:"balances"`  // Map local balances
	Timestamp int64                   `bson:"timestamp"` // Timestamp
}

func (a LocalAccount) IsEmpty() bool {
	return reflect.DeepEqual(a, LocalAccount{})
}

type LocalBalance struct {
	Asset        string   `bson:"asset"`         // Asset being tracked
	Amount       float32  `bson:"amount"`        // Amount of "asset" currently owned
	Usdt         float32  `bson:"usdt"`          // Usdt gotten by trading "asset"
	OperationIds []string `bson:"operationsIds"` // Operations where "asset" was traded back and forth for USDT
}

func (b LocalBalance) IsEmpty() bool {
	return reflect.DeepEqual(b, LocalBalance{})
}
