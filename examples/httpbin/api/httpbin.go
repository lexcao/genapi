package api

import "github.com/lexcao/genapi"

//go:generate go run ../../../cmd/genapi/main.go -file $GOFILE

// @BaseURL("https://httpbin.org")
type HttpBin interface {
	genapi.Interface

	// @GET("/get")
	// @Query("key", "value")
	// @Query("value", "{value}")
	// @Header("X-Hello", "world")
	Get(value string) (*Response, error)

	// @POST("/post")
	Post(body *Body) (*Response, error)
}

type Body struct {
	Hello string `json:"hello"`
}

type Response struct {
	Args    map[string]string `json:"args"`
	Headers map[string]string `json:"headers"`
	URL     string            `json:"url"`
}
