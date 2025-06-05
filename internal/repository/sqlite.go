package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	driver = "sqlite3"
)

type SQLContext struct {
	*sql.DB
}

func Connect(database string) (*SQLContext, error) {
	db, err := sql.Open(driver, database)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &SQLContext{DB: db}, nil
}

func (s *SQLContext) Disconnect() error {
	return s.Close()
}
