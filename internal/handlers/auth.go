package handlers

import (
	"net/http"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/", middleware.DefaultMiddleware(rootHandler))
	constants.ROUTER.HandleFunc("/auth/", authHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Redirect(w, r, "/lists/", http.StatusPermanentRedirect)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		authGET(w, r)
	case "POST":
		authPOST(w, r)
	case "DELETE":
		authDELETE(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

/* Login page */
func authGET(w http.ResponseWriter, r *http.Request) {
	web.LoginPage(w)
}

/* Login */
func authPOST(w http.ResponseWriter, r *http.Request) {
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
	http.Redirect(w, r, "/list/", http.StatusSeeOther)
}

/* Logout */
func authDELETE(w http.ResponseWriter, r *http.Request) {
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/auth/", http.StatusSeeOther)
}
