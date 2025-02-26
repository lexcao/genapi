package e2e

import (
	"testing"

	"github.com/lexcao/genapi"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	client := genapi.New[HttpBin]()

	t.Run("get", func(t *testing.T) {
		t.Parallel()
		resp, err := client.Get("test")
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "value", resp.Args["key"])
		require.Equal(t, "test", resp.Args["value"])
		require.Equal(t, "world", resp.Headers["X-Hello"])
		require.Equal(t, "global-value", resp.Headers["Global-Header"])
	})

	t.Run("post", func(t *testing.T) {
		t.Parallel()
		body := &Body{Hello: "world"}
		resp, err := client.Post(body)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "global-value", resp.Headers["Global-Header"])
		require.Equal(t, "https://httpbin.org/post", resp.URL)
		require.Equal(t, "world", resp.JSON["hello"])
	})
}
