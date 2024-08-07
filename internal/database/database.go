package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(name string) (*Database, error) {
	conn, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	db := &Database{conn}
	if err := db.init(); err != nil {
		defer db.Close()
		return nil, err
	}

	return db, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

func (db *Database) init() error {
	if err := db.initWorkItemsTable(); err != nil {
		return err
	}

	return nil
}
