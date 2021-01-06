package sdk

import (
	"context"
	sdkErr "github.com/Troublor/jasmine-vsys-go/sdk/error"
	"github.com/Troublor/jasmine-vsys-go/sdk/transport"
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"strings"
	"time"
)

type SDK struct {
	*transport.Provider
}

func New(endpoint string, netType transport.NetType) (*SDK, error) {
	p, err := transport.NewProvider(endpoint, netType)
	if err != nil {
		return nil, err
	}
	return &SDK{
		Provider: p,
	}, nil
}

func (sdk *SDK) VSYS2MinUnit(amountInVSYS float32) (amountInMinUnit int64) {
	return VSYS2MinUnit(amountInVSYS)
}

func (sdk *SDK) MinUnit2VSYS(amountInMinUnit int64) (amountInVSYS float32) {
	return MinUnit2VSYS(amountInMinUnit)
}

func VSYS2MinUnit(amountInVSYS float32) (amountInMinUnit int64) {
	amountInMinUnit = int64(amountInVSYS * TFCUnity)
	return amountInMinUnit
}

func MinUnit2VSYS(amountInMinUnit int64) (amountInVSYS float32) {
	amountInVSYS = float32(amountInMinUnit) / TFCUnity
	return amountInVSYS
}

func (sdk *SDK) CreateAccount() *Account {
	seed := vsys.GenerateSeed()
	return createAccount(seed, sdk.NetType)
}

func (sdk *SDK) RetrieveAccount(privateKey string) *Account {
	return retrieveAccount(privateKey, sdk.NetType)
}

func (sdk *SDK) DeployTFC(deployer *Account) (txId string, err error) {
	tx := deployer.vsysAcc.BuildRegisterContract(
		vsys.TokenContract,
		TFCTotalSupply,
		TFCUnity,
		"TFC",
		"VSystem blockchain TFC token",
	)
	var resp vsys.TransactionResponse
	err = sdk.Post("/contract/broadcast/register", nil, tx, &resp)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (sdk *SDK) Transfer(recipient Address, amount int64, sender *Account) (txId string, err error) {
	tx := sender.vsysAcc.BuildPayment(recipient, amount, "")
	var resp vsys.TransactionResponse
	err = sdk.Post("/vsys/broadcast/payment", nil, tx, &resp)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (sdk *SDK) BalanceOf(address Address) (int64, error) {
	var balance transport.VsysBalance
	err := sdk.Get("/addresses/balance/"+address, nil, &balance)
	if err != nil {
		return 0, err
	}
	return balance.Balance, nil
}

func (sdk *SDK) GetTransactionInfo(txId string) (tx transport.Transaction, err error) {
	err = sdk.Get("/transactions/info/"+txId, nil, &tx)
	return tx, err
}

func (sdk *SDK) TFCWithTokenId(tokenId string) *TFC {
	return NewTFCWithTokenId(tokenId, sdk.Provider)
}

func (sdk *SDK) TFCWithContractId(contractId string) *TFC {
	return NewTFCWithContractId(contractId, 0, sdk.Provider)
}

func (sdk *SDK) WaitForConfirmation(ctx context.Context, txId string, requiredConfirmationNumber int) (doneCh chan transport.Transaction, errCh chan error) {
	doneCh = make(chan transport.Transaction)
	errCh = make(chan error)
	if requiredConfirmationNumber < 0 {
		errCh <- sdkErr.NewError("requiredConfirmationNumber must be positive")
		return doneCh, errCh
	}
	go func() {
		// polling latest transaction info and height block
		poll := func() (confirmNumber int64, tx transport.Transaction, err error) {
			err = sdk.Get("/transactions/info/"+txId, nil, &tx)
			if err != nil {
				if sdkE, ok := err.(transport.VsysErr); ok {
					if strings.Contains(sdkE.Raw, "Transaction is not in blockchain") {
						return -1, tx, nil
					}
				}
				return confirmNumber, tx, err
			}
			if tx.Status != "Success" {
				return 0, transport.Transaction{}, sdkErr.NewTransactionFailureErr(tx.Status)
			}
			var height transport.LatestHeight
			err = sdk.Get("/blocks/height", nil, &height)
			if err != nil {
				return confirmNumber, tx, err
			}
			confirmNumber = height.Height - tx.Height
			return confirmNumber, tx, err
		}

		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				errCh <- sdkErr.NewError(ctx.Err().Error())
				ticker.Stop()
				return
			case <-ticker.C:
				confirmNumber, tx, err := poll()
				if err != nil {
					errCh <- err
					ticker.Stop()
					return
				}
				if confirmNumber >= int64(requiredConfirmationNumber) {
					doneCh <- tx
					ticker.Stop()
					return
				}
			}
		}
	}()
	return doneCh, errCh
}
