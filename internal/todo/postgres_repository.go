package todo

import (
	"database/sql"
	"errors"
	"fmt"
)

type PostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) Create(todo *Todo) error {
	query := `INSERT INTO todos (title, description, status, completed) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := r.DB.QueryRow(query, todo.Title, todo.Description, todo.Status, todo.Completed).Scan(&todo.ID, &todo.CreatedAt)
	return err
}

func (r *PostgresRepository) GetAll(status string) ([]Todo, error) {
	var rows *sql.Rows
	var err error

	if status == "" {
		rows, err = r.DB.Query(`SELECT id, title, description, status, completed, created_at FROM todos ORDER BY created_at DESC`)
	} else {
		rows, err = r.DB.Query(`SELECT id, title, description, status, completed, created_at FROM todos WHERE status=$1 ORDER BY created_at DESC`, status)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		var t Todo
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Completed, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func (r *PostgresRepository) GetByID(id int) (*Todo, error) {
	var t Todo
	query := `SELECT * FROM todos WHERE id=$1`
	err := r.DB.QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Completed, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresRepository) Update(todo *Todo) error {
	query := `UPDATE todos SET title=$1, description=$2, status=$3, completed=$4 WHERE id=$5`
	res, err := r.DB.Exec(query, todo.Title, todo.Description, todo.Status, todo.Completed, todo.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no todo found with id %d", todo.ID)
	}
	return nil
}

func (r *PostgresRepository) Delete(id int) error {
	res, err := r.DB.Exec(`DELETE FROM todos WHERE id=$1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no todo found with id %d", id)
	}
	return nil
}
