package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
)

func TestBindQuery(t *testing.T) {
	testBind(t, &queryBinding{}, []testCase{
		{
			name: "empty query",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{},
				},
			},
			expectedBindings: model.Bindings{},
		},
		{
			name: "one query",
			given: model.Method{
				Params: []model.Param{
					{Name: "sort", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "sort", Value: "{sort}"},
					},
				},
			},
			expectedBindings: model.Bindings{
				Queries: `url.Values{"sort":[]string{sort}}`,
			},
			expectedBinded:  []string{"sort"},
			expectedImports: []string{`"net/url"`},
		},
		{
			name: "multiple queries",
			given: model.Method{
				Params: []model.Param{
					{Name: "sort", Type: "string"},
					{Name: "page", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "sort", Value: "{sort}"},
						{Key: "page", Value: "{page}"},
					},
				},
			},
			expectedBindings: model.Bindings{
				Queries: `url.Values{"page":[]string{page}, "sort":[]string{sort}}`,
			},
			expectedBinded:  []string{"sort", "page"},
			expectedImports: []string{`"net/url"`},
		},
		{
			name: "query with constant value",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "sort", Value: "desc"},
					},
				},
			},
			expectedBindings: model.Bindings{
				Queries: `url.Values{"sort":[]string{"desc"}}`,
			},
			expectedImports: []string{`"net/url"`},
		},
		{
			name: "query param not found",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "sort", Value: "{sort}"},
					},
				},
			},
			expectedError: &ErrNotFound{Type: "query", Value: "sort"},
		},
		{
			name: "multiple values for same key",
			given: model.Method{
				Params: []model.Param{
					{Name: "sort1", Type: "string"},
					{Name: "sort2", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "sort", Value: "{sort1}"},
						{Key: "sort", Value: "{sort2}"},
					},
				},
			},
			expectedBindings: model.Bindings{
				Queries: `url.Values{"sort":[]string{sort1, sort2}}`,
			},
			expectedBinded:  []string{"sort1", "sort2"},
			expectedImports: []string{`"net/url"`},
		},
	})
}
