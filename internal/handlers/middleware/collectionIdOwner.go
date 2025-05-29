package middleware

import (
	"log"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
)

func CollectionIdOwner(pathVarName string) Middleware {

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

			collectionId, err := helper.GetPathVarInt(r, pathVarName)
			if err != nil {
				http.Error(w, "Invalid collectionId supplied in path", http.StatusBadRequest)
				log.Print(err)
				return
			}
			granted, err := database.UserOwnsCollection(userId, collectionId)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}
			if !granted {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
