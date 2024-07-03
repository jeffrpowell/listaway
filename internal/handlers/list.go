package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/list", middleware.DefaultMiddleware(listsHandler))
	constants.ROUTER.HandleFunc("/list/create", middleware.DefaultMiddleware(createListGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/list/namecheck", middleware.DefaultMiddleware(nameCheckGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}", middleware.DefaultMiddleware(listHandler))
}

func listsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listsGET(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		listPOST(w, r)
	case "DELETE":
		listDELETE(w, r)
	case "GET":
		listGET(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

/* Get all lists of a user */
func listsGET(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	lists, err := database.GetLists(userId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	listsPage := web.ListsPageParams(lists)
	web.ListsPage(w, listsPage)
}

/* Create list page */
func createListGET(w http.ResponseWriter, r *http.Request) {
	web.CreateListPage(w)
}

/* Checks if a list name is taken */
func nameCheckGET(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	taken, err := database.ListNameTaken(userId, r.URL.Query().Get("name"))
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if taken {
		http.Error(w, "List name already taken", http.StatusBadRequest)
	} else {
		w.Write([]byte(""))
	}
}

/* Rename list */
func listPOST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

/* Delete list */
func listDELETE(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

/* Get list-level details */
func listGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
