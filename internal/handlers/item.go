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
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/item", middleware.Chain(itemsGET, append(middleware.DefaultMiddlewareSlice, middleware.ListIdOwner("listId"))...)).Methods("GET")
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/item/{itemId:[0-9]+}", middleware.DefaultMiddlewareChain(itemHandler))
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		itemPOST(w, r)
	case "DELETE":
		itemDELETE(w, r)
	case "GET":
		itemGET(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

/* Get all items in a list */
func itemsGET(w http.ResponseWriter, r *http.Request) {
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
	listItemsPage := web.ListItemsPageParams(list, items)
	web.ListItemsPage(w, listItemsPage)
}

/* Update item */
func itemPOST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

/* Delete item */
func itemDELETE(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

/* Get item details */
func itemGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
