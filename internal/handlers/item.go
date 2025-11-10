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
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/item", middleware.Chain(itemPUT, append([]middleware.Middleware{middleware.ListIdOwner("listId")}, middleware.DefaultMiddlewareSlice...)...)).Methods("PUT")
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/item/create", middleware.Chain(createItemGET, append([]middleware.Middleware{middleware.ListIdOwner("listId")}, middleware.DefaultMiddlewareSlice...)...)).Methods("GET")
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/item/{itemId:[0-9]+}", middleware.Chain(itemHandler, append([]middleware.Middleware{middleware.ListIdOwner("listId")}, middleware.DefaultMiddlewareSlice...)...))
	constants.ROUTER.HandleFunc("/list/{listId:[0-9]+}/item/{itemId:[0-9]+}/edit", middleware.Chain(itemEditGET, append([]middleware.Middleware{middleware.ListIdOwner("listId")}, middleware.DefaultMiddlewareSlice...)...)).Methods("GET")
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		itemPOST(w, r)
	case "DELETE":
		itemDELETE(w, r)
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
	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)
	web.CreateEditItemPage(w, web.CreateItemParams(r, list, admin, instanceAdmin))
}

/* Create item */
func itemPUT(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	
	// Check if user can edit this list
	canEdit, err := database.UserCanEditList(userId, listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if !canEdit {
		http.Error(w, "Forbidden - you don't have permission to edit this list", http.StatusForbidden)
		return
	}
	
	var itemName string = r.FormValue("name")
	var url string = r.FormValue("url")
	priority, err := strconv.ParseInt(r.FormValue("priority"), 10, 64)
	var notes string = r.FormValue("notes")
	err = database.CreateItem(constants.ItemInsert{
		Name:     itemName,
		ListId:   uint64(listId),
		URL:      sql.NullString{String: url, Valid: url != ""},
		Priority: sql.NullInt64{Int64: priority, Valid: err == nil},
		Notes:    sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Add("Location", fmt.Sprintf("/list/%d", listId))
	w.WriteHeader(http.StatusNoContent)
}

/* Edit item page */
func itemEditGET(w http.ResponseWriter, r *http.Request) {
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	itemId, err := helper.GetPathVarInt(r, "itemId")
	if err != nil {
		http.Error(w, "Invalid itemId supplied", http.StatusBadRequest)
		log.Print(err)
		return
	}
	list, err := database.GetList(listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	item, err := database.GetItem(itemId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)
	web.CreateEditItemPage(w, web.EditItemParams(r, list, item, admin, instanceAdmin))
}

/* Update item */
func itemPOST(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	
	// Check if user can edit this list
	canEdit, err := database.UserCanEditList(userId, listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if !canEdit {
		http.Error(w, "Forbidden - you don't have permission to edit this list", http.StatusForbidden)
		return
	}
	
	itemId, err := helper.GetPathVarInt(r, "itemId")
	if err != nil {
		http.Error(w, "Invalid itemId supplied", http.StatusBadRequest)
		log.Print(err)
		return
	}
	var itemName string = r.FormValue("name")
	var url string = r.FormValue("url")
	priority, err := strconv.ParseInt(r.FormValue("priority"), 10, 64)
	var notes string = r.FormValue("notes")
	err = database.UpdateItem(itemId, constants.ItemInsert{
		Name:     itemName,
		ListId:   uint64(listId),
		URL:      sql.NullString{String: url, Valid: url != ""},
		Priority: sql.NullInt64{Int64: priority, Valid: err == nil},
		Notes:    sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Add("Location", fmt.Sprintf("/list/%d", listId))
	w.WriteHeader(http.StatusNoContent)
}

/* Delete item */
func itemDELETE(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	
	listId, _ := helper.GetPathVarInt(r, "listId") //err will trip in listIdOwner middleware first
	
	// Check if user can edit this list
	canEdit, err := database.UserCanEditList(userId, listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if !canEdit {
		http.Error(w, "Forbidden - you don't have permission to edit this list", http.StatusForbidden)
		return
	}
	
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
	w.WriteHeader(http.StatusNoContent)
}
