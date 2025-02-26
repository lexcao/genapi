package it

import (
	"context"
	"net/http"

	"github.com/lexcao/genapi"
)

// @Header("Global-Header", "global-value")
type TestAPI interface {
	genapi.Interface

	// @GET("/echo/{id}")
	// @Header("Echo-Header", "echo-value")
	// @Query("query", "{query}")
	GetEcho(id string, query string) (*genapi.Response, error)

	// @Post("/echo")
	// @Header("Echo-Header", "echo-value")
	// @Query("query", "{query}")
	PostEcho(ctx context.Context, body RequestBody, query string) (*ResponseBody, error)

	// @Post("/echo/error")
	// @Query("status_code", "{statusCode}")
	PostEchoError(statusCode string) error
}

type RequestBody struct {
	Message string `json:"message"`
}

type ResponseBody struct {
	Path    string      `json:"path"`
	Body    RequestBody `json:"body"`
	Headers http.Header `json:"headers"`
}
