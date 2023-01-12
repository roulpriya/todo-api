package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/roulpriya/todo-api/models"
	"github.com/roulpriya/todo-api/repositories"
	"net/http"
	"strconv"
	"time"
)

type CreateTodoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type TodoResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTodoResponse TodoResponse

/**
CreateTodo handles a request to create a Todo in the database
*/

func CreateTodo(repository repositories.TodoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateTodoRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo := models.Todo{
			Title:     request.Title,
			Content:   request.Content,
			Completed: false,
		}

		todo, err := repository.Create(todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, CreateTodoResponse{
			ID:        todo.ID,
			Title:     todo.Title,
			Content:   todo.Content,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt,
			UpdatedAt: todo.UpdatedAt,
		})
	}
}

type FindAllTodosResponse []TodoResponse

func FindAllTodos(repository repositories.TodoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := repository.FindAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response FindAllTodosResponse
		for _, todo := range todos {
			response = append(response, TodoResponse{
				ID:        todo.ID,
				Title:     todo.Title,
				Content:   todo.Content,
				Completed: todo.Completed,
				CreatedAt: todo.CreatedAt,
				UpdatedAt: todo.UpdatedAt,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

type FindTodoByIDResponse TodoResponse

func FindTodoByID(repository repositories.TodoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.FindByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, FindTodoByIDResponse{
			ID:        todo.ID,
			Title:     todo.Title,
			Content:   todo.Content,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt,
			UpdatedAt: todo.UpdatedAt,
		})
	}
}

type UpdateTodoRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
}

type UpdateTodoResponse TodoResponse

func UpdateTodo(repository repositories.TodoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var request UpdateTodoRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.FindByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		todo.Title = request.Title
		todo.Content = request.Content
		todo.Completed = request.Completed
		todo.UpdatedAt = time.Now()

		todo, err = repository.Update(todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, UpdateTodoResponse{
			ID:        todo.ID,
			Title:     todo.Title,
			Content:   todo.Content,
			Completed: todo.Completed,
			CreatedAt: todo.CreatedAt,
			UpdatedAt: todo.UpdatedAt,
		})
	}
}

func DeleteTodoByID(repository repositories.TodoRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.FindByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = repository.Delete(todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "todo deleted"})
	}
}
