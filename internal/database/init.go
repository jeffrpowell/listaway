package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/jeffrpowell/listaway/internal/constants"
	_ "github.com/lib/pq"
)

//go:embed init.sql
var initSQL string

func init() {
	// Fetch database connection parameters from environment variables
	dbUser := os.Getenv(constants.ENV_POSTGRES_USER)
	if dbUser == "" {
		dbUser = constants.DB_DEFAULT_USER
	}
	dbPassword := os.Getenv(constants.ENV_POSTGRES_PASSWORD)
	if dbPassword == "" {
		dbPassword = constants.DB_DEFAULT_PASSWORD
	}
	dbHost := os.Getenv(constants.ENV_POSTGRES_HOST)
	if dbHost == "" {
		dbHost = constants.DB_DEFAULT_HOST
	}
	dbName := os.Getenv(constants.ENV_POSTGRES_DB)
	if dbName == "" {
		dbName = constants.DB_DEFAULT_DB
	}

	// Construct the connection string
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)

	// Connect to the PostgreSQL database
	fmt.Printf("Attempting database connection: %s\n", connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Execute SQL queries
	fmt.Println("Running database initialization queries")
	_, err = db.Exec(initSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized successfully")
}
