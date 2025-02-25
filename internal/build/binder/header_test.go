package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
)

func TestBindHeader(t *testing.T) {
	testBind(t, &headerBinding{}, []testCase{
		{
			name: "empty header",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{},
				},
			},
			expectedBindings: model.Bindings{},
		},
		{
			name: "one header",
			given: model.Method{
				Params: []model.Param{
					{Name: "token", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
					},
				},
			},
			expectedBindings: model.Bindings{
				Headers: `http.Header{"Authorization":[]string{token}}`,
			},
			expectedBinded:  []string{"token"},
			expectedImports: []string{`"net/http"`},
		},
		{
			name: "multiple headers",
			given: model.Method{
				Params: []model.Param{
					{Name: "token", Type: "string"},
					{Name: "type", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
						{Key: "Content-Type", Values: []annotation.Variable{"{type}"}},
					},
				},
			},
			expectedBindings: model.Bindings{
				Headers: `http.Header{"Authorization":[]string{token}, "Content-Type":[]string{type}}`,
			},
			expectedBinded:  []string{"token", "type"},
			expectedImports: []string{`"net/http"`},
		},
		{
			name: "header with constant value",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Accept", Values: []annotation.Variable{"application/json"}},
					},
				},
			},
			expectedBindings: model.Bindings{
				Headers: `http.Header{"Accept":[]string{"application/json"}}`,
			},
			expectedImports: []string{`"net/http"`},
		},
		{
			name: "header param not found",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
					},
				},
			},
			expectedError: &ErrNotFound{Type: "header", Value: "token"},
		},
		{
			name: "multiple values for same header",
			given: model.Method{
				Params: []model.Param{
					{Name: "type1", Type: "string"},
					{Name: "type2", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Accept", Values: []annotation.Variable{"{type1}", "{type2}"}},
					},
				},
			},
			expectedBindings: model.Bindings{
				Headers: `http.Header{"Accept":[]string{type1, type2}}`,
			},
			expectedBinded:  []string{"type1", "type2"},
			expectedImports: []string{`"net/http"`},
		},
		{
			name: "mixed constant and variable values",
			given: model.Method{
				Params: []model.Param{
					{Name: "token", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Accept", Values: []annotation.Variable{"application/json"}},
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
					},
				},
			},
			expectedBindings: model.Bindings{
				Headers: `http.Header{"Accept":[]string{"application/json"}, "Authorization":[]string{token}}`,
			},
			expectedBinded:  []string{"token"},
			expectedImports: []string{`"net/http"`},
		},
	})
}
