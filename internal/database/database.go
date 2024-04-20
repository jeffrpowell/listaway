package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Init() {
	// Fetch database connection parameters from environment variables
	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		dbUser = "listaway"
	}
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		dbPassword = "listaway"
	}
	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		dbName = "listaway"
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

	// Read SQL file
	sqlFile, err := os.ReadFile("db/init.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Execute SQL queries
	fmt.Println("Running database initialization queries")
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized successfully")
}
