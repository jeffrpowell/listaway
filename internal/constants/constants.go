package constants

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Registered environment vars
const ENV_AUTH_KEY string = "LISTAWAY_AUTH_KEY" // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
const ENV_POSTGRES_USER string = "POSTGRES_USER"
const ENV_POSTGRES_PASSWORD string = "POSTGRES_PASSWORD"
const ENV_POSTGRES_HOST string = "POSTGRES_HOST"
const ENV_POSTGRES_DB string = "POSTGRES_DB"

// Database consts
const DB_DEFAULT_USER string = "listaway"
const DB_DEFAULT_PASSWORD string = "listaway"
const DB_DEFAULT_HOST string = "localhost"
const DB_DEFAULT_DB string = "listaway"
const DB_TABLE_LIST string = "listaway.list"
const DB_TABLE_SHARE string = "listaway.share"
const DB_TABLE_ITEM string = "listaway.item"

// Handler consts
const COOKIE_NAME_SESSION string = "session"

var (
	authKey                            = []byte(os.Getenv(ENV_AUTH_KEY))
	COOKIE_STORE *sessions.CookieStore = sessions.NewCookieStore(authKey)
	ROUTER       *mux.Router           = mux.NewRouter()
)

func init() {
	maxAge := 86400 * 7 // 7 days

	COOKIE_STORE.MaxAge(maxAge)
	COOKIE_STORE.Options.Path = "/"
	COOKIE_STORE.Options.HttpOnly = true
	COOKIE_STORE.Options.Secure = false
}
