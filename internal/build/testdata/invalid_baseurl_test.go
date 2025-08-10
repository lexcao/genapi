package testdata

import "github.com/lexcao/genapi"

// TestAPI with invalid base URL for e2e validation test
// @BaseURL("://invalid-url-format") 
type TestAPI interface {
	genapi.Interface
	
	// @GET("/test")
	GetTest() error
}