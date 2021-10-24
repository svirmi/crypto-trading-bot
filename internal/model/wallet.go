package model

import "reflect"

type Account struct {
	MakerCommission  int64
	TakerCommission  int64
	BuyerCommission  int64
	SellerCommission int64
	Balances         []Balance
}

func (a Account) IsEmpty() bool {
	return reflect.DeepEqual(a, Account{})
}

type Balance struct {
	Asset  string
	Amount string
}

func (b Balance) IsEmpty() bool {
	return reflect.DeepEqual(b, Balance{})
}
