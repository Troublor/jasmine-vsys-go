package sdk

import "github.com/virtualeconomy/go-v-sdk/vsys"

type Address = string

type Account struct {
	privateKey string
	address    Address
	vsysAcc    *vsys.Account
}

func createAccount(seed string, netType vsys.NetType) *Account {
	acc := vsys.InitAccount(netType)
	acc.BuildFromSeed(seed, 0)
	return &Account{
		privateKey: acc.PrivateKey(),
		address:    acc.Address(),
		vsysAcc:    acc,
	}
}

func retrieveAccount(privateKey string, netType vsys.NetType) *Account {
	acc := vsys.InitAccount(netType)
	acc.BuildFromPrivateKey(privateKey)
	return &Account{
		privateKey: acc.PrivateKey(),
		address:    acc.Address(),
		vsysAcc:    acc,
	}
}

func (acc *Account) PrivateKey() string {
	return acc.privateKey
}

func (acc *Account) Address() Address {
	return acc.address
}
