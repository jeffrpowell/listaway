package dbtest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	_ "github.com/lib/pq"
)

// TestDB provides a database connection for testing
type TestDB struct {
	DB *sql.DB
}

// Global variable to store the test database connection
var (
	testDbConnection *sql.DB
	testDbMutex      sync.RWMutex
)

// SetTestDatabase sets the test database connection to be used by all database functions during tests
func SetTestDatabase(db *sql.DB) {
	testDbMutex.Lock()
	defer testDbMutex.Unlock()
	testDbConnection = db
}

// ClearTestDatabase clears the test database connection
func ClearTestDatabase() {
	testDbMutex.Lock()
	defer testDbMutex.Unlock()
	testDbConnection = nil
}

// IsUsingTestDatabase returns true if a test database connection is currently set
func IsUsingTestDatabase() bool {
	testDbMutex.RLock()
	defer testDbMutex.RUnlock()
	return testDbConnection != nil
}

// GetTestDatabaseConnection returns the current test database connection if available
func GetTestDatabaseConnection() *sql.DB {
	testDbMutex.RLock()
	defer testDbMutex.RUnlock()
	return testDbConnection
}

// SetupTestDB creates a test database connection and returns a TestDB instance
func SetupTestDB(t *testing.T, initSQL string) *TestDB {
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
	fmt.Printf("Attempting test database connection: %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	fmt.Println("Running test database initialization queries")
	_, err = db.Exec(initSQL)
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	fmt.Println("Test database initialized successfully")

	// Set this connection as the test database connection
	SetTestDatabase(db)

	return &TestDB{DB: db}
}

// TeardownTestDB cleans up the database connection
func (tdb *TestDB) TeardownTestDB(t *testing.T) {
	// Clear the test database connection first
	ClearTestDatabase()

	// Then close the connection
	if tdb.DB != nil {
		if err := tdb.DB.Close(); err != nil {
			t.Errorf("Failed to close test database connection: %v", err)
		}
	}
}

// CleanupTables truncates all tables to start with a clean state
func (tdb *TestDB) CleanupTables(t *testing.T) {
	// The test database was created with the 'listaway' schema
	_, err := tdb.DB.Exec(`
		TRUNCATE TABLE listaway.list, listaway.item, listaway.user, listaway.reset_tokens, listaway.collection, listaway.collection_list CASCADE
	`)
	if err != nil {
		t.Fatalf("Failed to clean up tables: %v", err)
	}
}

// Helper function to get environment variable with default fallback
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
