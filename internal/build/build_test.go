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

func TestRunWithDefaultFileMode(t *testing.T) {
	outputFile := "testdata/github.gen.filemode.go"
	err := Run(Config{
		Filename: "testdata/github.go",
		Output:   outputFile,
		// FileMode: 0 (default)
	})
	require.NoError(t, err)

	// Check that file was created with default permissions (0600)
	info, err := os.Stat(outputFile)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0600), info.Mode().Perm())

	require.NoError(t, os.Remove(outputFile))
}

func TestRunWithCustomFileMode(t *testing.T) {
	outputFile := "testdata/github.gen.custom.go"
	err := Run(Config{
		Filename: "testdata/github.go",
		Output:   outputFile,
		FileMode: 0644,
	})
	require.NoError(t, err)

	// Check that file was created with custom permissions (0644)
	info, err := os.Stat(outputFile)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0644), info.Mode().Perm())

	require.NoError(t, os.Remove(outputFile))
}

func TestRunWithInvalidBaseURL(t *testing.T) {
	// Test that build fails with our validation error for invalid base URL
	err := Run(Config{
		Filename: "testdata/invalid_baseurl_test.go",
		Output:   "testdata/invalid_baseurl_test.gen.go",
	})
	defer func() {
		os.Remove("testdata/invalid_baseurl_test.gen.go") // cleanup any generated file
	}()
	
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid base URL '://invalid-url-format' in interface TestAPI: parse \"://invalid-url-format\": missing protocol scheme")
}

func TestRunWithValidBaseURL(t *testing.T) {
	// Test that build succeeds with valid base URL
	err := Run(Config{
		Filename: "testdata/valid_baseurl_test.go",
		Output:   "testdata/valid_baseurl_test.gen.actual.go",
	})
	require.NoError(t, err)

	actual, err := os.ReadFile("testdata/valid_baseurl_test.gen.actual.go")
	require.NoError(t, err)

	expect, err := os.ReadFile("testdata/valid_baseurl_test.gen.expect.go")
	require.NoError(t, err)

	require.Equal(t, expect, actual)
	require.NoError(t, os.Remove("testdata/valid_baseurl_test.gen.actual.go"))
}
