package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
)

func TestBindPathParam(t *testing.T) {
	testBind(t, &pathParamBinding{}, []testCase{
		{
			name: "no path params",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/users/me",
					},
				},
			},
			expectedBindings: model.MethodBindings{},
		},
		{
			name: "one path param",
			given: model.Method{
				Params: []model.Param{
					{Name: "owner", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/repos/{owner}",
					},
				},
			},
			expectedBindings: model.MethodBindings{
				PathParams: `map[string]string{
	"owner": owner,
}`,
			},
			expectedBinded: []string{"owner"},
		},
		{
			name: "multiple path params",
			given: model.Method{
				Params: []model.Param{
					{Name: "owner", Type: "string"},
					{Name: "repo", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/repos/{owner}/{repo}",
					},
				},
			},
			expectedBindings: model.MethodBindings{
				PathParams: `map[string]string{
	"owner": owner,
	"repo": repo,
}`,
			},
			expectedBinded: []string{"owner", "repo"},
		},
		{
			name: "path param not found",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/repos/{owner}",
					},
				},
			},
			expectedError: &ErrNotFound{Type: "path", Value: "owner"},
		},
		{
			name: "mixed constant and variable path",
			given: model.Method{
				Params: []model.Param{
					{Name: "owner", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/repos/{owner}/commits",
					},
				},
			},
			expectedBindings: model.MethodBindings{
				PathParams: `map[string]string{
	"owner": owner,
}`,
			},
			expectedBinded: []string{"owner"},
		},
		{
			name: "integer path param",
			given: model.Method{
				Params: []model.Param{
					{Name: "id", Type: "int"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/repos/{id}",
					},
				},
			},
			expectedBindings: model.MethodBindings{
				PathParams: `map[string]string{
	"id": strconv.Itoa(int(id)),
}`,
			},
			expectedBinded:  []string{"id"},
			expectedImports: []string{`"strconv"`},
		},
		{
			name: "float path param",
			given: model.Method{
				Params: []model.Param{
					{Name: "rating", Type: "float64"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/movies/{rating}",
					},
				},
			},
			expectedBindings: model.MethodBindings{
				PathParams: `map[string]string{
	"rating": strconv.FormatFloat(float64(rating), 'f', -1, 64),
}`,
			},
			expectedBinded:  []string{"rating"},
			expectedImports: []string{`"strconv"`},
		},
		{
			name: "boolean path param",
			given: model.Method{
				Params: []model.Param{
					{Name: "active", Type: "bool"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/users/{active}",
					},
				},
			},
			expectedBindings: model.MethodBindings{
				PathParams: `map[string]string{
	"active": strconv.FormatBool(active),
}`,
			},
			expectedBinded:  []string{"active"},
			expectedImports: []string{`"strconv"`},
		},
		{
			name: "mixed type path params",
			given: model.Method{
				Params: []model.Param{
					{Name: "owner", Type: "string"},
					{Name: "id", Type: "int"},
					{Name: "active", Type: "bool"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/repos/{owner}/{id}/{active}",
					},
				},
			},
			expectedBindings: model.MethodBindings{
				PathParams: `map[string]string{
	"owner": owner,
	"id": strconv.Itoa(int(id)),
	"active": strconv.FormatBool(active),
}`,
			},
			expectedBinded:  []string{"owner", "id", "active"},
			expectedImports: []string{`"strconv"`},
		},
	})
}
