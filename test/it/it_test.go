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

	t.Run("get echo numbers", func(t *testing.T) {
		resp, err := api.GetEchoNumbers(123, 1, 20)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "20", resp.Header.Get("X-Page-Size"))
		assert.Contains(t, resp.Request.URL.Path, "/echo/number/123")
		assert.Contains(t, resp.Request.URL.RawQuery, "page=1")
	})

	t.Run("get echo boolean", func(t *testing.T) {
		resp, err := api.GetEchoBoolean(true, false, true)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "true", resp.Header.Get("X-Debug"))
		assert.Contains(t, resp.Request.URL.Path, "/echo/boolean/true")
		assert.Contains(t, resp.Request.URL.RawQuery, "admin=false")
	})

	t.Run("get echo mixed", func(t *testing.T) {
		resp, err := api.GetEchoMixed(123, true, 1, false, 20, true)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "20", resp.Header.Get("X-Page-Size"))
		assert.Equal(t, "true", resp.Header.Get("X-Debug"))
		assert.Contains(t, resp.Request.URL.Path, "/echo/mixed/123/true")
		assert.Contains(t, resp.Request.URL.RawQuery, "page=1")
		assert.Contains(t, resp.Request.URL.RawQuery, "admin=false")
	})
}
