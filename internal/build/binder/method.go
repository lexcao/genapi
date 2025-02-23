package binder

type methodBinding struct {
	binding
}

func (b *methodBinding) Name() string {
	return "Method"
}

func (b *methodBinding) Bind(ctx *context) error {
	ctx.Method.Bindings.Method = ctx.Method.Annotations.RequestLine.Method
	return nil
}
