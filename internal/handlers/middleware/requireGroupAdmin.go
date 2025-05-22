package middleware

import (
	"log"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
)

func RequireGroupAdmin(userIdPathVarName string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			selfUserId, err := helper.GetUserId(r)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}
			admin, err := database.UserIsAdmin(selfUserId)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}
			//Check #1 - is this an admin user?
			if !admin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			targetUserId, err := helper.GetPathVarInt(r, userIdPathVarName)
			if err != nil {
				http.Error(w, "Invalid userId supplied in path", http.StatusBadRequest)
				log.Print(err)
				return
			}
			selfGroupId, err := database.GetUserGroupId(selfUserId)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}
			targetGroupId, err := database.GetUserGroupId(targetUserId)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}
			instanceAdmin, err := database.UserIsInstanceAdmin(selfUserId)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}
			// Check #2 - is this admin in the same group as the target user?
			// Instance admins are an exception to this rule
			if !instanceAdmin && selfGroupId != targetGroupId {
				http.Error(w, "Forbidden", http.StatusForbidden)
				log.Printf("Admin from group %d tried to access user from group %d", selfGroupId, targetGroupId)
				return
			}
			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
