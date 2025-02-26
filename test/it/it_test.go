package it

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/lexcao/genapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	ts := httptest.NewServer(newServer(t))
	t.Cleanup(ts.Close)

	api := genapi.New[TestAPI](
		genapi.WithBaseURL(ts.URL),
	)

	t.Run("get echo", func(t *testing.T) {
		resp, err := api.GetEcho("123", "query")
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "echo-value", resp.Header.Get("Echo-Header"))
		assert.Equal(t, "global-value", resp.Header.Get("Global-Header"))
	})

	t.Run("post echo", func(t *testing.T) {
		resp, err := api.PostEcho(context.Background(), RequestBody{Message: "hello"}, "query")
		require.NoError(t, err)
		assert.Equal(t, "/echo?query=query", resp.Path)
		assert.Equal(t, "hello", resp.Body.Message)
		assert.Equal(t, "echo-value", resp.Headers.Get("Echo-Header"))
		assert.Equal(t, "global-value", resp.Headers.Get("Global-Header"))
	})

	t.Run("post echo error", func(t *testing.T) {
		err := api.PostEchoError("500")
		require.Error(t, err)
		var e *genapi.Error
		require.ErrorAs(t, err, &e)
		assert.Equal(t, 500, e.Response.StatusCode)
	})
}
