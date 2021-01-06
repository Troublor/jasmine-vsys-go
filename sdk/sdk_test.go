package sdk

import (
	"context"
	"fmt"
	"github.com/Troublor/jasmine-vsys-go/sdk/transport"
	"math"
	"testing"
	"time"
)

func TestVSYS2MinUnit(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	minUnit := sdk.VSYS2MinUnit(1.1)
	if minUnit != 110000000 {
		t.Fatal("conversion error")
	}
}

func TestMinUnit2VSYS(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	minUnit := sdk.MinUnit2VSYS(110000000)
	if minUnit != 1.1 {
		t.Fatal("conversion error")
	}
}

func TestSDK_RetrieveAccount(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	acc := sdk.RetrieveAccount(testAccountPrivateKey)
	if acc.Address() != testAccountAddress {
		t.Fatal()
	}
}

func TestSDK_DeployTFCSync(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	deployer := sdk.RetrieveAccount(testAccountPrivateKey)
	txId, err := sdk.DeployTFC(deployer)
	if err != nil {
		t.Fatal(err)
	}
	doneCh, errCh := sdk.WaitForConfirmation(context.Background(), txId, 1)
	select {
	case tx := <-doneCh:
		if tx.ContractId == "" {
			t.Fatal()
		}
	case err := <-errCh:
		t.Fatal(err)
	}
}

func TestSDK_WaitForConfirmation_confirmed(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}

	txId := "CcfZwuj6zYSFtsXhYRCprtm874s1ZnU4Udkv6M9Neu2B"
	doneCh, errCh := sdk.WaitForConfirmation(context.Background(), txId, 1)
	select {
	case tx := <-doneCh:
		if tx.ContractId == "" {
			t.Fatal()
		}
	case err := <-errCh:
		t.Fatal(err)
	}
}

func TestSDK_WaitForConfirmation_unconfirmed(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}

	txId := "CcfZwuj6zYSFtsXhYRCprtm874s1ZnU4Udkv6M9Neu2B"
	doneCh, errCh := sdk.WaitForConfirmation(context.Background(), txId, math.MaxInt32)
	later := time.NewTimer(1 * time.Second)
	select {
	case <-later.C:
		return
	case <-doneCh:
		t.Fatal()
	case err := <-errCh:
		t.Fatal(err)
	}
}

func TestSDK_Transfer(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}

	account := sdk.RetrieveAccount(testAccountPrivateKey)
	txId, err := sdk.Transfer(account.Address(), 1, account)
	if err != nil {
		t.Fatal(err)
	}

	doneCh, errCh := sdk.WaitForConfirmation(context.Background(), txId, 0)
	select {
	case <-doneCh:
	case err := <-errCh:
		t.Fatal(err)
	}

	tx, err := sdk.GetTransactionInfo(txId)
	if err != nil {
		t.Fatal(err)
	}
	if tx.Status != "Success" {
		t.Fatal(tx.Status)
	}
}

func TestSDK_BalanceOf(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}

	account := sdk.RetrieveAccount(testAccountPrivateKey)
	balance, err := sdk.BalanceOf(account.Address())
	if balance <= 0 {
		t.Fatal()
	}
}

func TestSDK_GetTransactionInfo(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := sdk.GetTransactionInfo("EkDsc9ipUCXkJk6nWCZNSALmT6F8SiGQJJJcUsgdLSgT")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tx)
}
