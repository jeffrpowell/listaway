package handlers

import (
	"encoding/json"
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
	constants.ROUTER.HandleFunc("/list", middleware.DefaultMiddlewareChain(listsHandler))
	constants.ROUTER.HandleFunc("/list/create", middleware.DefaultMiddlewareChain(createListGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/list/namecheck", middleware.DefaultMiddlewareChain(nameCheckGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}", middleware.Chain(listHandler, append(middleware.DefaultMiddlewareSlice, middleware.ListIdOwner("listId"))...))
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/edit", middleware.Chain(editListGET, append(middleware.DefaultMiddlewareSlice, middleware.ListIdOwner("listId"))...)).Methods("GET")
}

func listsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listsGET(w, r)
	case "PUT":
		listPUT(w, r)
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

/* Create list */
func listPUT(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	var listName string = r.FormValue("name")
	//Don't trust client input
	taken, err := database.ListNameTaken(userId, listName)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if taken {
		http.Error(w, "List name already taken", http.StatusBadRequest)
	} else {
		id, err := database.CreateList(userId, listName)
		if err != nil {
			http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
			log.Print(err)
			return
		}
		w.Header().Add("Status", fmt.Sprint(http.StatusOK))
		w.Header().Add("Location", fmt.Sprintf("/list/%d", id))
		w.Write([]byte(""))
	}
}

/* Edit list page */
func editListGET(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	list, err := database.GetList(listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	editListPageParams := web.EditListParams(list)
	web.EditListPage(w, editListPageParams)
}

/* Rename list */
func listPOST(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	var listName string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&listName)
	if err != nil {
		http.Error(w, "Invalid input provided", http.StatusBadRequest)
		log.Print(err)
		return
	}
	//Don't trust client input
	taken, err := database.ListNameTaken(listId, listName)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if taken {
		http.Error(w, "List name already taken", http.StatusBadRequest)
	} else {
		err = database.UpdateList(listId, listName)
		if err != nil {
			http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
			log.Print(err)
			return
		}
		w.Header().Add("Status", fmt.Sprint(http.StatusOK))
		w.Write([]byte(""))
	}
}

/* Delete list */
func listDELETE(w http.ResponseWriter, r *http.Request) {
	listId, err := helper.GetPathVarInt(r, "listId")
	if err != nil {
		http.Error(w, "Invalid listId supplied in path", http.StatusBadRequest)
		log.Print(err)
		return
	}
	var confirmationName string
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&confirmationName)
	if err != nil {
		http.Error(w, "Invalid input provided", http.StatusBadRequest)
		log.Print(err)
		return
	}
	deleted, err := database.DeleteList(listId, confirmationName)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if !deleted {
		http.Error(w, "Confirmation name did not match list name", http.StatusBadRequest)
		return
	}
	w.Header().Add("Status", fmt.Sprint(http.StatusOK))
	w.Header().Add("Location", "/list")
	w.Write([]byte(""))
}

/* Get list-level details */
func listGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
