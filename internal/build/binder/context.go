package binder

type contextBinding struct {
	binding
}

func (b *contextBinding) Name() string {
	return "Context"
}

func (b *contextBinding) Bind(ctx *context) error {
	for _, param := range ctx.Method.Params {
		if _, ok := ctx.BindedParams[param.Name]; ok {
			continue
		}

		if param.Type == "context.Context" {
			if ctx.Method.Bindings.Context != "" {
				return &ErrDuplicated{Type: "context.Context", Value: param.Name}
			}

			ctx.BindedParams[param.Name] = struct{}{}
			ctx.Method.Bindings.Context = param.Name
		}
	}

	return nil
}
