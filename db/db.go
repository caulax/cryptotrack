package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func GetSQLiteDBConnection(dbPath string) (*sql.DB, error) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err // Return nil and the error if it fails
	}

	// Ping the database to ensure the connection is valid
	if err := db.Ping(); err != nil {
		db.Close() // Make sure to close the connection if ping fails
		return nil, err
	}

	// Return the database connection
	return db, nil
}
