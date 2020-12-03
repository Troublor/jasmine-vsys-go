package sdk

import (
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	acc := createAccount("hello", vsys.Testnet)
	retrieved := retrieveAccount(acc.PrivateKey(), vsys.Testnet)
	if acc.Address() != retrieved.Address() {
		t.Fatal()
	}
}
