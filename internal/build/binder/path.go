package binder

type pathBinding struct{}

func (b *pathBinding) Name() string {
	return "Path"
}

func (b *pathBinding) Bind(ctx *context) error {
	ctx.Method.Bindings.Path = ctx.Method.Annotations.RequestLine.Path
	return nil
}
