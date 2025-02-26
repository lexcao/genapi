package binder

type bodyBinding struct{}

func (b *bodyBinding) Name() string {
	return "Body"
}

func (b *bodyBinding) Bind(ctx *context) error {
	for _, param := range ctx.Method.Params {
		if ctx.BindedParams.Contains(param.Name) {
			continue
		}

		if disallowedBodyPrimitives[param.Type] {
			continue
		}

		if ctx.Method.Bindings.Body != "" {
			return &ErrDuplicated{Type: "body", Value: param.Name}
		}

		ctx.BindedParams.Add(param.Name)
		ctx.Method.Bindings.Body = param.Name
	}

	return nil
}

var disallowedBodyPrimitives = map[string]bool{
	"context.Context": true,
	"bool":            true,
	"byte":            true,
	"complex64":       true,
	"complex128":      true,
	"error":           true,
	"float32":         true,
	"float64":         true,
	"int":             true,
	"int8":            true,
	"int16":           true,
	"int32":           true,
	"int64":           true,
	"rune":            true,
	"uint":            true,
	"uint8":           true,
	"uint16":          true,
	"uint32":          true,
	"uint64":          true,
	"uintptr":         true,
}
