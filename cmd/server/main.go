package main

import (
	"log"
	"net/http"

	"github.com/mahabubulhasibshawon/todo/internal/db"
	"github.com/mahabubulhasibshawon/todo/internal/handlers"
	"github.com/mahabubulhasibshawon/todo/internal/routes"
	"github.com/mahabubulhasibshawon/todo/internal/todo"
)

func main() {
	// connect db
	db.Connect()

	repo := todo.NewPostgresRepository(db.DB)
	handler := handlers.NewTodoHandler(repo)

	routes.RegisterTodoRoutes(handler)

	log.Println("server running at localhost : 8080")
	http.ListenAndServe(":8080",nil)
}