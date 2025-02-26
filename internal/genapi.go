package internal

import (
	"context"
	"net/http"
	"net/url"
)

type Interface interface {
	SetHttpClient(HttpClient)
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

type Config struct {
	BaseURL string
	Headers http.Header
}

type HttpClient interface {
	SetConfig(Config)
	Do(req *Request) (*Response, error)
}

type Option interface {
	Apply(Config)
}
