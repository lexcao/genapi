package genapi

import (
	"github.com/lexcao/genapi/internal"
	"github.com/lexcao/genapi/internal/runtime"
)

type Option = runtime.Option

func WithHttpClient(httpClient internal.HttpClient) Option {
	return func(options *runtime.Options) {
		options.HttpClient = httpClient
	}
}

func WithConfig(config internal.Config) Option {
	return func(options *runtime.Options) {
		options.Config = config
	}
}

func WithBaseURL(baseURL string) Option {
	return func(options *runtime.Options) {
		options.Config.BaseURL = baseURL
	}
}

func WithHeader(key, value string) Option {
	return func(options *runtime.Options) {
		options.Config.Headers.Add(key, value)
	}
}
