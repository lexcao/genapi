package binder

type contextBinding struct {
	binding
}

func (b *contextBinding) Name() string {
	return "Context"
}

func (b *contextBinding) Bind(ctx *context) error {
	for _, param := range ctx.Method.Params {
		if ctx.BindedParams.Contains(param.Name) {
			continue
		}

		if param.Type == "context.Context" {
			if ctx.Method.Bindings.Context != "" {
				return &ErrDuplicated{Type: "context.Context", Value: param.Name}
			}

			ctx.BindedParams.Add(param.Name)
			ctx.Method.Bindings.Context = param.Name
		}
	}

	return nil
}
