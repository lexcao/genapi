package e2e

import "github.com/lexcao/genapi"

//go:generate go run github.com/lexcao/genapi/cmd/genapi -file $GOFILE

// @BaseURL("https://httpbin.org")
// @Header("Global-Header", "global-value")
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
	JSON    map[string]any    `json:"json"`
}
