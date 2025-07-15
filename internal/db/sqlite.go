package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	driver = "sqlite3"
)

type SQLiteContext struct {
	*sql.DB
}

func Connect(database string) (*SQLiteContext, error) {
	db, err := sql.Open(driver, database)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &SQLiteContext{DB: db}, nil
}

func (s *SQLiteContext) WithTx(fn func(tx *sql.Tx) error) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLiteContext) Disconnect() error {
	return s.Close()
}
