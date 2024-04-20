package main

import (
	"fmt"
	"net/http"

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
	handlers.InitAuthHandlers()
	handlers.InitListHandlers()
	handlers.InitShareHandlers()
	handlers.InitItemHandlers()
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
