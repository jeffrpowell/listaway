package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
)

func InitAuthHandlers(r *mux.Router) {
	r.HandleFunc("/auth/", authHandler)
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
	session, _ := middleware.GetCookieStore().Get(r, constants.COOKIE_SESSION)

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

/* Logout */
func authDELETE(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.GetCookieStore().Get(r, constants.COOKIE_SESSION)

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
