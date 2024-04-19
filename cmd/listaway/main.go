package main

import (
	"fmt"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers"
)

func main() {
	database.Init()
	handlers.InitAuthHandlers()
	handlers.InitItemHandlers()
	handlers.InitShareHandlers()
	handlers.InitItemHandlers()
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
