package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/model"
	"github.com/lexcao/genapi/internal/parser/annotation"
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
			expectedBindings: model.Bindings{},
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
			expectedBindings: model.Bindings{
				PathParams: `map[string]string{"owner":owner}`,
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
			expectedBindings: model.Bindings{
				PathParams: `map[string]string{"owner":owner, "repo":repo}`,
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
			expectedBindings: model.Bindings{
				PathParams: `map[string]string{"owner":owner}`,
			},
			expectedBinded: []string{"owner"},
		},
	})
}
