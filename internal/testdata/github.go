package testdata

import (
	"context"

	"github.com/lexcao/genapi"
)

// GitHubService provides access to the GitHub API.
// @BaseURL("https://api.github.com")
// @Header("Accept", "application/vnd.github.v3+json")
type GitHub interface {
	genapi.Interface

	// @get("/users/{username}")
	GetUser(ctx context.Context, username string) (*User, error)

	// @get("/users/{username}/repos")
	// @query("sort", "created")
	// @query("direction", "desc")
	// @query("per_page", "{perPage}")
	ListRepositories(ctx context.Context, username string, perPage int) ([]*Repository, error)
}

type User struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Repository represents a GitHub repository.
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
