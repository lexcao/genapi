package testdata

import (
	"context"

	"github.com/lexcao/genapi"
)

// GitHub API
// @BaseURL("https://api.github.com")
// @Header("Accept", "application/vnd.github.v3+json")
type GitHub interface {
	genapi.Interface
	// ListRepositories list repositories for specific user
	// @GET("/repos/{owner}")
	// @Query("sort", "{sort}")
	ListRepositories(ctx context.Context, owner string, sort string) ([]Repository, error)

	// ListContributors list contributors for specific repository
	// @GET("/repos/{owner}/{repo}/contributors")
	ListContributors(ctx context.Context, owner string, repo string) ([]Contributor, error)

	// CreateIssue create issue for specific repository
	// @POST("/repos/{owner}/{repo}/issues")
	CreateIssue(ctx context.Context, issue Issue, owner string, repo string) error
}

type Repository struct{}
type Contributor struct{}
type Issue struct{}
