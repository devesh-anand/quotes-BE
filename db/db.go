package db

import (
	"database/sql"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func GetConnection() (*sql.DB, error) {
	// db, err := sql.Open("mysql", os.Getenv("DSN"))
	// if err != nil {
	// 	log.Fatalf("failed to connect: %v", err)
	// }

	// if err := db.Ping(); err != nil {
	// 	log.Fatalf("failed to ping: %v", err)
	// }
	// log.Println("Successfully connected to PlanetScale!")

	// return db, nil

	//connect to turso db
	url:=os.Getenv("DSN");
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
