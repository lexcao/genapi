package genapi

func NewConfig(opts ...Option) Config {
	var config Config
	for _, opt := range opts {
		opt.config(&config)
	}
	return config
}

func WithBaseURL(baseURL string) Option {
	return funcOption{
		configFn: func(config *Config) {
			config.BaseURL = baseURL
		},
	}
}

func WithHeader(key, value string) Option {
	return funcOption{
		configFn: func(config *Config) {
			config.Headers.Set(key, value)
		},
	}
}

func WithHttpClient(client HttpClient) Option {
	return funcOption{
		applyFn: func(c Interface) {
			c.SetHttpClient(client)
		},
	}
}

// ------------------- Impl Option ------------------
type Option interface {
	apply(Interface)
	config(*Config)
}

type funcOption struct {
	applyFn  func(Interface)
	configFn func(*Config)
}

func (f funcOption) apply(client Interface) {
	if f.applyFn != nil {
		f.applyFn(client)
	}
}

func (f funcOption) config(config *Config) {
	if f.configFn != nil {
		f.configFn(config)
	}
}
