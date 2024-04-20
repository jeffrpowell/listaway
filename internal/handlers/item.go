package handlers

import (
	"fmt"
	"net/http"
)

func InitItemHandlers() {
	http.HandleFunc("/list/{listId}/item", itemsHandler)
	http.HandleFunc("/list/{listId}/item/{itemId}", itemHandler)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		itemsGET(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
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
	fmt.Fprintf(w, "Hello, World!")
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
