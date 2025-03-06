package resty

import (
	"net/http"
	"testing"

	"github.com/lexcao/genapi"
)

func (c *HttpClient) GetClient() *http.Client {
	return c.client.GetClient()
}

func TestHttpClient(t *testing.T) {
	genapi.TestHttpClient(t, func() genapi.HttpClientTester { return DefaultClient })
}
