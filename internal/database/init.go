package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/dbtest"
	_ "github.com/lib/pq"
)

//go:embed init.sql
var initSQL string

// GetInitSQL returns the embedded SQL initialization script
func GetInitSQL() string {
	return initSQL
}

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



// getDatabaseConnection returns the appropriate database connection
// If a test database connection has been set, it returns that, otherwise
// it opens a new connection to the main database
func getDatabaseConnection() *sql.DB {
	// Check if we should use the test database
	if testDb := dbtest.GetTestDatabaseConnection(); testDb != nil {
		return testDb
	}
	
	// Use the main database
	db, err := sql.Open("postgres", constants.DB_CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
