package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/share", middleware.Chain(listShareHandler, append([]middleware.Middleware{middleware.ListIdOwner("listId")}, middleware.DefaultMiddlewareSlice...)...))
	constants.ROUTER.HandleFunc("/"+constants.SHARED_LIST_PATH+"/{shareCode}", middleware.DefaultPublicMiddlewareChain(shareGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/"+constants.SHARED_LIST_PATH+"/{shareCode}/items", middleware.DefaultPublicMiddlewareChain(sharedItemsGET)).Methods("GET")
	// New route for nested shared list view within a collection
	constants.ROUTER.HandleFunc("/"+constants.SHARED_COLLECTION_PATH+"/{collectionShareCode}/"+constants.SHARED_LIST_PATH+"/{listShareCode}", middleware.DefaultPublicMiddlewareChain(nestedShareGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/"+constants.SHARED_COLLECTION_PATH+"/{collectionShareCode}/"+constants.SHARED_LIST_PATH+"/{listShareCode}/items", middleware.DefaultPublicMiddlewareChain(nestedSharedItemsGET)).Methods("GET")
}

func listShareHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		listSharePUT(w, r)
	case "DELETE":
		listShareDELETE(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

/* Create share link */
func listSharePUT(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	code, err := database.GenerateShareCode(listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Write([]byte(code))
}

/* Unpublish */
func listShareDELETE(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	err := database.UnpublishShareCode(listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/* View shared list */
func shareGET(w http.ResponseWriter, r *http.Request) {
	shareCode := mux.Vars(r)["shareCode"]
	list, err := database.GetListFromShareCode(shareCode)
	if err != nil {
		if err == sql.ErrNoRows {
			sharedList404Page := web.SharedList404PageParams(shareCode)
			web.SharedList404Page(w, sharedList404Page)
			return
		}
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	items, err := database.GetListItems(int(list.Id))
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)
	sharedListItemsPage := web.SharedListItemsPageParams(r, shareCode, list, items, admin, instanceAdmin)
	web.SharedListItemsPage(w, sharedListItemsPage)
}

/* Shared list items JSON */
func sharedItemsGET(w http.ResponseWriter, r *http.Request) {
	shareCode := mux.Vars(r)["shareCode"]
	list, err := database.GetListFromShareCode(shareCode)
	if err != nil {
		if err == sql.ErrNoRows {
			sharedList404Page := web.SharedList404PageParams(shareCode)
			web.SharedList404Page(w, sharedList404Page)
			return
		}
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	items, err := database.GetListItems(int(list.Id))
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		log.Print(err)
	}
}

/* View shared list with parent collection context */
func nestedShareGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listShareCode := vars["listShareCode"]
	collectionShareCode := vars["collectionShareCode"]

	// Verify that the collection exists
	collection, err := database.GetCollectionFromShareCode(collectionShareCode)
	if err != nil {
		if err == sql.ErrNoRows {
			sharedCollection404Page := web.SharedCollection404PageParams(collectionShareCode)
			web.SharedCollection404Page(w, sharedCollection404Page)
			return
		}
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Get the list from the share code
	list, err := database.GetListFromShareCode(listShareCode)
	if err != nil {
		if err == sql.ErrNoRows {
			sharedList404Page := web.SharedList404PageParams(listShareCode)
			web.SharedList404Page(w, sharedList404Page)
			return
		}
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Verify the list belongs to the collection
	belongs, err := database.ListInCollection(int(collection.Id), int(list.Id))
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if !belongs {
		// If list doesn't belong to collection, return 404
		sharedList404Page := web.SharedList404PageParams(listShareCode)
		web.SharedList404Page(w, sharedList404Page)
		return
	}

	items, err := database.GetListItems(int(list.Id))
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)

	// Render with collection context
	sharedListItemsPage := web.NestedSharedListItemsPageParams(r, listShareCode, collectionShareCode, list, items, admin, instanceAdmin)
	web.SharedListItemsPage(w, sharedListItemsPage)
}

/* Nested shared list items JSON */
func nestedSharedItemsGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listShareCode := vars["listShareCode"]

	// Get the list from the share code
	list, err := database.GetListFromShareCode(listShareCode)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "List not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	items, err := database.GetListItems(int(list.Id))
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		log.Print(err)
	}
}
