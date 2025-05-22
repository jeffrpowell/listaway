package middleware

import (
	"log"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
)

func RequireInstanceAdmin() Middleware {
	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {
		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			// Check if user is an instance admin
			userId, err := helper.GetUserId(r)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}

			// We need to add a function to check if a user is an instance admin
			isInstanceAdmin, err := database.UserIsInstanceAdmin(userId)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}

			if !isInstanceAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
