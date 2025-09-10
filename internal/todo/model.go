package todo

import "time"

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}
