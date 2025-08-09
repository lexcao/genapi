package runtime

import (
	"github.com/lexcao/genapi/internal"
	"github.com/lexcao/genapi/internal/runtime/registry"
	"github.com/lexcao/genapi/pkg/clients/http"
)

func New[T internal.Interface](opts ...Option) (T, error) {
	api, config, err := registry.New[T]()
	if err != nil {
		var zero T
		return zero, err
	}

	// build options
	options := &Options{
		HttpClient: http.DefaultClient,
	}
	if config, ok := config.(internal.Config); ok {
		options.Config = config
	}

	// apply options
	options.apply(opts...)

	// finish initialization
	api.SetHttpClient(options.client())
	return api, nil
}

func Register[api internal.Interface, client internal.Interface](config internal.Config) {
	registry.Register[api, client](config)
}
