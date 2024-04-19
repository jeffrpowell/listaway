package handlers

import (
	"fmt"
	"net/http"
)

func InitItemHandlers() {
	http.HandleFunc("GET /list/{listId}/item", itemGET)
	http.HandleFunc("/list/{listId}/item/{itemId}", itemHandler)
}

/* Get all items in a list */
func listItemGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
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
