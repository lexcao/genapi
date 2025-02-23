package binder

import (
	"fmt"
	"net/http"
)

type headerBinding struct {
	binding
}

func (b *headerBinding) Name() string {
	return "Header"
}

func (b *headerBinding) Bind(ctx *context) error {
	if len(ctx.Method.Annotations.Headers) == 0 {
		return nil
	}

	values := http.Header{}
	headers := ctx.Method.Annotations.Headers

	for _, header := range headers {
		for _, value := range header.Values {
			if value.IsVariable() {
				escaped := value.Escape()
				if _, ok := ctx.ParamsByName[escaped]; ok {
					ctx.BindedParams[escaped] = struct{}{}
				} else {
					return &ErrNotFound{Type: "header", Value: escaped}
				}
			}
			values.Add(header.Key, value.String())
		}
	}

	result := fmt.Sprintf("%#v", values)

	for _, header := range headers {
		for _, value := range header.Values {
			result = replaceVariable(result, value)
		}
	}

	ctx.Method.Bindings.Headers = result
	return nil
}
