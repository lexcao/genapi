package client

import (
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	httpmock.Activate()
	os.Exit(m.Run())
}
