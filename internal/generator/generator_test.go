package generator

import (
	"os"
	"strings"
	"testing"

	"github.com/lexcao/genapi/internal/model"
	"github.com/lexcao/genapi/internal/parser/annotation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFile(t *testing.T) {
	err := GenerateFile("test.go", []model.Interface{
		{
			Name:    "GitHub",
			Package: "generator",
			Methods: []model.Method{
				{
					Name: "ListRepositories",
					Params: []model.Param{
						{Name: "ctx", Type: "context.Context"},
						{Name: "owner", Type: "string"},
					},
					Results: []model.Param{
						{Type: "[]Repository"},
						{Type: "error"},
					},
					Bindings: &model.Bindings{
						Results: &model.BindingResults{
							Assignment: "resp, err",
							Statement:  "genapi.HandleResponse[[]Repository](resp, err)",
						},
						PathParams: "map[string]string{\"owner\": owner}",
						Path:       "/users/{owner}/repos",
						Method:     "GET",
					},
				},
			},
		},
	})
	require.NoError(t, err)

	// check file exists
	_, err = os.Stat("test.gen.go")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.Remove("test.gen.go")
	})
}

func TestGenerateInterface(t *testing.T) {
	actual, err := generateInterface(tmplInterface, model.Interface{
		Name: "Test",
	})
	require.NoError(t, err)
	expected := `
type implTest struct {
	client genapi.HttpClient
}

// setHttpClient implments genapi.Interface
func (i *implTest) setHttpClient(client genapi.HttpClient) {
	i.client = client
}
`
	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(actual))
}

