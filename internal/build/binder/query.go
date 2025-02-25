package binder

import (
	"fmt"
	"net/url"
)

type queryBinding struct {
	binding
}

func (b *queryBinding) Name() string {
	return "Query"
}

func (b *queryBinding) Bind(ctx *context) error {
	if len(ctx.Method.Annotations.Queries) == 0 {
		return nil
	}

	ctx.Method.Interface.Imports.Add(`"net/url"`)
	values := url.Values{}
	queries := ctx.Method.Annotations.Queries

	for _, query := range queries {
		if query.Value.IsVariable() {
			escaped := query.Value.Escape()
			if _, ok := ctx.ParamsByName[escaped]; ok {
				ctx.BindedParams.Add(escaped)
			} else {
				return &ErrNotFound{Type: "query", Value: escaped}
			}
		}
		values.Add(query.Key, query.Value.String())
	}

	result := fmt.Sprintf("%#v", values)

	for _, query := range queries {
		result = replaceVariable(result, query.Value)
	}

	ctx.Method.Bindings.Queries = result
	return nil
}
