package generator

import (
	"testing"

	"github.com/lexcao/genapi/internal/model"
	"github.com/lexcao/genapi/internal/parser/annotation"
	"github.com/stretchr/testify/assert"
)

func TestParseQueries(t *testing.T) {
	tests := []struct {
		name    string
		queries []annotation.Query
		expect  string
	}{
		{
			name:    "empty query",
			queries: []annotation.Query{},
			expect:  "url.Values{}",
		},
		{
			name: "one query",
			queries: []annotation.Query{
				{Key: "sort", Value: "desc"},
			},
			expect: `url.Values{"sort":[]string{"desc"}}`,
		},
		{
			name: "one query with variable",
			queries: []annotation.Query{
				{Key: "sort", Value: "{sort}"},
			},
			expect: `url.Values{"sort":[]string{sort}}`,
		},
		{
			name: "multiple queries",
			queries: []annotation.Query{
				{Key: "sort", Value: "desc"},
				{Key: "page", Value: "1"},
			},
			expect: `url.Values{"page":[]string{"1"}, "sort":[]string{"desc"}}`,
		},
		{
			name: "multiple queries with variable",
			queries: []annotation.Query{
				{Key: "sort", Value: "desc"},
				{Key: "page", Value: "{page}"},
			},
			expect: `url.Values{"page":[]string{page}, "sort":[]string{"desc"}}`,
		},
		{
			name: "one query with multiple values",
			queries: []annotation.Query{
				{Key: "sort", Value: "desc"},
				{Key: "sort", Value: "asc"},
			},
			expect: `url.Values{"sort":[]string{"desc", "asc"}}`,
		},
		{
			name: "multiple queries with multiple values",
			queries: []annotation.Query{
				{Key: "sort", Value: "desc"},
				{Key: "sort", Value: "asc"},
				{Key: "page", Value: "1"},
				{Key: "page", Value: "2"},
			},
			expect: `url.Values{"page":[]string{"1", "2"}, "sort":[]string{"desc", "asc"}}`,
		},
		{
			name: "multiple queries with multiple variables",
			queries: []annotation.Query{
				{Key: "sort", Value: "{sort1}"},
				{Key: "sort", Value: "{sort2}"},
				{Key: "page", Value: "{page1}"},
				{Key: "page", Value: "{page2}"},
			},
			expect: `url.Values{"page":[]string{page1, page2}, "sort":[]string{sort1, sort2}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := parseQueries(test.queries)
			assert.Equal(t, test.expect, actual)
		})
	}
}

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers []annotation.Header
		expect  string
	}{
		{
			name:    "empty header",
			headers: []annotation.Header{},
			expect:  "http.Header{}",
		},
		{
			name: "one header",
			headers: []annotation.Header{
				{Key: "Accept", Values: []annotation.Variable{"application/json"}},
			},
			expect: `http.Header{"Accept":[]string{"application/json"}}`,
		},
		{
			name: "one header with variable",
			headers: []annotation.Header{
				{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
			},
			expect: `http.Header{"Authorization":[]string{token}}`,
		},
		{
			name: "multiple headers",
			headers: []annotation.Header{
				{Key: "Accept", Values: []annotation.Variable{"application/json"}},
				{Key: "Content-Type", Values: []annotation.Variable{"application/json"}},
			},
			expect: `http.Header{"Accept":[]string{"application/json"}, "Content-Type":[]string{"application/json"}}`,
		},
		{
			name: "multiple headers with variable",
			headers: []annotation.Header{
				{Key: "Accept", Values: []annotation.Variable{"application/json"}},
				{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
			},
			expect: `http.Header{"Accept":[]string{"application/json"}, "Authorization":[]string{token}}`,
		},
		{
			name: "one header with multiple values",
			headers: []annotation.Header{
				{Key: "Accept", Values: []annotation.Variable{"application/json", "application/xml"}},
			},
			expect: `http.Header{"Accept":[]string{"application/json", "application/xml"}}`,
		},
		{
			name: "multiple headers with multiple values",
			headers: []annotation.Header{
				{Key: "Accept", Values: []annotation.Variable{"application/json", "application/xml"}},
				{Key: "Content-Type", Values: []annotation.Variable{"application/json", "text/plain"}},
			},
			expect: `http.Header{"Accept":[]string{"application/json", "application/xml"}, "Content-Type":[]string{"application/json", "text/plain"}}`,
		},
		{
			name: "multiple headers with multiple variables",
			headers: []annotation.Header{
				{Key: "Accept", Values: []annotation.Variable{"{type1}", "{type2}"}},
				{Key: "Authorization", Values: []annotation.Variable{"{token1}", "{token2}"}},
			},
			expect: `http.Header{"Accept":[]string{type1, type2}, "Authorization":[]string{token1, token2}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := parseHeaders(test.headers)
			assert.Equal(t, test.expect, actual)
		})
	}
}

func TestParsePathParams(t *testing.T) {
	tests := []struct {
		name   string
		method model.Method
		expect string
	}{
		{
			name: "no path params",
			method: model.Method{Annotations: annotation.MethodAnnotations{
				RequestLine: annotation.RequestLine{
					Path: "/users/me",
				},
			}},
			expect: "map[string]string{}",
		},
		{
			name: "one path param",
			method: model.Method{
				Params: []model.Param{{Name: "owner", Type: "string"}},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/repos/{owner}",
					},
				},
			},
			expect: `map[string]string{"owner":owner}`,
		},
		{
			name: "multiple path params",
			method: model.Method{
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
			expect: `map[string]string{"owner":owner, "repo":repo}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := parsePathParams(test.method)
			assert.Equal(t, test.expect, actual)
		})
	}
}
