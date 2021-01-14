package transport

import (
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"time"
)

type ErrorResponse struct {
	Code    int                    `json:"error"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"-"`
}

// response schema for /node/status
type NodeStatus struct {
	BlockchainHeight int64     `json:"blockchainHeight"`
	StateHeight      int64     `json:"stateHeight"`
	UpdatedTimestamp int64     `json:"updatedTimestamp"`
	UpdatedDate      time.Time `json:"updatedDate"`
}

// response schema for /blocks/height
type LatestHeight struct {
	Height int64 `json:"height"`
}

// response schema for /transactions/info/{txId}
type Transaction struct {
	vsys.Transaction

	// override fields
	Contract map[string]interface{} `json:"contract"`

	Proofs []struct {
		ProofType string `json:"proofType"`
		PublicKey string `json:"publicKey"`
		Address   string `json:"address"`
		Signature string `json:"signature"`
	} `json:"proofs"`

	Status     string `json:"status"`
	FeeCharged int64  `json:"feeCharged"`
	Height     int64  `json:"height"`
}

// response schema for /transactions/list
type TransactionList struct {
	TotalCount   int           `json:"totalCount"`
	Size         int           `json:"size"`
	Transactions []Transaction `json:"transactions"`
}

// response schema for /utils/hash/secure
type SecureCryptographicHash struct {
	Message string `json:"message"`
	Hash    string `json:"hash"`
}

// response schema for /addresses/signText/{address}
type SignedMessage struct {
	Message   string `json:"message"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
}

// response schema for /addresses/balance/{address}
type VsysBalance struct {
	Address       string `json:"address"`
	Confirmations int    `json:"confirmations"`
	Balance       int64  `json:"balance"`
}

// response schema for /contract/balance/{address}/{tokenId}
type TokenBalance struct {
	Address string `json:"address/contractId"`
	TokenId string `json:"tokenId"`
	Balance int64  `json:"balance"`
	Unity   int64  `json:"unity"`
	Height  int64  `json:"height"`
}
