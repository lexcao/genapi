package build

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	err := Run(Config{
		Filename: "testdata/github.go",
		Output:   "testdata/github.gen.actual.go",
	})
	require.NoError(t, err)

	actual, err := os.ReadFile("testdata/github.gen.actual.go")
	require.NoError(t, err)

	expect, err := os.ReadFile("testdata/github.gen.expect.go")
	require.NoError(t, err)

	require.Equal(t, expect, actual)
	require.NoError(t, os.Remove("testdata/github.gen.actual.go"))
}
