package helper

import (
	"net/http"
	"strconv"

	"errors"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
)

func GetUserId(r *http.Request) (int, error) {
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)

	// Retrieve our struct and type-assert it
	val := session.Values["userId"]
	if userId, ok := val.(int); ok {
		return userId, nil
	}
	return 0, errors.New("bad userid")
}

func GetPathVarInt(r *http.Request, pathNodeName string) (int, error) {
	return strconv.Atoi(mux.Vars(r)[pathNodeName])
}
