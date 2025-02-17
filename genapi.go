package genapi

import (
	"net/http"
)

type Interface interface {
	setHttpClient(HttpClient)
}

func New[T Interface](opts ...Option) T {
	panic("not implemented")
}

type Request = http.Request
type Response = http.Response

type HttpClient interface {
	Do(req *Request) (*Response, error)
}
