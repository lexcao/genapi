package binder

import (
	"github.com/lexcao/genapi/internal/build/binder/printer"
	"github.com/lexcao/genapi/internal/build/model"
)

type pathParamBinding struct{}

func (b *pathParamBinding) Name() string {
	return "PathParam"
}

func (b *pathParamBinding) Bind(ctx *context) error {
	if len(ctx.Method.Annotations.RequestLine.PathParams()) == 0 {
		return nil
	}

	values := map[string]model.BindedVariable{}
	pathParams := ctx.Method.Annotations.RequestLine.PathParams()

	for _, variable := range pathParams {
		bindedValue := model.BindedVariable{Variable: variable}
		escaped := variable.Escape()
		if variable.IsVariable() {
			if param, ok := ctx.ParamsByName[escaped]; ok {
				if param.Type != "string" {
					ctx.Method.Interface.Imports.Add(`"strconv"`)
				}
				ctx.BindedParams.Add(escaped)
				bindedValue.Param = &param
			} else {
				return &ErrNotFound{Type: "path", Value: escaped}
			}
		}
		values[escaped] = bindedValue
	}

	result := printer.PrintWith("map[string]string", func(p *printer.Printer) {
		seen := map[string]bool{}
		for _, variable := range pathParams {
			key := variable.Escape()
			value, ok := values[key]
			if !ok {
				continue
			}
			if seen[key] {
				continue
			}
			seen[key] = true

			p.KeyValueLine(func(p *printer.Printer) {
				p.Quote(key)
			}, func(p *printer.Printer) {
				printVariable(p, value)
			})
		}
	})

	ctx.Method.Bindings.PathParams = result
	return nil
}
