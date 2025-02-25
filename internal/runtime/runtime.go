package runtime

import (
	"github.com/lexcao/genapi/internal"
	"github.com/lexcao/genapi/internal/runtime/client"
	"github.com/lexcao/genapi/internal/runtime/registry"
)

func New[T internal.Interface]() T {
	api, config := registry.New[T]()

	httpClient := client.DefaultClient
	if config, ok := config.(internal.Config); ok {
		httpClient.SetConfig(config)
	}

	api.SetHttpClient(httpClient)
	return api
}

func Register[api internal.Interface, client internal.Interface](config internal.Config) {
	registry.Register[api, client](config)
}
