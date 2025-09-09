package routes

import (
	"net/http"

	"github.com/mahabubulhasibshawon/todo/handlers"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", handlers.GetTodos)
	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTodoByID(w, r)
		case http.MethodPut:
			handlers.UpdateTodo(w, r)
		case http.MethodDelete:
			handlers.DeleteTodo(w, r)
		default:
			http.Error(w, "method not allowed", 405)
		}
	})
	mux.HandleFunc("/create-todo", handlers.CreateTodo)

	return mux
}
