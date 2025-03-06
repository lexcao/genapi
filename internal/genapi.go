package internal

import (
	"context"
	"net/http"
	"net/url"
)

// Interface is the interface that mark this interface should be a genapi client that can be generated
type Interface interface {
	SetHttpClient(HttpClient)
}

type Request struct {
	Body       any
	Method     string
	Path       string
	Queries    url.Values
	Header     http.Header
	PathParams map[string]string
	Context    context.Context
}

type Response = http.Response

type Config struct {
	BaseURL string
	Header  http.Header
}

// HttpClient is the genapi client runtime, the generated client will call this interface to send requests
// You can provide your own implementation of HttpClient to use a different HTTP client
type HttpClient interface {
	// SetConfig sets the config for the client, global config
	SetConfig(Config)

	// Do sends a request and returns a response
	Do(req *Request) (*Response, error)
}

type Option interface {
	Apply(Config)
}
