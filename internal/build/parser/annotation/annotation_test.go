package annotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInterfaceAnnotations(t *testing.T) {
	tests := []struct {
		name     string
		comments []string
		want     InterfaceAnnotations
	}{
		{
			name:     "no annotations",
			comments: []string{},
			want:     InterfaceAnnotations{},
		},
		{
			name: "base url",
			comments: []string{
				`// @BaseURL("https://api.example.com")`,
			},
			want: InterfaceAnnotations{BaseURL: BaseURL{Value: "https://api.example.com"}},
		},
		{
			name: "base url with multiple annotations",
			comments: []string{
				`// @BaseURL("https://api.example.com1")`,
				`// @BaseURL("https://api.example.com2")`,
			},
			want: InterfaceAnnotations{BaseURL: BaseURL{Value: "https://api.example.com2"}},
		},
		{
			name: "multiple annotations",
			comments: []string{
				`// @BaseURL("https://api.example.com")`,
				`// @Header("Accept", "application/json")`,
			},
			want: InterfaceAnnotations{
				BaseURL: BaseURL{Value: "https://api.example.com"},
				Headers: []Header{
					{Key: "Accept", Values: []Variable{"application/json"}},
				},
			},
		},
		{
			name: "multiple headers",
			comments: []string{
				`// @Header("Accept", "application/json")`,
				`// @Header("Content-Type", "application/json")`,
			},
			want: InterfaceAnnotations{
				Headers: []Header{
					{Key: "Accept", Values: []Variable{"application/json"}},
					{Key: "Content-Type", Values: []Variable{"application/json"}},
				},
			},
		},
		{
			name: "headers with multiple values",
			comments: []string{
				`// @Header("Accept", "application/json", "application/xml")`,
			},
			want: InterfaceAnnotations{
				Headers: []Header{
					{Key: "Accept", Values: []Variable{"application/json", "application/xml"}},
				},
			},
		},
		{
			name: "unsupported annotation",
			comments: []string{
				`// @Unknown("GET", "POST")`,
			},
			want: InterfaceAnnotations{},
		},
		{
			name: "invalid annotation format",
			comments: []string{
				`// @BaseURL`,
			},
			want: InterfaceAnnotations{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ParseInterfaceAnnotations(test.comments)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestParseMethodAnnotations(t *testing.T) {
	tests := []struct {
		name     string
		comments []string
		want     MethodAnnotations
	}{
		{
			name:     "no annotations",
			comments: []string{},
			want:     MethodAnnotations{},
		},
		{
			name: "single query",
			comments: []string{
				`// @Query("page", "1")`,
			},
			want: MethodAnnotations{
				Queries: []Query{{Key: "page", Value: "1"}},
			},
		},
		{
			name: "multiple queries",
			comments: []string{
				`// @Query("page", "1")`,
				`// @Query("limit", "10")`,
			},
			want: MethodAnnotations{
				Queries: []Query{
					{Key: "page", Value: "1"},
					{Key: "limit", Value: "10"},
				},
			},
		},
		{
			name: "single header",
			comments: []string{
				`// @Header("Accept", "application/json")`,
			},
			want: MethodAnnotations{
				Headers: []Header{
					{Key: "Accept", Values: []Variable{"application/json"}},
				},
			},
		},
		{
			name: "multiple headers",
			comments: []string{
				`// @Header("Accept", "application/json")`,
				`// @Header("Content-Type", "application/json")`,
			},
			want: MethodAnnotations{
				Headers: []Header{
					{Key: "Accept", Values: []Variable{"application/json"}},
					{Key: "Content-Type", Values: []Variable{"application/json"}},
				},
			},
		},
		{
			name: "header with multiple values",
			comments: []string{
				`// @Header("Accept", "application/json", "application/xml")`,
			},
			want: MethodAnnotations{
				Headers: []Header{
					{Key: "Accept", Values: []Variable{"application/json", "application/xml"}},
				},
			},
		},
		{
			name: "mixed queries and headers",
			comments: []string{
				`// @Query("page", "1")`,
				`// @Header("Accept", "application/json")`,
				`// @Query("limit", "10")`,
				`// @Header("Content-Type", "application/json")`,
			},
			want: MethodAnnotations{
				Queries: []Query{
					{Key: "page", Value: "1"},
					{Key: "limit", Value: "10"},
				},
				Headers: []Header{
					{Key: "Accept", Values: []Variable{"application/json"}},
					{Key: "Content-Type", Values: []Variable{"application/json"}},
				},
			},
		},
		{
			name: "invalid query format",
			comments: []string{
				`// @Query("page")`,
			},
			want: MethodAnnotations{},
		},
		{
			name: "unsupported annotation",
			comments: []string{
				`// @Unknown("GET", "POST")`,
			},
			want: MethodAnnotations{},
		},
		{
			name: "alias",
			comments: []string{
				`// @Get("https://api.example.com")`,
			},
			want: MethodAnnotations{
				RequestLine: RequestLine{
					Method: "Get",
					Path:   "https://api.example.com",
				},
			},
		},
		{
			name: "alias with no alias",
			comments: []string{
				`// @Get("https://api.example.com")`,
				`// @RequestLine("Post", "https://api.example.com")`,
			},
			want: MethodAnnotations{
				RequestLine: RequestLine{
					Method: "Post",
					Path:   "https://api.example.com",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ParseMethodAnnotations(test.comments)
			assert.Equal(t, test.want, got)
		})
	}
}
