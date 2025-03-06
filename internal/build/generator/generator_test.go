package generator

import (
	"strings"
	"testing"

	"github.com/lexcao/genapi/internal/build/common"
	"github.com/lexcao/genapi/internal/build/model"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFile(t *testing.T) {
	actual, err := GenerateFile("test.go", []model.Interface{
		{
			Name:    "GitHub",
			Package: "generator",
			Imports: common.SetOf(
				`"context"`,
				`"net/http"`,
				`"github.com/lexcao/genapi"`,
			),
			Bindings: &model.InterfaceBindings{
				Config: `genapi.Config{BaseURL: "https://api.github.com", Headers: http.Header{"Accept": []string{"application/vnd.github.v3+json"}}}`,
			},
			Methods: []model.Method{
				{
					Name:      "ListRepositories",
					Interface: &model.Interface{Name: "GitHub"},
					Params: []model.Param{
						{Name: "ctx", Type: "context.Context"},
						{Name: "owner", Type: "string"},
					},
					Results: []model.Param{
						{Type: "error"},
					},
					Bindings: &model.MethodBindings{
						Results: &model.BindingResults{
							Assignment: "resp, err",
							Statement:  "genapi.HandleResponse0(resp, err)",
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

	expect := `// CODE GENERATED BY genapi. DO NOT EDIT.
package generator

import (
	"context"
	"github.com/lexcao/genapi"
	"net/http"
)

type implGitHub struct {
	client genapi.HttpClient
}

// SetHttpClient implments genapi.Interface
func (i *implGitHub) SetHttpClient(client genapi.HttpClient) {
	i.client = client
}

func (i *implGitHub) ListRepositories(ctx context.Context, owner string) error {
	resp, err := i.client.Do(&genapi.Request{
		Method:     "GET",
		Path:       "/users/{owner}/repos",
		PathParams: map[string]string{"owner": owner},
	})
	return genapi.HandleResponse0(resp, err)
}

func init() {
	genapi.Register[GitHub, *implGitHub](
		genapi.Config{BaseURL: "https://api.github.com", Headers: http.Header{"Accept": []string{"application/vnd.github.v3+json"}}},
	)
}
`
	assert.Equal(t, expect, string(actual))
}

func TestGenerateFileWithMultipleInterfaces(t *testing.T) {
	actual, err := GenerateFile("test.go", []model.Interface{
		{
			Name:    "GitHub",
			Package: "generator",
			Imports: common.SetOf(
				`"context"`,
				`"github.com/lexcao/genapi"`,
				`"net/http"`,
			),
			Bindings: &model.InterfaceBindings{
				Config: `genapi.Config{BaseURL: "https://api.github.com"}`,
			},
			Methods: []model.Method{
				{
					Name:      "ListRepositories",
					Interface: &model.Interface{Name: "GitHub"},
					Params: []model.Param{
						{Name: "ctx", Type: "context.Context"},
					},
					Results: []model.Param{
						{Type: "error"},
					},
					Bindings: &model.MethodBindings{
						Results: &model.BindingResults{
							Assignment: "resp, err",
							Statement:  "genapi.HandleResponse0(resp, err)",
						},
						Path:   "/repos",
						Method: "GET",
					},
				},
			},
		},
		{
			Name:    "Weather",
			Package: "generator",
			Imports: common.SetOf(
				`"context"`,
				`"github.com/lexcao/genapi"`,
				`"net/url"`,
			),
			Bindings: &model.InterfaceBindings{
				Config: `genapi.Config{BaseURL: "https://api.weather.com"}`,
			},
		},
	})
	require.NoError(t, err)

	expect := `// CODE GENERATED BY genapi. DO NOT EDIT.
package generator

import (
	"context"
	"github.com/lexcao/genapi"
	"net/http"
	"net/url"
)

type implGitHub struct {
	client genapi.HttpClient
}

// SetHttpClient implments genapi.Interface
func (i *implGitHub) SetHttpClient(client genapi.HttpClient) {
	i.client = client
}

func (i *implGitHub) ListRepositories(ctx context.Context) error {
	resp, err := i.client.Do(&genapi.Request{
		Method: "GET",
		Path:   "/repos",
	})
	return genapi.HandleResponse0(resp, err)
}

type implWeather struct {
	client genapi.HttpClient
}

// SetHttpClient implments genapi.Interface
func (i *implWeather) SetHttpClient(client genapi.HttpClient) {
	i.client = client
}

func init() {
	genapi.Register[GitHub, *implGitHub](
		genapi.Config{BaseURL: "https://api.github.com"},
	)
	genapi.Register[Weather, *implWeather](
		genapi.Config{BaseURL: "https://api.weather.com"},
	)
}
`
	assert.Equal(t, expect, string(actual))
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

// SetHttpClient implments genapi.Interface
func (i *implTest) SetHttpClient(client genapi.HttpClient) {
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
				Interface: &model.Interface{Name: "Client"},
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
				Interface: &model.Interface{Name: "Client"},
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
		Path: "/repos/{owner}",
		PathParams: map[string]string{
			"owner": owner,
		},
	})
}
`,
		},
		{
			name: "no params one result",
			method: model.Method{
				Name:      "OneResult",
				Interface: &model.Interface{Name: "Client"},
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
				Interface: &model.Interface{Name: "Client"},
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
		Path: "/repos/{owner}/{repo}",
		PathParams: map[string]string{
			"owner": owner,
			"repo":  repo,
		},
	})
	return genapi.HandleResponse[Result](resp, err)
}
`,
		},
		{
			name: "one param one request with imports",
			method: model.Method{
				Name:      "WithImports",
				Interface: &model.Interface{Name: "Client"},
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
				Interface: &model.Interface{Name: "Client"},
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
	Method: "GET",
	Path:   "/repos/{owner}",
	PathParams: map[string]string{
		"owner": owner,
	},
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
	Queries: url.Values{
		"sort": []string{
			sort,
		},
	},
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
						{Key: "page", Value: "{page}"},
						{Key: "sort", Value: "desc"},
					},
				},
			},
			expected: `
i.client.Do(&genapi.Request{
	Queries: url.Values{
		"page": []string{
			strconv.Itoa(int(page)),
		},
		"sort": []string{
			"desc",
		},
	},
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
	Header: http.Header{
		"Authorization": []string{
			token,
		},
	},
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
	Header: http.Header{
		"Authorization": []string{
			token,
		},
		"Content-Type": []string{
			"application/json",
		},
	},
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
			test.method.Interface = &model.Interface{}
			actual, err := generateMethod(tmplMethodBody, test.method)
			require.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(test.expected), strings.TrimSpace(actual))
		})
	}
}
