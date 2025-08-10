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
	// Create test file with invalid base URL
	testFile := "testdata/invalid_baseurl_test.go"
	testContent := `package testdata

import "github.com/lexcao/genapi"

// TestAPI with invalid base URL for e2e validation test
// @BaseURL("://invalid-url-format") 
type TestAPI interface {
	genapi.Interface
	
	// @GET("/test")
	GetTest() error
}`

	err := os.WriteFile(testFile, []byte(testContent), 0600)
	require.NoError(t, err)
	defer func() {
		os.Remove(testFile)
		os.Remove("testdata/invalid_baseurl_test.gen.go") // cleanup any generated file
	}()

	// Test that build fails with our validation error
	err = Run(Config{
		Filename: testFile,
		Output:   "testdata/invalid_baseurl_test.gen.go",
	})
	
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid base URL")
	require.Contains(t, err.Error(), "://invalid-url-format")
	require.Contains(t, err.Error(), "TestAPI")
	require.Contains(t, err.Error(), "missing protocol scheme")
}

func TestRunWithValidBaseURL(t *testing.T) {
	// Create test file with valid base URL
	testFile := "testdata/valid_baseurl_test.go"
	testContent := `package testdata

import "github.com/lexcao/genapi"

// TestAPI with valid base URL for e2e validation test
// @BaseURL("https://api.example.com") 
type TestAPI interface {
	genapi.Interface
	
	// @GET("/test")
	GetTest() error
}`

	err := os.WriteFile(testFile, []byte(testContent), 0600)
	require.NoError(t, err)
	defer func() {
		os.Remove(testFile)
		os.Remove("testdata/valid_baseurl_test.gen.go")
	}()

	// Test that build succeeds with valid base URL
	err = Run(Config{
		Filename: testFile,
		Output:   "testdata/valid_baseurl_test.gen.go",
	})
	
	require.NoError(t, err)
	
	// Verify generated file contains correct base URL
	generated, err := os.ReadFile("testdata/valid_baseurl_test.gen.go")
	require.NoError(t, err)
	
	generatedContent := string(generated)
	require.Contains(t, generatedContent, `BaseURL: "https://api.example.com"`)
	require.Contains(t, generatedContent, "implTestAPI")
}
