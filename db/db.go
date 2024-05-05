package db

import (
	"database/sql"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func GetConnection() (*sql.DB, error) {
	//connect to turso db
	url:=os.Getenv("DSN");
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
