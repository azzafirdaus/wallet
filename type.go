package main

import "time"

// Response ...
type Response struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}

// ResponseError ...
type ResponseError struct {
	Error string `json:"error,omitempty"`
}

// ResponseInitAccount ...
type ResponseInitAccount struct {
	Token string                    `json:"token,omitempty"`
	Error *ResponseInitAccountError `json:"error,omitempty"`
}

// ResponseInitAccountError ...
type ResponseInitAccountError struct {
	CustomerXID []string `json:"customer_xid,omitempty"`
}

// ResponseWallet ...
type ResponseWallet struct {
	Wallet ResponseWalletDetail `json:"wallet,omitempty"`
}

// ResponseWalletDetail ...
type ResponseWalletDetail struct {
	ID         string     `json:"id,omitempty"`
	OwnedBy    string     `json:"owned_by,omitempty"`
	Status     string     `json:"status,omitempty"`
	EnabledAt  *time.Time `json:"enabled_at,omitempty"`
	DisabledAt *time.Time `json:"disabled_at,omitempty"`
	Balance    int        `json:"balance"`
}

// ResponseDeposit ...
type ResponseDeposit struct {
	Deposit ResponseDepositDetail `json:"deposit,omitempty"`
}

// ResponseDepositDetail ...
type ResponseDepositDetail struct {
	ID          string    `json:"id,omitempty"`
	DepositedBy string    `json:"deposited_by,omitempty"`
	Status      string    `json:"status,omitempty"`
	DepositedAt time.Time `json:"deposited_at,omitempty"`
	Amount      int       `json:"amount,omitempty"`
	ReferenceID string    `json:"reference_id,omitempty"`
}

// ResponseWithdrawal ...
type ResponseWithdrawal struct {
	Withdrawal ResponseWithdrawalDetail `json:"withdrawal,omitempty"`
}

// ResponseWithdrawalDetail ...
type ResponseWithdrawalDetail struct {
	ID          string    `json:"id,omitempty"`
	WithdrawnBy string    `json:"withdrawn_by,omitempty"`
	Status      string    `json:"status,omitempty"`
	WithdrawnAt time.Time `json:"withdrawn_at,omitempty"`
	Amount      int       `json:"amount,omitempty"`
	ReferenceID string    `json:"reference_id,omitempty"`
}
