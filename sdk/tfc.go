package sdk

import "github.com/virtualeconomy/go-v-sdk/vsys"

type tfc struct {
	contractId string
}

func newTFC(tfcId string) *tfc {
	return &tfc{
		contractId: vsys.TokenId2ContractId(tfcId),
	}
}

func (t *tfc) claimTFC(recipient Address, amount int64, admin *Account) (doneCh chan interface{}, errCh chan error) {
	contract := vsys.Contract{
		Amount: amount,
	}
	issueData := t.contract.BuildIssueData()
tx:
	admin.vsysAcc.BuildExecuteContract(
		t.contractId,
	)
}
