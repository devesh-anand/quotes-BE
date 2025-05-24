package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func GetConnection() (*sql.DB, error) {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DSN environment variable is not set")
	}

	log.Println("Attempting to connect to database...")
	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close() // Clean up the connection if ping fails
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}
