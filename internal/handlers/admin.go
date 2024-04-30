package handlers

import (
	"log"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/admin/register", registerAdminHandler)
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
	if newUserIsInvalid(newUser) {
		http.Error(w, "Invalid user provided", http.StatusBadRequest)
		return
	}
	err := database.RegisterUser(newUser)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
	} else {
		// Set user as authenticated
		session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)
		session.Values["authenticated"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/list/", http.StatusSeeOther)
	}

}

func newUserIsInvalid(newUser constants.UserRegister) bool {
	return newUser.Email == "" || newUser.Name == "" || newUser.Password == ""
}
