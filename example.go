package main

import (
	"context"
	"fmt"
	"github.com/Troublor/jasmine-vsys-go/sdk"
	"github.com/Troublor/jasmine-vsys-go/sdk/transport"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	const confirmationRequirement = 0

	// initialize SDK object
	var endpoint = transport.Endpoint[transport.Testnet] // the endpoint of testnet and mainnet is predefined in sdk
	var netType = transport.Testnet
	vsysSdk, err := sdk.New(endpoint, netType)
	checkErr(err)

	// create account
	newAcc := vsysSdk.CreateAccount()
	fmt.Println(newAcc.Address(), newAcc.PrivateKey())

	// retrieve account
	adminPrivateKey := "Bokp6eDDyiumnxRVVMvWQmRCK6crc1QA3fNBtuN59ubh"
	admin := vsysSdk.RetrieveAccount(adminPrivateKey)

	// deploy TFC Token contract
	txId, err := vsysSdk.DeployTFC(admin) // send the deploy contract transaction to blockchain, return transaction id
	checkErr(err)

	// wait for enough block confirmation for the transaction
	var tfcContract *sdk.TFC
	doneCh, errCh := vsysSdk.WaitForConfirmation(context.Background(), txId, confirmationRequirement)
	select {
	case tx := <-doneCh:
		// get tfc contract object
		tfcContract = vsysSdk.TFCWithContractId(tx.ContractId)
		fmt.Println("contract deployed")
	case err := <-errCh:
		panic(err)
	}

	fmt.Println(sdk.TFCTotalSupply) // total supply of TFC Token, 20 billion
	fmt.Println(sdk.TFCUnity)       // decimal of TFC Token, 1e8, which means 8 decimals.
	// Note that this is different from Ethereum (18 decimals) due the the limitation of VSystem

	waitForConfirmation := func(txId string) {
		doneCh, errCh = vsysSdk.WaitForConfirmation(context.Background(), txId, confirmationRequirement)
		select {
		case <-doneCh:
			fmt.Println("transaction confirmed")
		case err := <-errCh:
			checkErr(err)
		}
	}

	// Mint TFC Token, can only be called by deployer of contract
	txId, err = tfcContract.Mint(admin.Address(), 1*sdk.TFCUnity, admin) // mint 1 TFC
	checkErr(err)
	// wait for enough block confirmation for the transaction
	waitForConfirmation(txId)

	// Transfer TFC Token
	txId, err = tfcContract.Transfer(newAcc.Address(), 1*sdk.TFCUnity, admin) // transfer 1 TFC
	checkErr(err)
	// wait for enough block confirmation for the transaction
	waitForConfirmation(txId)

	// check TFC balance
	balance, err := tfcContract.BalanceOf(admin.Address())
	fmt.Println(balance) // balance is in the minimal unity (1e-8)

	// Transfer VSYS
	txId, err = vsysSdk.Transfer(newAcc.Address(), 1*sdk.VSYSUnity, admin) // transfer 1 VSYS
	checkErr(err)
	// wait for enough block confirmation for the transaction
	waitForConfirmation(txId)

	// check VSYS balance
	balance, err = vsysSdk.BalanceOf(admin.Address())
	fmt.Println(balance) // balance is in the minimal unity (1e-8)
}
