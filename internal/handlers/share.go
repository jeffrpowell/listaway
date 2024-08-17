package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
)

func init() {
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/share", middleware.Chain(listShareHandler, append(middleware.DefaultMiddlewareSlice, middleware.ListIdOwner("listId"))...))
	constants.ROUTER.HandleFunc(constants.SHARED_LIST_PATH+"/{shareCode}", shareGET).Methods("GET")
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
	w.Header().Add("Status", fmt.Sprint(http.StatusOK))
	w.Write([]byte(code))
}

/* Unpublish */
func listShareDELETE(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

/* View shared list */
func shareGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
