package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/jeffrpowell/listaway/internal/constants"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte(os.Getenv("LISTAWAY_AUTH_KEY"))
	store = sessions.NewCookieStore(key)
)

func GetCookieStore() *sessions.CookieStore {
	return store
}

func RequireAuth() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Check if user has a valid session
			session, _ := store.Get(r, constants.COOKIE_SESSION)
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
