package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GetSQLiteDBConnection(dbPath string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %w", err)
	}

	// Check if the connection is valid
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	// Return the database connection
	return db, nil
}
