package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBind(t *testing.T) {
	cases := []struct {
		name             string
		given            model.Method
		expectedBindings model.Bindings
		expectedImports  []string
	}{
		{
			name:             "empty",
			given:            model.Method{},
			expectedBindings: model.Bindings{},
		},
		{
			name: "all bindings",
			given: model.Method{
				Name: "GetRepo",
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Method: "GET",
						Path:   "/repos/{owner}/{repo}",
					},
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
						{Key: "Content-Type", Values: []annotation.Variable{"application/json"}},
					},
					Queries: []annotation.Query{
						{Key: "sort", Value: "{sort}"},
						{Key: "page", Value: "1"},
						{Key: "perPage", Value: "10"},
					},
				},
				Params: []model.Param{
					{Name: "ctx", Type: "context.Context"},
					{Name: "owner", Type: "string"},
					{Name: "repo", Type: "string"},
					{Name: "request", Type: "Request"},
					{Name: "sort", Type: "string"},
					{Name: "token", Type: "string"},
				},
			},
			expectedBindings: model.Bindings{
				Method:     "GET",
				Path:       "/repos/{owner}/{repo}",
				Context:    "ctx",
				Body:       "request",
				PathParams: `map[string]string{"owner":owner, "repo":repo}`,
				Headers:    `http.Header{"Authorization":[]string{token}, "Content-Type":[]string{"application/json"}}`,
				Queries:    `url.Values{"page":[]string{"1"}, "perPage":[]string{"10"}, "sort":[]string{sort}}`,
			},
			expectedImports: []string{`"net/http"`, `"net/url"`},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.given.Interface = &model.Interface{}
			err := BindMethod(&c.given)
			require.NoError(t, err)
			assert.Equal(t, c.expectedBindings, *c.given.Bindings)
			assert.ElementsMatch(t, c.expectedImports, c.given.Interface.Imports.Slices())
		})
	}
}
