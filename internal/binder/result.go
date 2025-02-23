package binder

type resultBinding struct {
	binding
}

func (b *resultBinding) Name() string {
	return "Result"
}

func (b *resultBinding) Bind(ctx *context) error {
	return nil
}
