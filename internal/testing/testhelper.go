package testing

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// TestDB provides a database connection for testing
type TestDB struct {
	DB *sql.DB
}

// SetupTestDB creates a test database connection and returns a TestDB instance
func SetupTestDB(t *testing.T) *TestDB {
	// Use environment variables with fallback to test defaults
	testDbHost := getEnvOrDefault("TEST_POSTGRES_HOST", "localhost")
	testDbUser := getEnvOrDefault("TEST_POSTGRES_USER", "listaway")
	testDbPassword := getEnvOrDefault("TEST_POSTGRES_PASSWORD", "password")
	testDbName := getEnvOrDefault("TEST_POSTGRES_DATABASE", "listaway_test")
	testDbPort := getEnvOrDefault("TEST_POSTGRES_PORT", "5432")

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		testDbHost, testDbPort, testDbUser, testDbPassword, testDbName,
	)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	return &TestDB{DB: db}
}

// TeardownTestDB cleans up the database connection
func (tdb *TestDB) TeardownTestDB(t *testing.T) {
	if tdb.DB != nil {
		if err := tdb.DB.Close(); err != nil {
			t.Errorf("Failed to close test database connection: %v", err)
		}
	}
}

// CleanupTables truncates all tables to start with a clean state
func (tdb *TestDB) CleanupTables(t *testing.T) {
	// The test database was created with the 'listaway' schema just like the main database
	_, err := tdb.DB.Exec(`
		TRUNCATE TABLE listaway.list, listaway.item, listaway.user, listaway.reset_tokens, listaway.collection, listaway.collection_list CASCADE
	`)
	if err != nil {
		t.Fatalf("Failed to clean up tables: %v", err)
	}
}

// SetupTestRouter creates a router for HTTP handler testing
func SetupTestRouter() {
	// This will be filled in as we understand more about how the router is configured
	// It might need to use a different router than the production one
}

// helper function to get environment variable with default fallback
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
