package middleware

import (
	"log"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
)

// ListIdViewer middleware checks if a user can view a list (owns it or it's shared with their group)
func ListIdViewer(pathVarName string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			userId, err := helper.GetUserId(r)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}

			listId, err := helper.GetPathVarInt(r, pathVarName)
			if err != nil {
				http.Error(w, "Invalid listId supplied in path", http.StatusBadRequest)
				log.Print(err)
				return
			}
			
			// Check if user can view the list (owns it OR it's shared with their group)
			canView, err := database.UserCanViewList(userId, listId)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}
			if !canView {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
