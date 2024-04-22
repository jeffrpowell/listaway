package main

import (
	"fmt"
	"net/http"

	"github.com/jeffrpowell/listaway/internal/constants"
	_ "github.com/jeffrpowell/listaway/internal/database" //blank import to run init()
	_ "github.com/jeffrpowell/listaway/internal/handlers" //blank import to run init()
)

func main() {
	fmt.Println("####################################")
	fmt.Println("#             LISTAWAY             #")
	fmt.Println("####################################")
	fmt.Println()
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", constants.ROUTER)
}
