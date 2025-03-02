package binder

import (
	"github.com/lexcao/genapi/internal/build/binder/printer"
	"github.com/lexcao/genapi/internal/build/model"
)

type queryBinding struct{}

func (b *queryBinding) Name() string {
	return "Query"
}

func (b *queryBinding) Bind(ctx *context) error {
	if len(ctx.Method.Annotations.Queries) == 0 {
		return nil
	}

	ctx.Method.Interface.Imports.Add(`"net/url"`)
	queries := ctx.Method.Annotations.Queries
	values := map[string]bindedVariablesPrinter{}

	for _, query := range queries {
		variable := query.Value
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
				return &ErrNotFound{Type: "query", Value: escaped}
			}
		}
		values[query.Key] = append(values[query.Key], bindedValue)
	}

	result := printer.PrintWith("url.Values", func(p *printer.Printer) {
		seen := map[string]bool{}
		for _, query := range queries {
			key := query.Key
			values, ok := values[key]
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
				p.Item(values)
			})
		}
	})

	ctx.Method.Bindings.Queries = result
	return nil
}
