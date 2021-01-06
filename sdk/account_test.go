package sdk

import (
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"testing"
)

func TestCreateRetrieveAccount(t *testing.T) {
	acc := createAccount("hello", vsys.Testnet)
	retrieved := retrieveAccount(acc.PrivateKey(), vsys.Testnet)
	if acc.Address() != retrieved.Address() {
		t.Fatal()
	}
}

func TestRetrieveAccount(t *testing.T) {
	acc := retrieveAccount(testAccountPrivateKey, vsys.Testnet)
	if acc.Address() != testAccountAddress {
		t.Fatal()
	}
}
