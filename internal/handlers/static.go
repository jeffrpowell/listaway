package handlers

import (
	"net/http"

	"github.com/jeffrpowell/listaway/internal/constants"
)

func init() {
	fs := http.FileServer(http.Dir("assets/")) //build system will generate this directory
	constants.ROUTER.Handle("/static/", http.StripPrefix("/static/", fs))
}
