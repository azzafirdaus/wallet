package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"
)

var (
	database *sql.DB
)

func initDB() {
	os.Remove("wallet.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating wallet.db...")
	file, err := os.Create("wallet.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("wallet.db created")

	database, _ = sql.Open("sqlite3", "./wallet.db") // Open the created SQLite File

	createTable(database) // Create Database Tables
}

func createTable(db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error createTable BeginTx: " + err.Error())
		return
	}

	_, err = tx.ExecContext(ctx, createUserTable)
	if err != nil {
		tx.Rollback()
		log.Println("Error createTable ExecContext: " + err.Error())
		return
	}

	_, err = tx.ExecContext(ctx, createSessionTable)
	if err != nil {
		tx.Rollback()
		log.Println("Error createTable ExecContext: " + err.Error())
		return
	}

	_, err = tx.ExecContext(ctx, createWalletTable)
	if err != nil {
		tx.Rollback()
		log.Println("Error createTable ExecContext: " + err.Error())
		return
	}

	_, err = tx.ExecContext(ctx, createTransactionTable)
	if err != nil {
		tx.Rollback()
		log.Println("Error createTable ExecContext: " + err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error createTable Commit: " + err.Error())
		return
	}

	log.Println("table created")
}

func insertUser(db *sql.DB, ID string) (err error) {
	statement, err := db.Prepare(insertUserSQL)
	if err != nil {
		log.Println("Error insertUser Prepare: " + err.Error())
		return
	}

	if ID == "" {
		ID = generateUUID()
	}

	_, err = statement.Exec(ID)
	if err != nil {
		log.Println("Error insertUser Exec: " + err.Error())
	}

	return
}

func createSession(db *sql.DB, userID, sessionID string) (err error) {
	statement, err := db.Prepare(insertSessionSQL)
	if err != nil {
		log.Println("Error insertUser Prepare: " + err.Error())
		return
	}

	_, err = statement.Exec(sessionID, userID, statusActive)
	if err != nil {
		log.Println("Error insertUser Exec: " + err.Error())
	}

	return
}

func getSession(db *sql.DB, sessionID string) (userID string, err error) {
	row := db.QueryRow(
		checkSessionSQL,
		sessionID,
		statusActive,
	)

	err = row.Scan(&userID)
	if err != nil {
		log.Println("Error getSession Scan: " + err.Error())
	}

	return
}

func getWalletByUserID(db *sql.DB, userID string) (wallet Wallet, err error) {
	row := db.QueryRow(
		getWalletByUserIDSQL,
		userID,
	)

	err = row.Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.Status,
		&wallet.EnableTime,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error getWalletByUserID Scan: " + err.Error())
	}

	return
}

func updateWalletStatusByID(db *sql.DB, ID string, status int) (err error) {
	statement, err := db.Prepare(updateWalletStatusByUserIDSQL)
	if err != nil {
		log.Println("Error updateWalletStatusByID Prepare: " + err.Error())
		return
	}

	now := time.Now()

	_, err = statement.Exec(status, ID, now)
	if err != nil {
		log.Println("Error updateWalletStatusByID Exec: " + err.Error())
	}

	return
}

func createWallet(db *sql.DB, userID string, balance int) (wallet Wallet, err error) {
	statement, err := db.Prepare(insertWalletSQL)
	if err != nil {
		log.Println("Error createWallet Prepare: " + err.Error())
		return
	}

	id := generateUUID()

	now := time.Now()
	_, err = statement.Exec(id, userID, balance, statusActive, now)
	if err != nil {
		log.Println("Error createWallet Exec: " + err.Error())
	}

	wallet = Wallet{
		ID:         id,
		UserID:     userID,
		Balance:    balance,
		Status:     statusActive,
		EnableTime: now,
	}

	return
}

func getTransactionByReferenceID(db *sql.DB, referenceID string, transactionType int) (transaction WalletTransaction, err error) {
	row := db.QueryRow(
		getTransactionByReferenceIDSQL,
		referenceID,
		transactionType,
	)

	err = row.Scan(
		&transaction.ID,
		&transaction.WalletID,
		&transaction.Type,
		&transaction.Amount,
		&transaction.ReferenceID,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error getWalletByUserID Scan: " + err.Error())
	}

	return
}

func updateBalance(db *sql.DB, walletID, referenceID string, amount, total, transactionType int) (transaction WalletTransaction, err error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error updateBalance BeginTx: " + err.Error())
		return
	}

	_, err = tx.ExecContext(ctx,
		updateWalletBalanceByIDSQL,
		total,
		walletID,
	)
	if err != nil {
		tx.Rollback()
		log.Println("Error updateBalance ExecContext: " + err.Error())
		return
	}

	now := time.Now()
	transaction = WalletTransaction{
		ID:          generateUUID(),
		WalletID:    walletID,
		Type:        transactionType,
		Amount:      amount,
		ReferenceID: referenceID,
		CreateTime:  now,
	}

	_, err = tx.ExecContext(ctx,
		insertTransactionSQL,
		transaction.ID,
		transaction.WalletID,
		transaction.Type,
		transaction.Amount,
		transaction.ReferenceID,
		transaction.CreateTime,
	)

	if err != nil {
		tx.Rollback()
		log.Println("Error updateBalance ExecContext: " + err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error updateBalance Commit: " + err.Error())
		return
	}

	return
}

func displayStudents(db *sql.DB) {
	row, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		log.Println("Student: ", code, " ", name, " ", program)
	}
}
