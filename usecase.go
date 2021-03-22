package main

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// InitAccount ...
func InitAccount(userID string) (sessionID string, err error) {
	err = insertUser(database, userID)
	if err != nil {
		log.Println("Error InitAccount insertUser: " + err.Error())
		return
	}

	sessionID = generateSessionID(userID)

	err = createSession(database, userID, sessionID)
	if err != nil {
		log.Println("Error InitAccount createSession: " + err.Error())
		return
	}

	return
}

// EnableWallet ...
func EnableWallet(userID string) (status bool, wallet Wallet, err error) {
	wallet, err = getWalletByUserID(database, userID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error EnableWallet getWalletByUserID: " + err.Error())
		return
	}
	err = nil

	if wallet.ID == "" {
		wallet, err = createWallet(database, userID, defaultBalance)
		if err != nil {
			log.Println("Error EnableWallet createWallet: " + err.Error())
			return
		}
	} else {
		if wallet.Status == 1 {
			log.Println("Info EnableWallet wallet already enabled")
			return
		}

		err = updateWalletStatusByID(database, wallet.ID, statusActive)
		if err != nil {
			log.Println("Error EnableWallet updateWalletStatusByID: " + err.Error())
			return
		}
	}

	status = true
	return
}

// ViewBalance ...
func ViewBalance(userID string) (status bool, wallet Wallet, err error) {
	return viewBalance(userID)
}

func viewBalance(userID string) (status bool, wallet Wallet, err error) {
	wallet, err = getWalletByUserID(database, userID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error ViewBalance getWalletByUserID: " + err.Error())
		return
	}
	err = nil

	if wallet.ID == "" || wallet.Status == 0 {
		log.Println("Info wallet disabled")
		return
	}

	status = true
	return
}

// DisableWallet ...
func DisableWallet(userID string) (status bool, wallet Wallet, err error) {
	wallet, err = getWalletByUserID(database, userID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error DisableWallet getWalletByUserID: " + err.Error())
		return
	}
	err = nil

	if wallet.ID == "" {
		log.Println("Wallet is disabled")
		return

	}

	if wallet.Status == 0 {
		log.Println("Info DisableWallet wallet already enabled")
		return
	}

	err = updateWalletStatusByID(database, wallet.ID, statusInactive)
	if err != nil {
		log.Println("Error DisableWallet updateWalletStatusByID: " + err.Error())
		return
	}

	status = true
	return
}

// Deposit ...
func Deposit(userID, referenceID string, amount int) (status bool, transaction WalletTransaction, err error) {
	status, wallet, err := viewBalance(userID)
	if err != nil {
		log.Println("Error Deposit viewBalance: " + err.Error())
		return
	}

	if !status {
		return
	}

	transaction, err = getTransactionByReferenceID(database, referenceID, depositType)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error Deposit getTransactionByReferenceID: " + err.Error())
		return
	}
	err = nil

	if transaction.ID != "" {
		err = errors.New("Reference id must be unique")
		return
	}

	total := wallet.Balance + amount
	transaction, err = updateBalance(database, wallet.ID, referenceID, amount, total, depositType)
	if err != nil {
		log.Println("Error Deposit updateBalance: " + err.Error())
		return
	}

	status = true
	return
}

// Withdrawal ...
func Withdrawal(userID, referenceID string, amount int) (status bool, transaction WalletTransaction, err error) {
	status, wallet, err := viewBalance(userID)
	if err != nil {
		log.Println("Error Withdrawal viewBalance: " + err.Error())
		return
	}

	if !status {
		return
	}

	if amount > wallet.Balance {
		err = errors.New("Insufficient balance")
		return
	}

	transaction, err = getTransactionByReferenceID(database, referenceID, withdrawalType)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error Withdrawal getTransactionByReferenceID: " + err.Error())
		return
	}
	err = nil

	if transaction.ID != "" {
		err = errors.New("Reference id must be unique")
		return
	}

	total := wallet.Balance - amount
	transaction, err = updateBalance(database, wallet.ID, referenceID, amount, total, withdrawalType)
	if err != nil {
		log.Println("Error Withdrawal updateBalance: " + err.Error())
		return
	}

	status = true
	return
}

func generateSessionID(userID string) (sessionID string) {
	hasher := sha1.New()
	hasher.Write([]byte(userID))

	sessionID = fmt.Sprintf("%x", hasher.Sum(nil))

	return
}
