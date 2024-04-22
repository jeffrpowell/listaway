package handlers

import (
	"fmt"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/constants"
)

func init() {
	constants.ROUTER.HandleFunc("/list/{listId}/share", listShareHandler)
}

func listShareHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		listSharePOST(w, r)
	case "DELETE":
		listShareDELETE(w, r)
	case "GET":
		listShareGET(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

/* Change share link */
func listSharePOST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

/* Unpublish */
func listShareDELETE(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

/* Get share link */
func listShareGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
