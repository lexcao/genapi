package genapi

import (
	"github.com/lexcao/genapi/internal"
	"github.com/lexcao/genapi/internal/runtime"
)

type Interface = internal.Interface
type Request = internal.Request
type Response = internal.Response
type Config = internal.Config
type HttpClient = internal.HttpClient

func New[T Interface]() T {
	return runtime.New[T]()
}

func Register[api Interface, client Interface](config Config) {
	runtime.Register[api, client](config)
}
