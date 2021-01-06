package transport

import (
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"testing"
)

const testAccountPrivate = "Bokp6eDDyiumnxRVVMvWQmRCK6crc1QA3fNBtuN59ubh"
const testAccountAddress = "ATuNyDibv3n2KpYEKHy56c94kGQdRdT7D7E"

func TestNewProvider_normalEndpoint(t *testing.T) {
	p, err := NewProvider(Endpoint[Testnet], Testnet)
	if err != nil {
		t.Fatal(err)
	}
	if p.Endpoint != Endpoint[Testnet] {
		t.Fatal()
	}
}

func TestNewProvider_endpointWithTailingSlash(t *testing.T) {
	p, err := NewProvider(Endpoint[Testnet]+"/", Testnet)
	if err != nil {
		t.Fatal(err)
	}
	if p.Endpoint != Endpoint[Testnet] {
		t.Fatal()
	}
}

func TestNewProvider_invalidEndpoint(t *testing.T) {
	p, err := NewProvider("https://google.com", Testnet)
	if err == nil {
		t.Fatal()
	}
	if p != nil {
		t.Fatal()
	}
}

func TestProvider_Get(t *testing.T) {
	p, err := NewProvider(Endpoint[Testnet], Testnet)
	if err != nil {
		t.Fatal(err)
	}

	var height LatestHeight
	err = p.Get("/blocks/height", nil, &height)
	if err != nil {
		t.Fatal(err)
	}

	if height.Height <= 0 {
		t.Fatal(height)
	}
}

func TestProvider_GetWithQuery(t *testing.T) {
	p, err := NewProvider(Endpoint[Testnet], Testnet)
	if err != nil {
		t.Fatal(err)
	}

	var txs TransactionList
	err = p.Get("/transactions/list", map[string]string{
		"address": testAccountAddress,
		"txType":  "8",
		"limit":   "1",
		"offset":  "0",
	}, &txs)
	if err != nil {
		t.Fatal(err)
	}

	if txs.Size != 1 || len(txs.Transactions) != txs.Size {
		t.Fatal()
	}
}

func TestProvider_PostWithStringData(t *testing.T) {
	p, err := NewProvider(Endpoint[Testnet], Testnet)
	if err != nil {
		t.Fatal(err)
	}

	var hash SecureCryptographicHash
	err = p.Post("/utils/hash/secure", nil, "abc", &hash)
	if err != nil {
		t.Fatal(err)
	}

	if hash.Hash != "B9rbM2PGvFEqHw6T9rBeDi9Yrni6f1YC55zYoJY4QEFv" {
		t.Fatal()
	}
}

func TestProvider_PostWithJSONData(t *testing.T) {
	p, err := NewProvider(Endpoint[Testnet], Testnet)
	if err != nil {
		t.Fatal(err)
	}

	acc := vsys.InitAccount(p.NetType)
	acc.BuildFromPrivateKey(testAccountPrivate)
	tx := acc.BuildPayment(acc.Address(), 1, acc.Address())
	var resp map[string]interface{}
	err = p.Post("/vsys/broadcast/payment", nil, tx, &resp)
	if err != nil {
		t.Fatal(err)
	}
}
