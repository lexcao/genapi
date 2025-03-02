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

	// @GET("/echo/number/{id}")
	// @Query("page", "{page}")
	// @Header("X-Page-Size", "{pageSize}")
	GetEchoNumbers(id int, page int, pageSize int) (*genapi.Response, error)

	// @GET("/echo/boolean/{enabled}")
	// @Query("admin", "{isAdmin}")
	// @Header("X-Debug", "{debug}")
	GetEchoBoolean(enabled bool, isAdmin bool, debug bool) (*genapi.Response, error)

	// @GET("/echo/mixed/{id}/{enabled}")
	// @Query("page", "{page}")
	// @Query("admin", "{isAdmin}")
	// @Header("X-Page-Size", "{pageSize}")
	// @Header("X-Debug", "{debug}")
	GetEchoMixed(id int, enabled bool, page int, isAdmin bool, pageSize int, debug bool) (*genapi.Response, error)
}

type RequestBody struct {
	Message string `json:"message"`
}

type ResponseBody struct {
	Path    string      `json:"path"`
	Body    RequestBody `json:"body"`
	Headers http.Header `json:"headers"`
}
