package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "github.com/mahabubulhasibshawon/todo/internal/todo"
)

type TodoHandler struct {
    Repo todo.Repository
}

// Constructor for TodoHandler
func NewTodoHandler(repo todo.Repository) *TodoHandler {
    return &TodoHandler{Repo: repo}
}

//////////////////////////////////////////////////
// POST /todos — Create a new todo
//////////////////////////////////////////////////
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    var t todo.Todo
    err := json.NewDecoder(r.Body).Decode(&t)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if t.Status == "" {
        t.Status = "due" // Default status
    }

    err = h.Repo.Create(&t)
    if err != nil {
        http.Error(w, "Failed to create todo: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(t)
}

//////////////////////////////////////////////////
// GET /todos — Get all todos (optional ?status=)
//////////////////////////////////////////////////
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
    status := r.URL.Query().Get("status")

    todos, err := h.Repo.GetAll(status)
    if err != nil {
        http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

//////////////////////////////////////////////////
// GET /todos/get/{id}
//////////////////////////////////////////////////
func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
    id, err := getIDFromPath(r, "/todos/get/")
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    todoItem, err := h.Repo.GetByID(id)
    if err != nil {
        http.Error(w, "DB error", http.StatusInternalServerError)
        return
    }
    if todoItem == nil {
        http.NotFound(w, r)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todoItem)
}

//////////////////////////////////////////////////
// PUT /todos/update/{id}
//////////////////////////////////////////////////
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
    id, err := getIDFromPath(r, "/todos/update/")
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var t todo.Todo
    err = json.NewDecoder(r.Body).Decode(&t)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    t.ID = id // Ensure the correct ID is set

    err = h.Repo.Update(&t)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Todo updated"})
}

//////////////////////////////////////////////////
// DELETE /todos/delete/{id}
//////////////////////////////////////////////////
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
    id, err := getIDFromPath(r, "/todos/delete/")
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    err = h.Repo.Delete(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // 204 No Content
}

//////////////////////////////////////////////////
// Helper function to extract ID from URL path
//////////////////////////////////////////////////
func getIDFromPath(r *http.Request, prefix string) (int, error) {
    idStr := strings.TrimPrefix(r.URL.Path, prefix)
    return strconv.Atoi(idStr)
}
