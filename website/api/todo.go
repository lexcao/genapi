package api

import (
	"context"

	"github.com/lexcao/genapi"
)

//go:generate go run github.com/lexcao/genapi/cmd/genapi -file $GOFILE

// @BaseURL("https://jsonplaceholder.typicode.com")
type TodoAPI interface {
	genapi.Interface

	// @GET("/todos")
	GetTodos(ctx context.Context) ([]Todo, error)

	// @GET("/todos/{id}")
	GetTodo(ctx context.Context, id string) (Todo, error)

	// @POST("/todos")
	CreateTodo(ctx context.Context, todo Todo) (Todo, error)

	// @PUT("/todos/{id}")
	UpdateTodo(ctx context.Context, id string, todo Todo) (Todo, error)

	// @DELETE("/todos/{id}")
	DeleteTodo(ctx context.Context, id string) error
}

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
