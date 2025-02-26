package runtime

import (
	"github.com/lexcao/genapi/internal"
	"github.com/lexcao/genapi/internal/runtime/client"
	"github.com/lexcao/genapi/internal/runtime/registry"
)

func New[T internal.Interface](opts ...Option) T {
	api, config := registry.New[T]()

	// build options
	options := &Options{
		HttpClient: client.DefaultClient,
	}
	if config, ok := config.(internal.Config); ok {
		options.Config = config
	}

	// apply options
	options.apply(opts...)

	// finish initialization
	api.SetHttpClient(options.client())
	return api
}

func Register[api internal.Interface, client internal.Interface](config internal.Config) {
	registry.Register[api, client](config)
}
