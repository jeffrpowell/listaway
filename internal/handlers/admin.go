package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/admin/register", middleware.DefaultPublicMiddlewareChain(registerAdminHandler))
}

func registerAdminHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		registerAdminGET(w, r)
	case "POST":
		registerAdminPOST(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

/* Register admin page */
func registerAdminGET(w http.ResponseWriter, r *http.Request) {
	params := web.RegisterAdminParams(database.AdminUserExists())
	web.RegisterAdmin(w, params)
}

/* Submit new admin user */
func registerAdminPOST(w http.ResponseWriter, r *http.Request) {
	if database.AdminUserExists() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	newUser := constants.UserRegister{
		Email:    r.FormValue("email"),
		Name:     r.FormValue("name"),
		Password: r.FormValue("password"),
		Admin:    true,
	}
	if invalid, reason := newUserIsInvalid(newUser); invalid {
		http.Error(w, reason, http.StatusBadRequest)
		return
	}
	err := database.RegisterUser(newUser)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
	} else {
		constants.ADMIN_EXISTS = true
		// Set user as authenticated
		session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)
		session.Values["authenticated"] = true
		session.Save(r, w)
		w.Header().Add("Location", "/auth")
		w.WriteHeader(http.StatusOK)
	}
}

func newUserIsInvalid(newUser constants.UserRegister) (bool, string) {
	if strings.TrimSpace(newUser.Email) == "" {
		return true, "Email is required"
	}
	if strings.TrimSpace(newUser.Name) == "" {
		return true, "Name is required"
	}
	if strings.TrimSpace(newUser.Password) == "" {
		return true, "Password is required"
	}
	return false, ""
}
