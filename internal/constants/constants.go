package constants

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Registered environment vars
const (
	ENV_AUTH_KEY          string = "LISTAWAY_AUTH_KEY" // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	ENV_PORT              string = "PORT"
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
	defaultPort         string = "8080"
)

var PORT string = loadEnvWithDefault(ENV_PORT, defaultPort)

var (
	authKey                            = []byte(os.Getenv(ENV_AUTH_KEY))
	COOKIE_STORE *sessions.CookieStore = sessions.NewCookieStore(authKey)
	ROUTER       *mux.Router           = mux.NewRouter()
	ADMIN_EXISTS                       = false
)

// Random consts
const (
	DefaultN                  = 8
	charSetUnambiguousUpper   = "ABCDEFGHJKLMNPQRTUVWYXZ"
	charSetUnambiguousLower   = "abcdefghjklmnpqrtuvwyxz"
	charSetUnambiguousNumeric = "2346789"
	CHARSET_UNAMBIGUOUS       = charSetUnambiguousUpper + charSetUnambiguousLower + charSetUnambiguousNumeric
)

func init() {
	maxAge := 86400 * 7 // 7 days

	COOKIE_STORE.MaxAge(maxAge)
	COOKIE_STORE.Options.Path = "/"
	COOKIE_STORE.Options.HttpOnly = true
	COOKIE_STORE.Options.Secure = false
	COOKIE_STORE.Options.SameSite = http.SameSiteLaxMode
}

func loadEnvWithDefault(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}

func getDbConnectionString() string {
	// Fetch database connection parameters from environment variables
	dbUser := loadEnvWithDefault(ENV_POSTGRES_USER, DB_DEFAULT_USER)
	dbPassword := loadEnvWithDefault(ENV_POSTGRES_PASSWORD, DB_DEFAULT_PASSWORD)
	dbHost := loadEnvWithDefault(ENV_POSTGRES_HOST, DB_DEFAULT_HOST)
	dbName := loadEnvWithDefault(ENV_POSTGRES_DB, DB_DEFAULT_DB)

	// Construct the connection string
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
}
