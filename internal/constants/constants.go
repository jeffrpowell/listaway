package constants

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Registered environment vars
const (
	ENV_AUTH_KEY          string = "LISTAWAY_AUTH_KEY" // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	ENV_POSTGRES_USER     string = "POSTGRES_USER"
	ENV_POSTGRES_PASSWORD string = "POSTGRES_PASSWORD"
	ENV_POSTGRES_HOST     string = "POSTGRES_HOST"
	ENV_POSTGRES_DB       string = "POSTGRES_DB"
)

// Database consts
const (
	DB_DEFAULT_USER     string = "listaway"
	DB_DEFAULT_PASSWORD string = "listaway"
	DB_DEFAULT_HOST     string = "localhost"
	DB_DEFAULT_DB       string = "listaway"
	DB_TABLE_LIST       string = "listaway.list"
	DB_TABLE_USER       string = "listaway.user"
	DB_TABLE_ITEM       string = "listaway.item"
)

var DB_CONNECTION_STRING string = getDbConnectionString()

// Handler consts
const (
	COOKIE_NAME_SESSION string = "session"
	SHARED_LIST_PATH    string = "sharedlist"
)

var (
	authKey                            = []byte(os.Getenv(ENV_AUTH_KEY))
	COOKIE_STORE *sessions.CookieStore = sessions.NewCookieStore(authKey)
	ROUTER       *mux.Router           = mux.NewRouter()
)

// Random consts
const (
	DefaultN                  = 8
	charSetUnambiguousUpper   = "ABCDEFGHJKLMNPQRTUVWYXZ"
	charSetUnambiguousLower   = "abcdefghjklmnpqrtuvwyxz"
	charSetUnambiguousNumeric = "2346789"
	CharSetUnambiguous        = charSetUnambiguousUpper + charSetUnambiguousLower + charSetUnambiguousNumeric
)

func init() {
	maxAge := 86400 * 7 // 7 days

	COOKIE_STORE.MaxAge(maxAge)
	COOKIE_STORE.Options.Path = "/"
	COOKIE_STORE.Options.HttpOnly = true
	COOKIE_STORE.Options.Secure = false
}

func getDbConnectionString() string {
	// Fetch database connection parameters from environment variables
	dbUser := os.Getenv(ENV_POSTGRES_USER)
	if dbUser == "" {
		dbUser = DB_DEFAULT_USER
	}
	dbPassword := os.Getenv(ENV_POSTGRES_PASSWORD)
	if dbPassword == "" {
		dbPassword = DB_DEFAULT_PASSWORD
	}
	dbHost := os.Getenv(ENV_POSTGRES_HOST)
	if dbHost == "" {
		dbHost = DB_DEFAULT_HOST
	}
	dbName := os.Getenv(ENV_POSTGRES_DB)
	if dbName == "" {
		dbName = DB_DEFAULT_DB
	}

	// Construct the connection string
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
}
