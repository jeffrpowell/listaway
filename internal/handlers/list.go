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
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}", middleware.Chain(listHandler, append([]middleware.Middleware{middleware.ListIdOwner("listId")}, middleware.DefaultMiddlewareSlice...)...))
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/edit", middleware.Chain(editListGET, append([]middleware.Middleware{middleware.ListIdOwner("listId")}, middleware.DefaultMiddlewareSlice...)...)).Methods("GET")
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/items", middleware.Chain(listItemsGET, append([]middleware.Middleware{middleware.ListIdOwner("listId")}, middleware.DefaultMiddlewareSlice...)...)).Methods("GET")
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

	// Get user's lists
	lists, err := database.GetLists(userId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Get user's collections
	collections, err := database.GetCollections(userId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)
	listsPage := web.ListsPageParams(r, lists, collections, admin, instanceAdmin)
	web.ListsPage(w, listsPage)
}

/* Create list page */
func createListGET(w http.ResponseWriter, r *http.Request) {
	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)
	params := web.CreateListParams(r, admin, instanceAdmin)
	web.CreateListPage(w, params)
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
		var description string = r.FormValue("description")
		id, err := database.CreateList(userId, listName, description)
		if err != nil {
			http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
			log.Print(err)
			return
		}
		w.Header().Add("Location", fmt.Sprintf("/list/%d", id))
		w.WriteHeader(http.StatusOK)
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
	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)
	editListPageParams := web.EditListParams(r, list, admin, instanceAdmin)
	web.EditListPage(w, editListPageParams)
}

/* Update list */
func listPOST(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	var listParams constants.ListPostParams
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&listParams)
	if err != nil {
		http.Error(w, "Invalid input provided", http.StatusBadRequest)
		log.Print(err)
		return
	}
	//Don't trust client input
	taken, err := database.ListNameTaken(listId, listParams.Name)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if taken {
		http.Error(w, "List name already taken", http.StatusBadRequest)
	} else {
		err = database.UpdateList(listId, listParams)
		if err != nil {
			http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
			log.Print(err)
			return
		}
		w.WriteHeader(http.StatusOK)
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
	w.Header().Add("Location", "/list")
	w.WriteHeader(http.StatusOK)
}

/* View list items page */
func listGET(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	list, err := database.GetList(listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	items, err := database.GetListItems(listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)
	listItemsPage := web.ListItemsPageParams(r, list, items, admin, instanceAdmin)
	web.ListItemsPage(w, listItemsPage)
}

/* List items JSON */
func listItemsGET(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	items, err := database.GetListItems(listId)
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
