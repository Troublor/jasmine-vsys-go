package sdk

import (
	"github.com/virtualeconomy/go-v-sdk/vsys"
)

type NetType vsys.NetType

type SDK struct {
	netType NetType
}

func New(endpoint string, netType NetType) *SDK {
	vsys.InitApi(endpoint, vsys.NetType(netType))
	return &SDK{}
}

func (sdk *SDK) VSYS2MinUnit(amountInVSYS float32) (amountInMinUnit int64) {
	return VSYS2MinUnit(amountInVSYS)
}

func (sdk *SDK) MinUnit2VSYS(amountInMinUnit int64) (amountInVSYS float32) {
	return MinUnit2VSYS(amountInMinUnit)
}

func VSYS2MinUnit(amountInVSYS float32) (amountInMinUnit int64) {
	amountInMinUnit = int64(amountInVSYS * 1e8)
	return amountInMinUnit
}

func MinUnit2VSYS(amountInMinUnit int64) (amountInVSYS float32) {
	amountInVSYS = float32(amountInMinUnit) / 1e8
	return amountInVSYS
}

func (sdk *SDK) CreateAccount() *Account {
	seed, _ := generateRandomString(64)
	return createAccount(seed, vsys.NetType(sdk.netType))
}

func (sdk *SDK) RetrieveAccount(privateKey string) *Account {
	return retrieveAccount(privateKey, vsys.NetType(sdk.netType))
}

func (sdk *SDK) DeployTFCSync(deployer *Account) (tfcAddress Address, err error) {
	tfcAddressCh, errCh := sdk.DeployTFC(deployer)
	select {
	case addr := <-tfcAddressCh:
		return addr, nil
	case err := <-errCh:
		return "", err
	}
}

func (sdk *SDK) DeployTFC(deployer *Account) (tfcAddressCh chan Address, errCh chan error) {
	tfcAddressCh = make(chan Address, 1)
	errCh = make(chan error, 1)
	go func() {
		tx := deployer.vsysAcc.BuildRegisterContract(
			vsys.TokenContract,
			2000000000,
			18,
			"TFC",
			"VSystem blockchain TFC token",
		)
		resp, err := vsys.SendRegisterContractTx(tx)
		if err != nil {
			errCh <- err
		} else {
			tfcAddressCh <- resp.Recipient
		}
	}()
	return tfcAddressCh, errCh
}
