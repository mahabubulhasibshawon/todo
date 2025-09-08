package main

import (
	"errors"
	"sync"
	"time"
)

type TodoService struct {
	todos []Todo
	nextID int
	mu sync.Mutex
}

func NewTodoService() *TodoService {
	return &TodoService{
		todos: []Todo{},
		nextID: 1,
	}
}

func (s *TodoService) GetAll() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Todo{}, s.todos...)
}

func (s *TodoService) GetById(id int) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.todos {
		if t.ID == id {
			return t, nil
		}
	}
	return Todo{}, errors.New("todo not found")
}

func (s *TodoService) Create(title, description string, completed bool) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := Todo{
		ID: s.nextID,
		Title: title,
		Description: description,
		Completed: completed,
		CreatedAT: time.Now(),
	}
	s.nextID++
	s.todos = append(s.todos, t)
	return t
}

func (s *TodoService) Update(id int, title, description string, completed bool) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.todos {
		if t.ID == id {
			s.todos[i].Title = title
			s.todos[i].Description = description
			s.todos[i].Completed = completed
			return s.todos[i],nil
		}
	}
	return Todo{}, errors.New("todo not found")
}

func (s *TodoService) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.todos {
		if t.ID == id {
			s.todos = append(s.todos[:i],s.todos[i+1:]...)
			return nil
		}
	}
	return errors.New("todo not found")
}