package handlers

import (
	"log"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/", middleware.DefaultMiddlewareChain(rootHandler)).Methods("GET")
	constants.ROUTER.HandleFunc("/auth", middleware.DefaultPublicMiddlewareChain(authHandler))
	constants.ADMIN_EXISTS = database.AdminUserExists()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/list", http.StatusPermanentRedirect)
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

	userId, err := database.LoginUser(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if userId == -1 {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Values["userId"] = userId
	session.Save(r, w)
	w.Header().Add("Location", "/list")
	w.WriteHeader(http.StatusOK)
}

/* Logout */
func authDELETE(w http.ResponseWriter, r *http.Request) {
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)

	// Revoke users authentication
	session.Values["authenticated"] = false
	delete(session.Values, "userId")
	session.Options.MaxAge = -1
	session.Save(r, w)
	w.Header().Add("Location", "/auth")
	w.WriteHeader(http.StatusNoContent)
}
