package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers"
)

func main() {
	fmt.Println("####################################")
	fmt.Println("#             LISTAWAY             #")
	fmt.Println("####################################")
	fmt.Println()
	database.Init()
	fmt.Println("Standing up web server")
	constants.Init()
	r := mux.NewRouter()
	handlers.InitAuthHandlers(r)
	handlers.InitListHandlers(r)
	handlers.InitShareHandlers(r)
	handlers.InitItemHandlers(r)
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
