package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	"github.com/jeffrpowell/listaway/internal/constants"
	_ "github.com/lib/pq"
)

//go:embed init.sql
var initSQL string

func init() {
	fmt.Printf("Attempting database connection: %s\n", constants.DB_CONNECTION_STRING)
	db := getDatabaseConnection()
	defer db.Close()

	fmt.Println("Running database initialization queries")
	_, err := db.Exec(initSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized successfully")
}

func getDatabaseConnection() *sql.DB {
	db, err := sql.Open("postgres", constants.DB_CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
