package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/roulpriya/todo-api/handlers"
	"github.com/roulpriya/todo-api/repositories"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	db, err := sqlx.Connect("postgres",
		"postgres://postgres:postgres@localhost:5432/todo_api?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	//Run database migrations

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		} else {
			fmt.Println("Database up to date")
		}
	}

	// Create repositories
	todoRepository := repositories.NewTodoRepository(db)

	//some external agent continuously checks server is reachable

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/todos", handlers.CreateTodo(todoRepository))
	r.GET("/todos", handlers.FindAllTodos(todoRepository))
	r.GET("/todos/:id", handlers.FindTodoByID(todoRepository))
	r.PUT("/todos/:id", handlers.UpdateTodo(todoRepository))
	r.DELETE("/todos/:id", handlers.DeleteTodoByID(todoRepository))

	err = r.Run()

	if err != nil {
		log.Fatal(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
