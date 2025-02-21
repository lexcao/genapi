package genapi

import (
	"context"
	"net/http"
	"net/url"
)

type Interface interface {
	setHttpClient(HttpClient)
}

func New[T Interface](opts ...Option) T {
	panic("not implemented")
}

type Request struct {
	Body       any
	Method     string
	Path       string
	Queries    url.Values
	Headers    http.Header
	PathParams map[string]string
	Context    context.Context
}

type Response = http.Response

type HttpClient interface {
	Do(req *Request) (*Response, error)
}
