package middleware

import (
	"net/http"

	"github.com/jeffrpowell/listaway/internal/constants"
)

func RequireAuth() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Check if user has a valid session
			session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
