package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mahabubulhasibshawon/todo/db"
	"github.com/mahabubulhasibshawon/todo/model"
)

// --- CORS ---
func handleCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handlePreflightReq(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(200)
	}
}

func sendData(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// --- Handlers ---
func GetTodos(w http.ResponseWriter, r *http.Request) {
	handleCors(w)
	handlePreflightReq(w, r)
	if r.Method != http.MethodGet {
		http.Error(w, "use GET method", 400)
		return
	}

	rows, err := db.DB.Query("SELECT id, title, description, completed, created_at FROM todos")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	todos := []model.Todo{}
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		todos = append(todos, t)
	}
	sendData(w, todos, 200)
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	handleCors(w)
	handlePreflightReq(w, r)
	if r.Method != http.MethodGet {
		http.Error(w, "use GET method", 400)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", 400)
		return
	}

	var t model.Todo
	err = db.DB.QueryRow("SELECT id, title, description, completed, created_at FROM todos WHERE id=$1", id).
		Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.CreatedAt)
	if err != nil {
		http.Error(w, "todo not found", 404)
		return
	}

	sendData(w, t, 200)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	handleCors(w)
	handlePreflightReq(w, r)
	if r.Method != http.MethodPost {
		http.Error(w, "use POST method", 400)
		return
	}

	var t model.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}

	err := db.DB.QueryRow(
		"INSERT INTO todos (title, description, completed, created_at) VALUES ($1,$2,$3,$4) RETURNING id, created_at",
		t.Title, t.Description, t.Completed, time.Now(),
	).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	sendData(w, t, 201)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	handleCors(w)
	handlePreflightReq(w, r)
	if r.Method != http.MethodPut {
		http.Error(w, "use PUT method", 400)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", 400)
		return
	}

	var t model.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}

	res, err := db.DB.Exec(
		"UPDATE todos SET title=$1, description=$2, completed=$3 WHERE id=$4",
		t.Title, t.Description, t.Completed, id,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "todo not found", 404)
		return
	}

	_ = db.DB.QueryRow("SELECT created_at FROM todos WHERE id=$1", id).Scan(&t.CreatedAt)
	t.ID = id

	sendData(w, t, 200)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	handleCors(w)
	handlePreflightReq(w, r)
	if r.Method != http.MethodDelete {
		http.Error(w, "use DELETE method", 400)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", 400)
		return
	}

	res, err := db.DB.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "todo not found", 404)
		return
	}

	w.WriteHeader(204)
}
