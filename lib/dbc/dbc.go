package dbc

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	errorHandler "github.com/yonmey/subtracker-api/lib/errorHandler"
)

func DbInit() {
	schema := `CREATE TABLE IF NOT EXISTS subscriptions (
    id INTEGER PRIMARY KEY,
    name TEXT NULL,
    duration INTEGER NULL);`

	db, err := Connect()
	errorHandler.CheckErr(err)

	_, err = db.Exec(schema)
	errorHandler.CheckErr(err)
}

func Connect() (*sqlx.DB, error) {
	return sqlx.Connect("sqlite3", "subscriptions")
}
