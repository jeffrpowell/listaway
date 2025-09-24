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

	// SMTP configuration for password reset emails
	ENV_SMTP_HOST     string = "SMTP_HOST"
	ENV_SMTP_PORT     string = "SMTP_PORT"
	ENV_SMTP_USER     string = "SMTP_USER"
	ENV_SMTP_PASSWORD string = "SMTP_PASSWORD"
	ENV_SMTP_FROM     string = "SMTP_FROM"
	ENV_SMTP_SECURE   string = "SMTP_SECURE" // true for TLS/SSL, false for unencrypted
	ENV_APP_URL       string = "APP_URL"     // base URL for the application (for reset links)

	// OIDC configuration
	ENV_OIDC_ENABLED        string = "OIDC_ENABLED"        // true/false to enable OIDC authentication
	ENV_OIDC_PROVIDER_URL   string = "OIDC_PROVIDER_URL"   // OIDC provider URL (e.g., https://accounts.google.com)
	ENV_OIDC_CLIENT_ID      string = "OIDC_CLIENT_ID"      // OAuth2 client ID
	ENV_OIDC_CLIENT_SECRET  string = "OIDC_CLIENT_SECRET"  // OAuth2 client secret
	ENV_OIDC_REDIRECT_URL   string = "OIDC_REDIRECT_URL"   // OAuth2 redirect URL
	ENV_OIDC_SCOPES         string = "OIDC_SCOPES"         // OAuth2 scopes (space-separated)
)

// Database consts
const (
	DB_DEFAULT_USER          string = "listaway"
	DB_DEFAULT_PASSWORD      string = "listaway"
	DB_DEFAULT_HOST          string = "localhost"
	DB_DEFAULT_DB            string = "listaway"
	DB_TABLE_LIST            string = "listaway.list"
	DB_TABLE_USER            string = "listaway.user"
	DB_TABLE_ITEM            string = "listaway.item"
	DB_TABLE_RESET           string = "listaway.reset_tokens"
	DB_TABLE_COLLECTION      string = "listaway.collection"
	DB_TABLE_COLLECTION_LIST string = "listaway.collection_list"
)

var DB_CONNECTION_STRING string = getDbConnectionString()

// SMTP configuration with defaults
var (
	SMTP_HOST     string = loadEnvWithDefault(ENV_SMTP_HOST, "")
	SMTP_PORT     string = loadEnvWithDefault(ENV_SMTP_PORT, "587")
	SMTP_USER     string = loadEnvWithDefault(ENV_SMTP_USER, "")
	SMTP_PASSWORD string = loadEnvWithDefault(ENV_SMTP_PASSWORD, "")
	SMTP_FROM     string = loadEnvWithDefault(ENV_SMTP_FROM, "noreply@listaway.dev")
	SMTP_SECURE   string = loadEnvWithDefault(ENV_SMTP_SECURE, "true")
	APP_URL       string = loadEnvWithDefault(ENV_APP_URL, "http://localhost:8080")
)

// OIDC configuration with defaults
var (
	OIDC_ENABLED        string = loadEnvWithDefault(ENV_OIDC_ENABLED, "false")
	OIDC_PROVIDER_URL   string = loadEnvWithDefault(ENV_OIDC_PROVIDER_URL, "")
	OIDC_CLIENT_ID      string = loadEnvWithDefault(ENV_OIDC_CLIENT_ID, "")
	OIDC_CLIENT_SECRET  string = loadEnvWithDefault(ENV_OIDC_CLIENT_SECRET, "")
	OIDC_REDIRECT_URL   string = loadEnvWithDefault(ENV_OIDC_REDIRECT_URL, "")
	OIDC_SCOPES         string = loadEnvWithDefault(ENV_OIDC_SCOPES, "openid profile email")
)

// Handler consts
const (
	COOKIE_NAME_SESSION    string = "session"
	SHARED_LIST_PATH       string = "sharedlist"
	SHARED_COLLECTION_PATH string = "sharedcollection"
	defaultPort            string = "8080"
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
