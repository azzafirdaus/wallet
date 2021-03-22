package main

const (
	createUserTable = `
		CREATE TABLE user (
			id TEXT NOT NULL PRIMARY KEY
		);
	`

	createSessionTable = `
		CREATE TABLE session (
			id TEXT NOT NULL PRIMARY KEY,
			user_id TEXT NOT NULL,
			status INTEGER NOT NULL
		);
	`

	createWalletTable = `
		CREATE TABLE wallet (
			id TEXT NOT NULL PRIMARY KEY,
			user_id TEXT NOT NULL,
			balance INTEGER NOT NULL,
			status INTEGER NOT NULL,
			enable_time DATETIME
		);
	`
	createTransactionTable = `
		CREATE TABLE wallet_transaction (
			id TEXT NOT NULL PRIMARY KEY,
			wallet_id TEXT NOT NULL,
			type INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			reference_id TEXT NOT NULL,
			create_time DATETIME
		);
	`

	insertUserSQL = `
		INSERT INTO user 
			(id) 
		VALUES 
			(?)
		;
	`

	insertSessionSQL = `
		INSERT INTO session 
			(id, user_id, status) 
		VALUES 
			(?,?,?)
		;
	`

	checkSessionSQL = `
		SELECT
			user_id
		FROM
			session
		WHERE 
			id = $1 AND
			status = $2
		;
	`

	insertWalletSQL = `
		INSERT INTO wallet 
			(id, user_id, balance, status, enable_time) 
		VALUES 
			(?,?,?,?,?)
		;
	`

	getWalletByUserIDSQL = `
		SELECT
			id,
			user_id,
			balance,
			status,
			enable_time
		FROM
			wallet
		WHERE
			user_id = $1
	`

	getTransactionByReferenceIDSQL = `
		SELECT
			id,
			wallet_id,
			type,
			amount,
			reference_id
		FROM
			wallet_transaction
		WHERE
			reference_id = $1 AND
			type = $2
	`

	insertTransactionSQL = `
		INSERT INTO wallet_transaction 
			(id, wallet_id, type, amount, reference_id) 
		VALUES 
			(?,?,?,?,?)
		;
	`

	updateWalletStatusByUserIDSQL = `
		UPDATE
			wallet
		SET
			status = $1,
			enable_time = $2
		WHERE
			user_id = $3
	`

	updateWalletBalanceByIDSQL = `
		UPDATE
			wallet
		SET
			balance = $1
		WHERE
			id = $2
	`
)
