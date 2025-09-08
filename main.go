package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var service = NewTodoService()

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		todos := service.GetAll()
		json.NewEncoder(w).Encode(todos)
	case http.MethodPost:
		var input Todo
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t := service.Create(input.Title, input.Description, input.Completed)
		json.NewEncoder(w).Encode(t)
	default:
		http.Error(w, "method not allowed", 400)
	}
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		t, err := service.GetById(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(t)
	case http.MethodPut:
		var input Todo
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t, err := service.Update(id, input.Title, input.Description, input.Completed)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(t)

	case http.MethodDelete:
		err := service.Delete(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	fmt.Println("==== server running at port 8080: ====")
	http.HandleFunc("/todos", todosHandler)
	http.HandleFunc("/todos/", todoHandler)
	http.ListenAndServe(":8080", nil)
}
