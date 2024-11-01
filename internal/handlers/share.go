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
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/share", middleware.Chain(listShareHandler, append(middleware.DefaultMiddlewareSlice, middleware.ListIdOwner("listId"))...))
	constants.ROUTER.HandleFunc("/"+constants.SHARED_LIST_PATH+"/{shareCode}", middleware.DefaultPublicMiddlewareChain(shareGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/"+constants.SHARED_LIST_PATH+"/{shareCode}/items", middleware.DefaultPublicMiddlewareChain(sharedItemsGET)).Methods("GET")
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
	sharedListItemsPage := web.SharedListItemsPageParams(shareCode, list, items, admin)
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
