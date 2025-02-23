package binder

import "fmt"

type pathParamBinding struct {
	binding
}

func (b *pathParamBinding) Name() string {
	return "PathParam"
}

func (b *pathParamBinding) Bind(ctx *context) error {
	if len(ctx.Method.Annotations.RequestLine.PathParams()) == 0 {
		return nil
	}

	values := map[string]string{}
	pathParams := ctx.Method.Annotations.RequestLine.PathParams()

	for _, param := range pathParams {
		escaped := param.Escape()
		if param.IsVariable() {
			if _, ok := ctx.ParamsByName[escaped]; ok {
				ctx.BindedParams[escaped] = struct{}{}
			} else {
				return &ErrNotFound{Type: "path", Value: escaped}
			}
		}
		values[escaped] = param.String()
	}

	result := fmt.Sprintf("%#v", values)

	for _, param := range pathParams {
		result = replaceVariable(result, param)
	}

	ctx.Method.Bindings.PathParams = result
	return nil
}
