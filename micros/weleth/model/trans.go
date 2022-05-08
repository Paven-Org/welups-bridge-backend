package model

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	WelTokenFromEth map[string]string = make(map[string]string)
	EthTokenFromWel map[string]string = make(map[string]string)

	EthereumTk = "0x0000000000000000000000000000000000000000"
	WelupsTk   = "W9yD14Nj9j7xAB4dbGeiX9h8unkKHxuTtb"
)

const (
	StatusSuccess = "confirmed"
	StatusUnknown = "unconfirmed"
	StatusPending = "pending"

	RequestDoubleClaimed = "doubleclaimed"
	RequestExpired       = "expired"
	RequestSuccess       = "success"
)

var (
	ErrAlreadyClaimed     = fmt.Errorf("Already claimed")
	ErrRequestPending     = fmt.Errorf("Request pending")
	ErrUnrecognizedStatus = fmt.Errorf("Unrecognized transaction status")
)

type ClaimRequest struct {
	Txid   int64  `db:"tx_id"`
	ReqID  string `db:"request_id"`
	Status string `db:"status"`
}

type WelEthEvent = WelCashinEthTrans
type WelCashinEthTrans struct {
	ID    int64  `json:"id,omitempty" db:"id,omitempty"`
	ReqID string `json:"request_id,omitempty" db:"request_id,omitempty"`

	DepositTxHash string `json:"deposit_tx_hash" db:"deposit_tx_hash"`
	ClaimTxHash   string `json:"claim_tx_hash" db:"claim_tx_hash"`

	WelTokenAddr string `json:"wel_token_addr" db:"wel_token_addr,omitempty"`
	EthTokenAddr string `json:"eth_token_addr" db:"eth_token_addr,omitempty"`

	WelWalletAddr string `json:"wel_wallet_addr,omitempty" db:"wel_wallet_addr"`
	EthWalletAddr string `json:"eth_wallet_addr" db:"eth_wallet_addr"`

	NetworkID string `json:"network_id" db:"network_id"`

	Amount string `json:"amount" db:"amount"`

	Fee string `json:"fee" db:"fee"`

	DepositStatus string `json:"deposit_status" db:"deposit_status"`
	ClaimStatus   string `json:"claim_status" db:"claim_status"`

	DepositAt time.Time    `json:"deposit_at" db:"deposit_at"`
	ClaimAt   sql.NullTime `json:"claim_at" db:"claim_at"`
}

type EthWelEvent = EthCashoutWelTrans
type EthCashoutWelTrans struct {
	ID    int64  `json:"id,omitempty" db:"id,omitempty"`
	ReqID string `json:"request_id,omitempty" db:"request_id,omitempty"`

	DepositTxHash string `json:"deposit_tx_hash" db:"deposit_tx_hash"`
	ClaimTxHash   string `json:"claim_tx_hash" db:"claim_tx_hash"`

	WelTokenAddr string `json:"wel_token_addr" db:"wel_token_addr,omitempty"`
	EthTokenAddr string `json:"eth_token_addr" db:"eth_token_addr,omitempty"`

	WelWalletAddr string `json:"wel_wallet_addr,omitempty" db:"wel_wallet_addr"`
	EthWalletAddr string `json:"eth_wallet_addr" db:"eth_wallet_addr"`

	NetworkID string `json:"network_id" db:"network_id"`

	Amount string `json:"amount" db:"amount"`

	Fee string `json:"fee" db:"fee"`

	DepositStatus string `json:"deposit_status" db:"deposit_status"`
	ClaimStatus   string `json:"claim_status" db:"claim_status"`

	DepositAt time.Time    `json:"deposit_at" db:"deposit_at"`
	ClaimAt   sql.NullTime `json:"claim_at" db:"claim_at"`
}

//----------------------------------------------------------------//
const (
	Tx2TrUnconfirmed = "unconfirmed"
	Tx2TrIsCashin    = "isCashin"
	Tx2TrExpired     = "expired"

	EthCashinWelUnconfirmed = "unconfirmed"
	EthCashinWelConfirmed   = "confirmed"
	EthCashinWelFailed      = "failed"
)

var (
	ErrTx2TreasuryNotFound = fmt.Errorf("Tx to treasury not found")
)

type TxToTreasury struct {
	TxID         string `json:"tx_id" db:"tx_id"`
	FromAddress  string `json:"from_address" db:"from_address"`
	TreasuryAddr string `json:"treasury_address" db:"treasury_address"`
	TokenAddr    string `json:"token_address" db:"token_address"`

	Amount string `json:"amount" db:"amount"`
	TxFee  string `json:"tx_fee" db:"tx_fee"`

	Status string `json:"status" db:"status"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type EthCashinWelTrans struct {
	ID int64 `json:"id,omitempty" db:"id,omitempty"`

	EthTxHash      string `json:"eth_tx_hash,omitempty" db:"eth_tx_hash,omitempty"`
	WelIssueTxHash string `json:"wel_issue_tx_hash" db:"wel_issue_tx_hash,omitempty"`

	EthTokenAddr string `json:"eth_token_addr" db:"eth_token_addr,omitempty"`
	WelTokenAddr string `json:"wel_token_addr" db:"wel_token_addr,omitempty"`

	EthWalletAddr string `json:"eth_wallet_addr" db:"eth_wallet_addr,omitempty"`
	WelWalletAddr string `json:"wel_wallet_addr,omitempty" db:"wel_wallet_addr,omitempty"`

	NetworkID string `json:"network_id" db:"network_id,omitempty"`

	Total         string `json:"total" db:"total,omitempty"`
	Amount        string `json:"amount" db:"amount,omitempty"`
	CommissionFee string `json:"commission_fee" db:"commission_fee,omitempty"`

	Status string `json:"status" db:"status"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IssuedAt  time.Time `json:"issued_at" db:"issued_at,omitempty"`
}
