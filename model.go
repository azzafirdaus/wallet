package main

import "time"

// Wallet ...
type Wallet struct {
	ID         string    `db:"id"`
	UserID     string    `db:"user_id"`
	Balance    int       `db:"balance"`
	Status     int       `db:"status"`
	EnableTime time.Time `db:"enable_time"`
}

// WalletTransaction ...
type WalletTransaction struct {
	ID          string    `db:"id"`
	WalletID    string    `db:"wallet_id"`
	Type        int       `db:"type"`
	Amount      int       `db:"amount"`
	ReferenceID string    `db:"reference_id"`
	CreateTime  time.Time `db:"create_time"`
}
