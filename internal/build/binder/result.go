package binder

import (
	"fmt"

	"github.com/lexcao/genapi/internal/build/model"
)

type resultBinding struct{}

func (b *resultBinding) Name() string {
	return "Result"
}

var errUnsupportedMultipleResults = &ErrUnsupported{Message: "multiple results, only [T | Response], [error] and [T | Response, error] are supported"}

func (b *resultBinding) Bind(ctx *context) error {
	var results model.BindingResults
	switch len(ctx.Method.Results) {
	case 0:
		return nil
	case 1:
		typ := ctx.Method.Results[0].Type
		switch typ {
		case "error":
			results = model.BindingResults{
				Assignment: "resp, err",
				Statement:  "genapi.HandleResponse0(resp, err)",
			}
		case "genapi.Response":
			results = model.BindingResults{
				Assignment: "resp, _",
				Statement:  "*resp",
			}
		case "*genapi.Response":
			results = model.BindingResults{
				Assignment: "resp, _",
				Statement:  "resp",
			}
		default:
			if disallowedResultPrimitives[typ] {
				return &ErrUnsupported{Message: fmt.Sprintf("type [%s] is not supported as a result type", typ)}
			}
			results = model.BindingResults{
				Assignment: "resp, err",
				Statement:  fmt.Sprintf("genapi.MustHandleResponse[%s](resp, err)", typ),
			}
		}
	case 2:
		type0 := ctx.Method.Results[0].Type
		type1 := ctx.Method.Results[1].Type
		if type1 != "error" {
			return errUnsupportedMultipleResults
		}
		switch type0 {
		case "genapi.Response":
			results = model.BindingResults{
				Assignment: "resp, err",
				Statement:  "*resp, err",
			}
		case "*genapi.Response":
			results = model.BindingResults{
				Assignment: "resp, err",
				Statement:  "resp, err",
			}
		default:
			if disallowedResultPrimitives[type0] {
				return &ErrUnsupported{Message: fmt.Sprintf("type [%s] is not supported as a result type", type0)}
			}
			results = model.BindingResults{
				Assignment: "resp, err",
				Statement:  fmt.Sprintf("genapi.HandleResponse[%s](resp, err)", type0),
			}
		}
	default:
		return errUnsupportedMultipleResults
	}

	ctx.Method.Bindings.Results = &results

	return nil
}

var disallowedResultPrimitives = map[string]bool{
	"context.Context": true,
	"string":          true,
	"bool":            true,
	"byte":            true,
	"complex64":       true,
	"complex128":      true,
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
