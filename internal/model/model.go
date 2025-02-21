package model

import (
	"github.com/lexcao/genapi/internal/parser/annotation"
)

type Interface struct {
	Name        string
	Package     string
	Methods     []Method
	Annotations annotation.InterfaceAnnotations
}

type Method struct {
	Name        string
	Interface   string
	Params      []Param
	Results     []Param
	Annotations annotation.MethodAnnotations

	// computed fields
	bindingParams *BindingParams
}

// TODO: add unit tests
// TODO: Introduce a Binder to handle this
// Refactor: Binder + functions
func (m Method) BindingParams() *BindingParams {
	if m.bindingParams == nil {
		bindingParams := &BindingParams{}

		bindedParams := map[string]struct{}{}
		paramsByName := map[string]Param{}
		for _, param := range m.Params {
			paramsByName[param.Name] = param
		}

		for _, query := range m.Annotations.Queries {
			value := query.Value.Escape()
			if param, ok := paramsByName[value]; ok {
				bindedParams[value] = struct{}{}
				bindingParams.Queries = append(bindingParams.Queries, param)
			}
		}

		for _, header := range m.Annotations.Headers {
			for _, value := range header.Values {
				value := value.Escape()
				if param, ok := paramsByName[value]; ok {
					bindedParams[value] = struct{}{}
					bindingParams.Headers = append(bindingParams.Headers, param)
				}
			}
		}

		for _, path := range m.Annotations.RequestLine.PathParams() {
			value := path.Escape()
			if param, ok := paramsByName[value]; ok {
				bindedParams[value] = struct{}{}
				bindingParams.Path = append(bindingParams.Path, param)
			}
		}

		for _, param := range m.Params {
			if _, ok := bindedParams[param.Name]; ok {
				continue
			}

			if param.Type == "context.Context" {
				bindingParams.Context = &param
			} else {
				bindingParams.Body = &param
			}
		}

		m.bindingParams = bindingParams
	}
	return m.bindingParams
}

type BindingParams struct {
	Context *Param
	Body    *Param
	Queries []Param
	Headers []Param
	Path    []Param
}

type Param struct {
	Name string
	Type string
}
