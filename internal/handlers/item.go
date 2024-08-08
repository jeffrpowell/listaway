package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/item", middleware.Chain(itemPUT, append(middleware.DefaultMiddlewareSlice, middleware.ListIdOwner("listId"))...)).Methods("PUT")
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/item/create", middleware.Chain(createItemGET, append(middleware.DefaultMiddlewareSlice, middleware.ListIdOwner("listId"))...)).Methods("GET")
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/item/{itemId:[0-9]+}", middleware.Chain(itemHandler, append(middleware.DefaultMiddlewareSlice, middleware.ListIdOwner("listId"))...))
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

/* Create item page */
func createItemGET(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	list, err := database.GetList(listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	web.CreateItemPage(w, web.CreateItemParams(list))
}

/* Create item */
func itemPUT(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	var itemName string = r.FormValue("name")
	var url string = r.FormValue("url")
	priority, err := strconv.ParseInt(r.FormValue("priority"), 10, 32)
	var notes string = r.FormValue("notes")
	id, err := database.CreateItem(constants.ItemInsert{
		Name:     itemName,
		ListId:   uint64(listId),
		URL:      sql.NullString{String: url, Valid: notes != ""},
		Priority: sql.NullInt64{Int64: priority, Valid: err != nil},
		Notes:    sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Add("Status", fmt.Sprint(http.StatusOK))
	w.Header().Add("Location", fmt.Sprintf("/list/%d", id))
	w.Write([]byte(""))
}

/* Update item */
func itemPOST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

/* Delete item */
func itemDELETE(w http.ResponseWriter, r *http.Request) {
	itemId, err := helper.GetPathVarInt(r, "itemId")
	if err != nil {
		http.Error(w, "Invalid itemId supplied in path", http.StatusBadRequest)
		log.Print(err)
		return
	}
	err = database.DeleteItem(itemId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Add("Status", fmt.Sprint(http.StatusNoContent))
	w.Write([]byte(""))
}

/* Get item details */
func itemGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
