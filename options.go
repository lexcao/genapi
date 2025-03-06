package genapi

func WithHttpClient(httpClient HttpClient) Option {
	return func(options *Options) {
		options.HttpClient = httpClient
	}
}

func WithConfig(config Config) Option {
	return func(options *Options) {
		options.Config = config
	}
}

func WithBaseURL(baseURL string) Option {
	return func(options *Options) {
		options.Config.BaseURL = baseURL
	}
}

func WithHeader(key, value string) Option {
	return func(options *Options) {
		options.Config.Header.Add(key, value)
	}
}
