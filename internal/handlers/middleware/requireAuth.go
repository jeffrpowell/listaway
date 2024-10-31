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

			if !constants.ADMIN_EXISTS {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				http.Redirect(w, r, "/admin/register", http.StatusSeeOther)
				return
			}
			// Check if user has a valid session
			session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
                w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
                w.Header().Set("Pragma", "no-cache")
                w.Header().Set("Expires", "0")
				http.Redirect(w, r, "/auth", http.StatusSeeOther)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
