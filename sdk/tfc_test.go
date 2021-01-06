package sdk

import (
	"context"
	sdkErr "github.com/Troublor/jasmine-vsys-go/sdk/error"
	"github.com/Troublor/jasmine-vsys-go/sdk/transport"
	"testing"
)

const tokenId = "TWsacZbyKuXwxVBSadvEG9Wi6BW98QkVqrNU5yWJN"

func TestTFC_Mint(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	admin := sdk.RetrieveAccount(testAccountPrivateKey)
	tfc := sdk.TFCWithTokenId(tokenId)

	balanceOld, err := tfc.BalanceOf(admin.Address())
	if err != nil {
		t.Fatal(err)
	}

	txId, err := tfc.Mint(admin.Address(), 1, admin)
	if err != nil {
		t.Fatal(err)
	}

	doneCh, errCh := sdk.WaitForConfirmation(context.Background(), txId, 0)
	select {
	case <-doneCh:
	case err = <-errCh:
		t.Fatal(err)
	}

	balanceNew, err := tfc.BalanceOf(admin.Address())
	if err != nil {
		t.Fatal(err)
	}

	if balanceOld+1 != balanceNew {
		t.Fatal()
	}
}

func TestTFC_Mint_unauthorized(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	acc := sdk.RetrieveAccount("4GFFsdrXDfg5QDd7esPc2TLDzzKinpGA1Gt14pyggi8j")
	tfc := sdk.TFCWithTokenId(tokenId)

	txId, err := tfc.Mint(acc.Address(), 1, acc)
	if err != nil {
		t.Fatal(err)
	}

	doneCh, errCh := sdk.WaitForConfirmation(context.Background(), txId, 0)
	select {
	case tx := <-doneCh:
		t.Fatal(tx)
	case err = <-errCh:
		if err != sdkErr.InvalidCallerTxFailure {
			t.Fatal(err)
		}
	}
}

func TestTFC_BalanceOf(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	account := sdk.RetrieveAccount(testAccountPrivateKey)
	tfc := sdk.TFCWithTokenId(tokenId)

	balance, err := tfc.BalanceOf(account.Address())
	if err != nil {
		t.Fatal(err)
	}
	if balance <= 0 {
		t.Fatal()
	}
}

func TestTFC_Transfer(t *testing.T) {
	sdk, err := New(transport.Endpoint[transport.Testnet], transport.Testnet)
	if err != nil {
		t.Fatal(err)
	}
	admin := sdk.RetrieveAccount(testAccountPrivateKey)
	tfc := sdk.TFCWithTokenId(tokenId)

	balanceOld, err := tfc.BalanceOf(admin.Address())
	if err != nil {
		t.Fatal(err)
	}

	txId, err := tfc.Transfer(admin.Address(), 1, admin)
	if err != nil {
		t.Fatal(err)
	}

	doneCh, errCh := sdk.WaitForConfirmation(context.Background(), txId, 0)
	select {
	case <-doneCh:
	case err = <-errCh:
		t.Fatal(err)
	}

	balanceNew, err := tfc.BalanceOf(admin.Address())
	if err != nil {
		t.Fatal(err)
	}

	if balanceOld != balanceNew {
		t.Fatal()
	}
}
