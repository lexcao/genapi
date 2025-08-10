package testdata

import "github.com/lexcao/genapi"

// TestAPI with valid base URL for e2e validation test
// @BaseURL("https://api.example.com") 
type TestAPI interface {
	genapi.Interface
	
	// @GET("/test")
	GetTest() error
}