func TestGenerateMethod(t *testing.T) {
	tests := []struct {
		name     string
		method   model.Method
		expected string
	}{
		{
			name: "no params no results",
			method: model.Method{
				Name:      "NoParams",
				Interface: "Client",
			},
			expected: `
func (i *implClient) NoParams() {
	i.client.Do(&genapi.Request{})
}
`,
		},
		{
			name: "one param no results",
			method: model.Method{
				Name:      "OneParam",
				Interface: "Client",
				Params: []model.Param{
					{Name: "owner", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/repos/{owner}",
					},
				},
			},
			expected: `
func (i *implClient) OneParam(owner string) {
	i.client.Do(&genapi.Request{
		Path:       "/repos/{owner}",
		PathParams: map[string]string{"owner": owner},
	})
}
`,
		},
		{
			name: "no params one result",
			method: model.Method{
				Name:      "OneResult",
				Interface: "Client",
				Results: []model.Param{
					{Type: "error"},
				},
			},
			expected: `
func (i *implClient) OneResult() error {
	resp, err := i.client.Do(&genapi.Request{})
	return genapi.HandleResponse0(resp, err)
}
`,
		},
		{
			name: "two params two results",
			method: model.Method{
				Name:      "TwoParamsTwoResults",
				Interface: "Client",
				Params: []model.Param{
					{Name: "owner", Type: "string"},
					{Name: "repo", Type: "string"},
				},
				Results: []model.Param{
					{Type: "Result"},
					{Type: "error"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/repos/{owner}/{repo}",
					},
				},
			},
			expected: `
func (i *implClient) TwoParamsTwoResults(owner string, repo string) (Result, error) {
	resp, err := i.client.Do(&genapi.Request{
		Path:       "/repos/{owner}/{repo}",
		PathParams: map[string]string{"owner": owner, "repo": repo},
	})
	return genapi.HandleResponse[Result](resp, err)
}
`,
		},
		{
			name: "one param one request with imports",
			method: model.Method{
				Name:      "WithImports",
				Interface: "Client",
				Params: []model.Param{
					{Name: "ctx", Type: "context.Context"},
				},
				Results: []model.Param{
					{Type: "result.Result"},
				},
			},
			expected: `
func (i *implClient) WithImports(ctx context.Context) result.Result {
	resp, err := i.client.Do(&genapi.Request{
		Context: ctx,
	})
	return genapi.MustHandleResponse[result.Result](resp, err)
}
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := generateMethod(tmplMethod, test.method)
			require.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(test.expected), strings.TrimSpace(actual))
		})
	}
}

func TestGenerateMethodBody(t *testing.T) {
	tests := []struct {
		name     string
		method   model.Method
		expected string
	}{
		{
			name: "no body",
			method: model.Method{
				Name:      "NoBody",
				Interface: "Client",
			},
			expected: `
i.client.Do(&genapi.Request{})
`,
		},
		{
			name: "with request line",
			method: model.Method{
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Method: "GET",
						Path:   "/repos",
					},
				},
			},
			expected: `
i.client.Do(&genapi.Request{
	Method: "GET",
	Path:   "/repos",
})
`,
		},
		{
			name: "with request line with path params",
			method: model.Method{
				Params: []model.Param{
					{Name: "owner", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Method: "GET",
						Path:   "/repos/{owner}",
					},
				},
			},
			expected: `
i.client.Do(&genapi.Request{
	Method:     "GET",
	Path:       "/repos/{owner}",
	PathParams: map[string]string{"owner": owner},
})
`,
		},
		{
			name: "with query single",
			method: model.Method{
				Params: []model.Param{
					{Name: "sort", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "sort", Value: "{sort}"},
					},
				},
			},
			expected: `
i.client.Do(&genapi.Request{
	Queries: url.Values{"sort": []string{sort}},
})
`,
		},
		{
			name: "with query multiple",
			method: model.Method{
				Params: []model.Param{
					{Name: "page", Type: "int"},
				},
				Annotations: annotation.MethodAnnotations{
					Queries: []annotation.Query{
						{Key: "sort", Value: "desc"},
						{Key: "page", Value: "{page}"},
					},
				},
			},
			expected: `
i.client.Do(&genapi.Request{
	Queries: url.Values{"page": []string{page}, "sort": []string{"desc"}},
})
`,
		},
		{
			name: "with header single",
			method: model.Method{
				Params: []model.Param{
					{Name: "token", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
					},
				},
			},
			expected: `
i.client.Do(&genapi.Request{
	Headers: http.Header{"Authorization": []string{token}},
})
`,
		},
		{
			name: "with header multiple",
			method: model.Method{
				Params: []model.Param{
					{Name: "token", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
						{Key: "Content-Type", Values: []annotation.Variable{"application/json"}},
					},
				},
			},
			expected: `
i.client.Do(&genapi.Request{
	Headers: http.Header{"Authorization": []string{token}, "Content-Type": []string{"application/json"}},
})
`,
		},
		{
			name: "with context",
			method: model.Method{
				Params: []model.Param{
					{Name: "ctx", Type: "context.Context"},
				},
			},
			expected: `
i.client.Do(&genapi.Request{
	Context: ctx,
})
`,
		},
		{
			name: "with body",
			method: model.Method{
				Params: []model.Param{
					{Name: "body", Type: "[]byte"},
				},
			},
			expected: `
i.client.Do(&genapi.Request{
	Body: body,
})
`,
		},
		{
			name: "with results 1 any",
			method: model.Method{
				Results: []model.Param{
					{Type: "Result"},
				},
			},
			expected: `
resp, err := i.client.Do(&genapi.Request{})
return genapi.MustHandleResponse[Result](resp, err)
			`,
		},
		{
			name: "with results 1 genapi.Response",
			method: model.Method{
				Results: []model.Param{
					{Type: "genapi.Response"},
				},
			},
			expected: `
resp, _ := i.client.Do(&genapi.Request{})
return *resp
			`,
		},
		{
			name: "with results 1 *genapi.Response",
			method: model.Method{
				Results: []model.Param{
					{Type: "*genapi.Response"},
				},
			},
			expected: `
resp, _ := i.client.Do(&genapi.Request{})
return resp
			`,
		},
		{
			name: "with results 1 error",
			method: model.Method{
				Results: []model.Param{
					{Type: "error"},
				},
			},
			expected: `
resp, err := i.client.Do(&genapi.Request{})
return genapi.HandleResponse0(resp, err)
			`,
		},
		{
			name: "with results 2 any",
			method: model.Method{
				Results: []model.Param{
					{Type: "Result"},
					{Type: "error"},
				},
			},
			expected: `
resp, err := i.client.Do(&genapi.Request{})
return genapi.HandleResponse[Result](resp, err)
			`,
		},
		{
			name: "with results 2 genapi.Response",
			method: model.Method{
				Results: []model.Param{
					{Type: "genapi.Response"},
					{Type: "error"},
				},
			},
			expected: `
resp, err := i.client.Do(&genapi.Request{})
return *resp, err
			`,
		},
		{
			name: "with results 2 *genapi.Response",
			method: model.Method{
				Results: []model.Param{
					{Type: "*genapi.Response"},
					{Type: "error"},
				},
			},
			expected: `
resp, err := i.client.Do(&genapi.Request{})
return resp, err
			`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := generateMethod(tmplMethodBody, test.method)
			require.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(test.expected), strings.TrimSpace(actual))
		})
	}
}
