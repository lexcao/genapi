package genapi

type Option func()

func WithHeader(key, value string) Option {
	return func() {
		panic("not implemented")
	}
}
