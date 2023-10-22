package storage

import (
	"database/sql"
)

type Storage struct {
	db  *sql.DB
	log logger
}

type logger interface {
	Error(args ...interface{})
}

func New(db *sql.DB, log logger) *Storage {
	return &Storage{db: db, log: log}
}
