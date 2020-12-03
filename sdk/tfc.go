package sdk

import (
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"strconv"
)

const TFCUnity = 1e8
const TFCTotalSupply = 20e16

type tfc struct {
	contractId string
	tokenId    string
}

func newTFC(tokenId string) *tfc {
	return &tfc{
		contractId: vsys.TokenId2ContractId(tokenId),
		tokenId:    tokenId,
	}
}

func (t *tfc) ClaimTFCSync(recipient Address, amount int64, admin *Account) error {
	doneCh, errCh := t.ClaimTFC(recipient, amount, admin)
	select {
	case <-doneCh:
		return nil
	case err := <-errCh:
		return err
	}
}

func (t *tfc) ClaimTFC(recipient Address, amount int64, admin *Account) (doneCh chan interface{}, errCh chan error) {
	doneCh = make(chan interface{}, 1)
	errCh = make(chan error, 1)
	go func() {
		contract := vsys.Contract{
			Amount: amount,
		}
		funcData := contract.BuildIssueData()
		tx := admin.vsysAcc.BuildExecuteContract(
			t.contractId,
			vsys.FuncidxIssue,
			funcData,
			recipient+" claims "+strconv.FormatInt(amount, 10)+" TFC",
		)
		resp, err := vsys.SendExecuteContractTx(tx)
		if err != nil {
			errCh <- err
			return
		}
		if resp.Error != 0 {
			errCh <- vsysResponseError(resp.Error)
			return
		}
		close(doneCh)
	}()
	return doneCh, errCh
}

func (t *tfc) BalanceOf(address *Account) (balance int64, err error) {
	bal, err := address.vsysAcc.GetTokenBalance(t.tokenId)
	if err != nil {
		return 0, err
	}
	if bal.Error != 0 {
		return 0, vsysResponseError(bal.Error)
	}
	balance = bal.Balance
	return balance, nil
}

func (t *tfc) TransferSync(recipient Address, amount int64, sender *Account) error {
	doneCh, errCh := t.Transfer(recipient, amount, sender)
	select {
	case <-doneCh:
		return nil
	case err := <-errCh:
		return err
	}
}

func (t *tfc) Transfer(recipient Address, amount int64, sender *Account) (doneCh chan interface{}, errCh chan error) {
	doneCh = make(chan interface{}, 1)
	errCh = make(chan error, 1)
	go func() {
		tx := sender.vsysAcc.BuildSendTokenTransaction(
			t.tokenId,
			recipient,
			amount,
			false,
			sender.address+" transfer "+strconv.FormatInt(amount, 10)+" to "+recipient)
		resp, err := vsys.SendExecuteContractTx(tx)
		if err != nil {
			errCh <- err
			return
		}
		if resp.Error != 0 {
			errCh <- vsysResponseError(resp.Error)
			return
		}
		close(doneCh)
	}()
	return doneCh, errCh
}
