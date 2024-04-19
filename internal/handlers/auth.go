package handlers

import (
	"fmt"
	"net/http"
)

func InitAuthHandlers() {
	http.HandleFunc("/auth/", authHandler)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		authPOST(w, r)
	case "DELETE":
		authDELETE(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

/* Login */
func authPOST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

/* Logout */
func authDELETE(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
