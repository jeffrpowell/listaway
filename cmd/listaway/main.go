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
	if !constants.ADMIN_EXISTS {
		fmt.Println("No admin user has been created yet. Please register one at http://localhost:8080/admin/register")
	}
	http.ListenAndServe(":8080", constants.ROUTER)
}
