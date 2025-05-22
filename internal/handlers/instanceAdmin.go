package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	constants.ROUTER.HandleFunc("/admin/allusers", middleware.Chain(allUsersGET, append(middleware.DefaultMiddlewareSlice, middleware.RequireInstanceAdmin())...)).Methods("GET")
	constants.ROUTER.HandleFunc("/admin/user/{userId:[0-9]+}/toggleinstanceadmin", middleware.Chain(toggleUserInstanceAdmin, append(middleware.DefaultMiddlewareSlice, middleware.RequireInstanceAdmin())...)).Methods("POST")
}

// All Users page - for Instance Administrators
func allUsersGET(w http.ResponseWriter, r *http.Request) {
	selfId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Instance admin can see all users
	users, err := database.GetAllUsers()
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	isInstanceAdmin := true // We're in RequireInstanceAdmin middleware, so this is always true
	params := web.AllUsersPageParams(users, selfId, isInstanceAdmin)
	web.AllUsersPage(w, params)
}

// Toggle Instance Admin status
func toggleUserInstanceAdmin(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetPathVarInt(r, "userId")
	if err != nil {
		http.Error(w, "Invalid userId supplied in path", http.StatusBadRequest)
		log.Print(err)
		return
	}

	// Don't allow users to demote themselves
	selfId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if userId == selfId {
		http.Error(w, "Cannot change your own instance admin status", http.StatusForbidden)
		return
	}

	instanceAdmin, err := database.UserIsInstanceAdmin(userId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	err = database.SetUserInstanceAdmin(userId, !instanceAdmin)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Return the new status
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strconv.FormatBool(!instanceAdmin)))
}
