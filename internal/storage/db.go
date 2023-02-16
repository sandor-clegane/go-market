package storage

import (
	"database/sql"
)

const (
	initUsersQuery = "" +
		"CREATE TABLE IF NOT EXISTS users (" +
		"id VARCHAR(255) PRIMARY KEY, " +
		"login VARCHAR(255) UNIQUE NOT NULL, " +
		"password VARCHAR(255) NOT NULL)"
	initOrdersQuery = "" +
		"CREATE TABLE IF NOT EXISTS orders (" +
		"id BIGINT PRIMARY KEY, " +
		"status SMALLINT NOT NULL, " +
		"accrual_amount DECIMAL, " +
		"uploaded_at timestamptz NOT NULL, " +
		"user_id VARCHAR(255) REFERENCES users (id))"
	initWithdrawsQuery = "" +
		"CREATE TABLE IF NOT EXISTS withdraws (" +
		"id BIGINT PRIMARY KEY, " +
		"sum DECIMAL NOT NULL, " +
		"processed_at timestamptz NOT NULL, " +
		"user_id VARCHAR(255) REFERENCES users (id))"
	setTZQuery = "" +
		"set timezone = 'Europe/Moscow'"
)

var db *sql.DB

func ConnectAndInitDB(dbAddress string) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	dbConnection, err := sql.Open("postgres", dbAddress)
	if err != nil {
		return nil, err
	}
	err = initDB(dbConnection)
	if err != nil {
		return nil, err
	}
	return dbConnection, nil
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(initUsersQuery)
	if err != nil {
		return err
	}
	_, err = db.Exec(initOrdersQuery)
	if err != nil {
		return err
	}
	_, err = db.Exec(initWithdrawsQuery)
	if err != nil {
		return err
	}
	_, err = db.Exec(setTZQuery)
	if err != nil {
		return err
	}
	return nil
}
