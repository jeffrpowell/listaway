package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/admin/register", middleware.DefaultPublicMiddlewareChain(registerAdminHandler))
	constants.ROUTER.HandleFunc("/admin/users", middleware.Chain(userAdminGET, append(middleware.DefaultMiddlewareSlice, middleware.RequireAdmin())...)).Methods("GET")
	constants.ROUTER.HandleFunc("/admin/users/create", middleware.Chain(createUserHandler, append(middleware.DefaultMiddlewareSlice, middleware.RequireAdmin())...))
	constants.ROUTER.HandleFunc("/admin/user/{userId:[0-9]+}/listscount", middleware.Chain(userListCountGET, append(middleware.DefaultMiddlewareSlice, middleware.RequireGroupAdmin("userId"))...)).Methods("GET")
	constants.ROUTER.HandleFunc("/admin/user/{userId:[0-9]+}", middleware.Chain(deleteUser, append(middleware.DefaultMiddlewareSlice, middleware.RequireGroupAdmin("userId"))...)).Methods("DELETE")
	constants.ROUTER.HandleFunc("/admin/user/{userId:[0-9]+}/toggleadmin", middleware.Chain(toggleUserAdmin, append(middleware.DefaultMiddlewareSlice, middleware.RequireGroupAdmin("userId"))...)).Methods("POST")
}

func registerAdminHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		registerAdminGET(w, r)
	case "PUT":
		registerAdminPUT(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		createUserGET(w, r)
	case "PUT":
		createUserPUT(w, r)
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
func registerAdminPUT(w http.ResponseWriter, r *http.Request) {
	if database.AdminUserExists() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	newUser := constants.UserRegister{
		GroupId:       0,
		Email:         r.FormValue("email"),
		Name:          r.FormValue("name"),
		Password:      r.FormValue("password"),
		Admin:         true,
		InstanceAdmin: true,
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

/* User Admin page */
func userAdminGET(w http.ResponseWriter, r *http.Request) {
	selfId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	users, err := database.GetUsersInSameGroupAsUser(selfId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	admin := helper.IsUserAdmin(r)
	params := web.UserAdminPageParams(users, selfId, admin)
	web.UserAdminPage(w, params)
}

/* Create user page */
func createUserGET(w http.ResponseWriter, r *http.Request) {
	admin := helper.IsUserAdmin(r)
	params := web.CreateUserParams(admin)
	web.CreateUserPage(w, params)
}

/* Submit new user */
func createUserPUT(w http.ResponseWriter, r *http.Request) {
	var admin = r.FormValue("admin") == "on"
	selfId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	groupId, err := database.GetUserGroupId(selfId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	newUser := constants.UserRegister{
		GroupId:       groupId,
		Email:         r.FormValue("email"),
		Name:          r.FormValue("name"),
		Password:      r.FormValue("password"),
		Admin:         admin,
		InstanceAdmin: false,
	}
	if invalid, reason := newUserIsInvalid(newUser); invalid {
		http.Error(w, reason, http.StatusBadRequest)
		return
	}
	err = database.RegisterUser(newUser)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
	} else {
		w.Header().Add("Location", "/admin/users")
		w.WriteHeader(http.StatusNoContent)
	}
}

func userListCountGET(w http.ResponseWriter, r *http.Request) {
	targetUserId, err := helper.GetPathVarInt(r, "userId")
	if err != nil {
		http.Error(w, "Invalid userId provided in path", http.StatusBadRequest)
		log.Print(err)
		return
	}
	lists, err := database.GetLists(targetUserId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	_, err = fmt.Fprintf(w, "%d", len(lists))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetPathVarInt(r, "userId")
	if err != nil {
		http.Error(w, "Invalid userId supplied in path", http.StatusBadRequest)
		log.Print(err)
		return
	}
	err = database.DeleteUser(userId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toggleUserAdmin(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetPathVarInt(r, "userId")
	if err != nil {
		http.Error(w, "Invalid userId supplied in path", http.StatusBadRequest)
		log.Print(err)
		return
	}
	admin, err := database.UserIsAdmin(userId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	err = database.SetUserAdmin(userId, !admin)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if admin {
		w.Write([]byte("false"))
	} else {
		w.Write([]byte("true"))
	}
}
