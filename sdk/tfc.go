package sdk

import (
	"github.com/Troublor/jasmine-vsys-go/sdk/transport"
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"strconv"
)

const VSYSUnity = 1e8
const TFCUnity = 1e8              // 8 decimal precision
const TFCTotalSupply int64 = 20e9 // 20 billion TFC supply

type TFC struct {
	*transport.Provider
	ContractId string
	TokenId    string
}

func NewTFCWithTokenId(tokenId string, provider *transport.Provider) *TFC {
	return &TFC{
		Provider:   provider,
		ContractId: vsys.TokenId2ContractId(tokenId),
		TokenId:    tokenId,
	}
}

func NewTFCWithContractId(contractId string, tokenIndex int, provider *transport.Provider) *TFC {
	return &TFC{
		Provider:   provider,
		ContractId: contractId,
		TokenId:    vsys.ContractId2TokenId(contractId, tokenIndex),
	}
}

func (t *TFC) Mint(recipient Address, amount int64, admin *Account) (txId string, err error) {
	contract := vsys.Contract{
		Amount: amount,
	}
	funcData := contract.BuildIssueData()
	tx := admin.vsysAcc.BuildExecuteContract(
		t.ContractId,
		vsys.FuncidxIssue,
		funcData,
		recipient+" claims "+strconv.FormatInt(amount, 10)+" TFC",
	)
	var resp vsys.TransactionResponse
	err = t.Post("/contract/broadcast/execute", nil, tx, &resp)
	if err != nil {
		return "", err
	}
	return resp.Id, nil

}

func (t *TFC) BalanceOf(address Address) (int64, error) {
	var balance transport.TokenBalance
	err := t.Get("/contract/balance/"+address+"/"+t.TokenId, nil, &balance)
	if err != nil {
		return 0, err
	}
	return balance.Balance, nil
}

func (t *TFC) Transfer(recipient Address, amount int64, sender *Account) (txId string, err error) {
	tx := sender.vsysAcc.BuildSendTokenTransaction(
		t.TokenId,
		recipient,
		amount,
		false,
		sender.address+" transfer "+strconv.FormatInt(amount, 10)+" to "+recipient)
	var resp vsys.TransactionResponse
	err = t.Post("/contract/broadcast/execute", nil, tx, &resp)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}
