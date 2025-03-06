package genapi

import (
	"testing"

	"github.com/lexcao/genapi/internal"
)

func TestHttpClient(t *testing.T, createClient func() HttpClientTester) {
	internal.TestHttpClient(t, createClient)
}
