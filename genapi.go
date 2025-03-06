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
type HttpClientTester = internal.HttpClientTester

type Option = runtime.Option
type Options = runtime.Options

func New[T Interface](opts ...Option) T {
	return runtime.New[T](opts...)
}

func Register[api Interface, client Interface](config Config) {
	runtime.Register[api, client](config)
}
