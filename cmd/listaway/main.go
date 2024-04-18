package main

import (
	"fmt"
	"net/http"

	"github.com/jeffrpowell/listaway/pkg/database"
)

func main() {
	database.Init()
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
