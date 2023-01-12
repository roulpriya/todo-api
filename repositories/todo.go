package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/roulpriya/todo-api/models"
	"time"
)

type TodoRepository struct {
	db *sqlx.DB
}

func NewTodoRepository(db *sqlx.DB) TodoRepository {
	return TodoRepository{db: db}
}

func (t TodoRepository) Create(todo models.Todo) (models.Todo, error) {
	now := time.Now()
	var lastInsertId int64
	err := t.db.QueryRow(
		"INSERT INTO todo (title, content, completed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		todo.Title, todo.Content, todo.Completed, now, now).Scan(&lastInsertId)
	if err != nil {
		return models.Todo{}, err
	}

	return models.Todo{
		ID:        lastInsertId,
		Title:     todo.Title,
		Content:   todo.Content,
		Completed: todo.Completed,
	}, nil
}

func (t TodoRepository) FindAll() ([]models.Todo, error) {
	var todos []models.Todo
	err := t.db.Select(&todos, "SELECT * FROM todo")
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (t TodoRepository) FindByID(id int64) (models.Todo, error) {
	var todo models.Todo
	err := t.db.Get(&todo, "SELECT * FROM todo WHERE id = $1", id)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (t TodoRepository) Update(todo models.Todo) (models.Todo, error) {
	err := t.db.Get(&todo, "UPDATE todo SET title = $1, content = $2, completed = $3, updated_at = $4 WHERE id = $5",
		todo.Title, todo.Content, todo.Completed, time.Now(), todo.ID)
	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (t TodoRepository) Delete(todo models.Todo) error {
	_, err := t.db.Exec("DELETE FROM todo WHERE id = $1", todo.ID)
	return err
}
