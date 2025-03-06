package http

import (
	"net/http"
	"testing"

	"github.com/lexcao/genapi/internal"
)

func (c *HttpClient) GetClient() *http.Client {
	return c.client
}

func TestHttpClient(t *testing.T) {
	internal.TestHttpClient(t, func() internal.HttpClientTester { return DefaultClient })
}
