// CODE GENERATED BY genapi. DO NOT EDIT.
package api

import (
	"context"
	"github.com/lexcao/genapi"
	"net/http"
)

type implGitHub struct {
	client genapi.HttpClient
}

// SetHttpClient implments genapi.Interface
func (i *implGitHub) SetHttpClient(client genapi.HttpClient) {
	i.client = client
}

func (i *implGitHub) Contributors(ctx context.Context, owner string, repo string) ([]Contributor, error) {
	resp, err := i.client.Do(&genapi.Request{
		Method: "GET",
		Path:   "/repos/{owner}/{repo}/contributors",
		PathParams: map[string]string{
			"owner": owner,
			"repo":  repo,
		},
		Context: ctx,
	})
	return genapi.HandleResponse[[]Contributor](resp, err)
}

func (i *implGitHub) CreateIssue(ctx context.Context, issue Issue, owner string, repo string) error {
	resp, err := i.client.Do(&genapi.Request{
		Method: "POST",
		Path:   "/repos/{owner}/{repo}/issues",
		PathParams: map[string]string{
			"owner": owner,
			"repo":  repo,
		},
		Context: ctx,
		Body:    issue,
	})
	return genapi.HandleResponse0(resp, err)
}

func init() {
	genapi.Register[GitHub, *implGitHub](
		genapi.Config{
			BaseURL: "https://api.github.com",
			Header: http.Header{
				"Accept": []string{
					"application/vnd.github.v3+json",
				},
			},
		},
	)
}
