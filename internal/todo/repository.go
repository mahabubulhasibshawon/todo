package todo



type Repository interface {
	Create(todo *Todo) error
	GetAll(status string) ([]Todo, error)
	GetByID(id int) (*Todo, error)
	Update(todo *Todo) error
	Delete(id int) error
}