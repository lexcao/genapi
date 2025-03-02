package binder

import (
	"net/textproto"

	"github.com/lexcao/genapi/internal/build/binder/printer"
	"github.com/lexcao/genapi/internal/build/model"
)

type headerBinding struct{}

func (b *headerBinding) Name() string {
	return "Header"
}

func (b *headerBinding) Bind(ctx *context) error {
	if len(ctx.Method.Annotations.Headers) == 0 {
		return nil
	}

	ctx.Method.Interface.Imports.Add(`"net/http"`)
	headers := ctx.Method.Annotations.Headers
	bindedHeaders := map[string]bindedVariablesPrinter{}

	for _, header := range headers {
		for _, variable := range header.Values {
			bindedValue := model.BindedVariable{Variable: variable}
			if variable.IsVariable() {
				escaped := variable.Escape()
				if param, ok := ctx.ParamsByName[escaped]; ok {
					if param.Type != "string" {
						ctx.Method.Interface.Imports.Add(`"strconv"`)
					}

					ctx.BindedParams.Add(escaped)
					bindedValue.Param = &param
				} else {
					return &ErrNotFound{Type: "header", Value: escaped}
				}
			}

			key := textproto.CanonicalMIMEHeaderKey(header.Key)
			bindedHeaders[key] = append(bindedHeaders[key], bindedValue)
		}
	}

	headersValue := printer.Print(bindedHeaderPrinter{
		orderBy: headers,
		binded:  bindedHeaders,
	})

	ctx.Method.Bindings.Headers = headersValue
	return nil
}
