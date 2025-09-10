package routes

import (
	"net/http"

	"github.com/mahabubulhasibshawon/todo/internal/handlers"
)

func RegisterTodoRoutes(handler *handlers.TodoHandler) {
    http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
            handler.CreateTodo(w, r)
        case http.MethodGet:
            handler.GetAllTodos(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/todos/get/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            handler.GetTodoByID(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/todos/update/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPut {
            handler.UpdateTodo(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/todos/delete/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodDelete {
            handler.DeleteTodo(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
}
