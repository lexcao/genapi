package api

import (
	"context"

	"github.com/lexcao/genapi"
)

//go:generate go run ../../../cmd/genapi/main.go -file $GOFILE

// @BaseURL("https://api.github.com")
// @Header("Accept", "application/vnd.github.v3+json")
type GitHub interface {
	genapi.Interface

	// @GET("/repos/{owner}/{repo}/contributors")
	Contributors(ctx context.Context, owner string, repo string) ([]Contributor, error)

	// @POST("/repos/{owner}/{repo}/issues")
	CreateIssue(ctx context.Context, issue Issue, owner string, repo string) error
}

type Contributor struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}

type Issue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Assignees []string `json:"assignees"`
	Milestone int      `json:"milestone"`
	Labels    []string `json:"labels"`
}
