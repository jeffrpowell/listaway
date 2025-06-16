package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
)

// collectionListsPOST handles adding multiple lists to a collection at once
func collectionListsPOST(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		log.Print(err)
		return
	}

	// Get collection ID from the route parameters
	collectionId, _ := helper.GetPathVarInt(r, "collectionId") // Error already checked in middleware
	
	// Get list IDs from form (can be multiple)
	listIds := r.Form["listIds"]
	if len(listIds) == 0 {
		// No lists selected, redirect back to collection detail page
		http.Redirect(w, r, "/collections/"+strconv.Itoa(collectionId), http.StatusSeeOther)
		return
	}

	// Add each list to the collection
	for _, listIdStr := range listIds {
		listId, err := strconv.Atoi(listIdStr)
		if err != nil {
			continue // Skip invalid IDs
		}

		// Add list to collection
		err = database.AddListToCollection(collectionId, listId)
		if err != nil {
			log.Printf("Error adding list %d to collection %d: %v", listId, collectionId, err)
			// Continue with other lists even if one fails
		}
	}

	// Redirect back to collection detail page
	http.Redirect(w, r, "/collections/"+strconv.Itoa(collectionId), http.StatusSeeOther)
}
