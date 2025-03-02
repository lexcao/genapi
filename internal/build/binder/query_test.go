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
			expectedBindings: model.MethodBindings{},
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
			expectedBindings: model.MethodBindings{
				Queries: `url.Values{
	"sort": []string{
		sort,
	},
}`,
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
			expectedBindings: model.MethodBindings{
				Queries: `url.Values{
	"sort": []string{
		sort,
	},
	"page": []string{
		page,
	},
}`,
			},
			expectedBinded:  []string{"sort", "page"},
			expectedImports: []string{`"net/url"`},
		},
		{
			name: "query with constant value",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "sort", Value: "asc"},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Queries: `url.Values{
	"sort": []string{
		"asc",
	},
}`,
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
			name: "multiple values for same query",
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
			expectedBindings: model.MethodBindings{
				Queries: `url.Values{
	"sort": []string{
		sort1,
		sort2,
	},
}`,
			},
			expectedBinded:  []string{"sort1", "sort2"},
			expectedImports: []string{`"net/url"`},
		},
		{
			name: "integer query param",
			given: model.Method{
				Params: []model.Param{
					{Name: "page", Type: "int"},
				},
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "page", Value: "{page}"},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Queries: `url.Values{
	"page": []string{
		strconv.Itoa(int(page)),
	},
}`,
			},
			expectedBinded:  []string{"page"},
			expectedImports: []string{`"net/url"`, `"strconv"`},
		},
		{
			name: "float query param",
			given: model.Method{
				Params: []model.Param{
					{Name: "price", Type: "float64"},
				},
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "price", Value: "{price}"},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Queries: `url.Values{
	"price": []string{
		strconv.FormatFloat(float64(price), 'f', -1, 64),
	},
}`,
			},
			expectedBinded:  []string{"price"},
			expectedImports: []string{`"net/url"`, `"strconv"`},
		},
		{
			name: "boolean query param",
			given: model.Method{
				Params: []model.Param{
					{Name: "active", Type: "bool"},
				},
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "active", Value: "{active}"},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Queries: `url.Values{
	"active": []string{
		strconv.FormatBool(active),
	},
}`,
			},
			expectedBinded:  []string{"active"},
			expectedImports: []string{`"net/url"`, `"strconv"`},
		},
		{
			name: "mixed type query params",
			given: model.Method{
				Params: []model.Param{
					{Name: "name", Type: "string"},
					{Name: "page", Type: "int"},
					{Name: "active", Type: "bool"},
				},
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "name", Value: "{name}"},
						{Key: "page", Value: "{page}"},
						{Key: "active", Value: "{active}"},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Queries: `url.Values{
	"name": []string{
		name,
	},
	"page": []string{
		strconv.Itoa(int(page)),
	},
	"active": []string{
		strconv.FormatBool(active),
	},
}`,
			},
			expectedBinded:  []string{"name", "page", "active"},
			expectedImports: []string{`"net/url"`, `"strconv"`},
		},
	})
}
