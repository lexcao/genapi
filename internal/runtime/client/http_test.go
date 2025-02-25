package client

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/lexcao/genapi/internal"
	"github.com/stretchr/testify/require"
)

func TestHttpClient(t *testing.T) {
	var client = &HttpClient{
		config: internal.Config{
			BaseURL: "https://api.example.com",
		},
		client: http.DefaultClient,
	}

	requireBody := func(t *testing.T, resp *http.Response, expected string) {
		require.Equal(t, 200, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, expected, string(body))
	}

	t.Run("Simple GET", func(t *testing.T) {
		t.Cleanup(httpmock.Reset)

		httpmock.RegisterResponder("GET", "https://api.example.com/users",
			httpmock.NewStringResponder(200, `{"users": [{"id": 1, "name": "John Doe"}]}`))

		resp, err := client.Do(&internal.Request{
			Method: "GET",
			Path:   "/users",
		})

		require.NoError(t, err)
		requireBody(t, resp, `{"users": [{"id": 1, "name": "John Doe"}]}`)
	})

	t.Run("GET with all params", func(t *testing.T) {
		t.Cleanup(httpmock.Reset)

		httpmock.RegisterResponder("GET", "https://api.example.com/users/1?direction=desc&per_page=10&sort=created",
			httpmock.NewStringResponder(200, `{"user": {"id": 1, "name": "John Doe"}}`))

		resp, err := client.Do(&internal.Request{
			Method: "GET",
			Path:   "/users/{id}",
			PathParams: map[string]string{
				"id": "1",
			},
			Queries: url.Values{
				"direction": []string{"desc"},
				"per_page":  []string{"10"},
				"sort":      []string{"created"},
			},
		})

		require.NoError(t, err)
		requireBody(t, resp, `{"user": {"id": 1, "name": "John Doe"}}`)
	})

	t.Run("POST with headers", func(t *testing.T) {
		t.Cleanup(httpmock.Reset)

		httpmock.RegisterResponder("POST", "https://api.example.com/users",
			httpmock.NewStringResponder(200, `{"user": {"id": 1, "name": "John Doe"}}`))

		resp, err := client.Do(&internal.Request{
			Method: "POST",
			Path:   "/users",
			Body:   `{"name": "John Doe"}`,
			Headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
		})

		require.NoError(t, err)
		requireBody(t, resp, `{"user": {"id": 1, "name": "John Doe"}}`)
	})
}
