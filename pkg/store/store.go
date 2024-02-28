package store

import (
	"database/sql"
	"errors"
)

var store *sql.DB

func Open() error {
	db, err := sql.Open("sqlite3", "./store.db")
	if err != nil {
		return errors.New("Failed to open database")
	}

	store = db
	return nil
}

func Execute(query string, any ...any) (sql.Result, error) {
	return store.Exec(query, any)
}

func Prepare(query string) (*sql.Stmt, error) {
	return store.Prepare(query)
}

func Query(query string, any ...any) (*sql.Rows, error) {
	return store.Query(query, any)
}

func Close() error {
	return store.Close()
}
