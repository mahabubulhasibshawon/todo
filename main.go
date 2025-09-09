package main

import (
	"fmt"
	"net/http"
	"github.com/mahabubulhasibshawon/todo/db"
	"github.com/mahabubulhasibshawon/todo/routes"
)

func main() {
	// Connect DB
	db.Connect()

	mux := routes.SetupRoutes()

	fmt.Println("==== Server running on :8080 ====")
	http.ListenAndServe(":8080", mux)
}
