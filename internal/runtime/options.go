package runtime

import "github.com/lexcao/genapi/internal"

type Options struct {
	HttpClient internal.HttpClient
	Config     internal.Config
}

func (o *Options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func (o *Options) client() internal.HttpClient {
	client := o.HttpClient
	client.SetConfig(o.Config)
	return client
}

type Option = func(*Options)
