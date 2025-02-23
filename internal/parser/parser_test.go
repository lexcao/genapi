package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/lexcao/genapi/internal/model"
	"github.com/lexcao/genapi/internal/parser/annotation"
	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	interfaces, err := ParseFile("testdata/example.go")
	require.NoError(t, err)
	require.Len(t, interfaces, 1)

	actual := interfaces[0]
	expect := model.Interface{
		Name:    "GitHub",
		Package: "testdata",
		Imports: []string{
			`"context"`,
			`"github.com/lexcao/genapi"`,
		},
		Annotations: annotation.InterfaceAnnotations{
			BaseURL: annotation.BaseURL{Value: "https://api.github.com"},
			Headers: []annotation.Header{
				{Key: "Accept", Values: []annotation.Variable{"application/vnd.github.v3+json"}},
			},
		},
		Methods: []model.Method{
			{
				Name:      "ListRepositories",
				Interface: "GitHub",
				Params: []model.Param{
					{Name: "ctx", Type: "context.Context"},
					{Name: "owner", Type: "string"},
					{Name: "sort", Type: "string"},
				},
				Results: []model.Param{
					{Type: "[]Repository"},
					{Type: "error"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Method: "GET",
						Path:   "/repos/{owner}",
					},
					Queries: []annotation.Query{
						{Key: "sort", Value: "{sort}"},
					},
				},
			},
			{
				Name:      "ListContributors",
				Interface: "GitHub",
				Params: []model.Param{
					{Name: "ctx", Type: "context.Context"},
					{Name: "owner", Type: "string"},
					{Name: "repo", Type: "string"},
				},
				Results: []model.Param{
					{Type: "[]Contributor"},
					{Type: "error"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Method: "GET",
						Path:   "/repos/{owner}/{repo}/contributors",
					},
				},
			},
			{
				Name:      "CreateIssue",
				Interface: "GitHub",
				Params: []model.Param{
					{Name: "ctx", Type: "context.Context"},
					{Name: "issue", Type: "Issue"},
					{Name: "owner", Type: "string"},
					{Name: "repo", Type: "string"},
				},
				Results: []model.Param{
					{Type: "error"},
				},
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Method: "POST",
						Path:   "/repos/{owner}/{repo}/issues",
					},
				},
			},
		},
	}
	require.Equal(t, expect, actual)
}

func TestParseInterface(t *testing.T) {
	tests := []struct {
		name string
		code string
		want model.Interface
	}{
		{
			name: "no methods",
			code: `
			type Interface interface {}
			`,
			want: model.Interface{Name: "Interface", Package: "main"},
		},
		{
			name: "one method",
			code: `
			type Interface interface {
				Method()
			}
			`,
			want: model.Interface{Name: "Interface", Package: "main", Methods: []model.Method{{Name: "Method", Interface: "Interface"}}},
		},
		{
			name: "many methods",
			code: `
			// @BaseURL("https://api.example.com")
			type Interface interface {
				Test(a int)
				Test2(a int) error
			}
			`,
			want: model.Interface{
				Name:    "Interface",
				Package: "main",
				Methods: []model.Method{
					{Name: "Test", Interface: "Interface", Params: []model.Param{{Name: "a", Type: "int"}}},
					{Name: "Test2", Interface: "Interface", Params: []model.Param{{Name: "a", Type: "int"}}, Results: []model.Param{{Type: "error"}}},
				},
				Annotations: annotation.InterfaceAnnotations{
					BaseURL: annotation.BaseURL{Value: "https://api.example.com"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parseCodeNode(t, test.code, func(file *ast.File) func(ast.Node) bool {
				return func(n ast.Node) bool {
					if decl, ok := n.(*ast.GenDecl); ok {
						for _, spec := range decl.Specs {
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								if iface, ok := typeSpec.Type.(*ast.InterfaceType); ok {
									got := parseInterface(parseInterfaceParams{
										File:        file,
										Declaration: decl,
										TypeSpec:    typeSpec,
										Interface:   iface,
									})
									require.Equal(t, test.want, got)
									return false
								}
							}
						}
					}
					return true
				}
			})
		})
	}
}

func TestParseMethod(t *testing.T) {
	tests := []struct {
		name string
		code string
		want model.Method
	}{
		{
			name: "no params",
			code: `
			func test() {}
			`,
			want: model.Method{
				Name: "test",
			},
		},
		{
			name: "one param",
			code: `
			func test(a int) {}
			`,
			want: model.Method{
				Name: "test",
				Params: []model.Param{
					{Name: "a", Type: "int"},
				},
			},
		},
		{
			name: "one returns",
			code: `
			func test() error {}
			`,
			want: model.Method{
				Name: "test",
				Results: []model.Param{
					{Name: "", Type: "error"},
				},
			},
		},
		{
			name: "many params and returns",
			code: `
			func test(a int, b string) (bool, error) {}
			`,
			want: model.Method{
				Name: "test",
				Params: []model.Param{
					{Name: "a", Type: "int"},
					{Name: "b", Type: "string"},
				},
				Results: []model.Param{
					{Name: "", Type: "bool"},
					{Name: "", Type: "error"},
				},
			},
		},
		{
			name: "with imports",
			code: `
			func test(ctx context.Context) (context.Context, error) {}
			`,
			want: model.Method{
				Name: "test",
				Params: []model.Param{
					{Name: "ctx", Type: "context.Context"},
				},
				Results: []model.Param{
					{Name: "", Type: "error"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parseCodeNode(t, test.code, func(file *ast.File) func(ast.Node) bool {
				return func(n ast.Node) bool {
					if method, ok := n.(*ast.Field); ok {
						if fn, ok := method.Type.(*ast.FuncType); ok {
							got := parseMethod(method, fn)
							require.Equal(t, test.want, got)
							return false
						}
					}
					return true
				}
			})
		})
	}
}

func TestParseParams(t *testing.T) {
	tests := []struct {
		name string
		code string
		want []model.Param
	}{
		{
			name: "no params",
			code: `
			func test() {}
			`,
			want: nil,
		},
		{
			name: "one param",
			code: `
			func test(a int) {}
			`,
			want: []model.Param{
				{
					Name: "a",
					Type: "int",
				},
			},
		},
		{
			name: "two params",
			code: `
			func test(a int, b string) {}
			`,
			want: []model.Param{
				{
					Name: "a",
					Type: "int",
				},
				{
					Name: "b",
					Type: "string",
				},
			},
		},
		{
			name: "params with import",
			code: `
			import "context"
			func test(ctx context.Context, a *int) {}
			`,
			want: []model.Param{
				{
					Name: "ctx",
					Type: "context.Context",
				},
				{
					Name: "a",
					Type: "*int",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parseCodeNode(t, test.code, func(file *ast.File) func(ast.Node) bool {
				return func(n ast.Node) bool {
					if fn, ok := n.(*ast.FuncType); ok {
						got := parseParams(fn.Params)
						require.Equal(t, test.want, got)
						return false
					}
					return true
				}
			})
		})
	}
}

func TestHasEmbededInterface(t *testing.T) {
	typ := "genapi.Interface"

	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "no embeded interface",
			code: `
			type Interface interface {
				Method()
			}
			`,
			want: false,
		},
		{
			name: "embeded interface",
			code: `
			type Interface interface {
				genapi.Interface
				Method()
			}
			`,
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var found bool

			parseCodeNode(t, test.code, func(file *ast.File) func(ast.Node) bool {
				return func(n ast.Node) bool {
					if typeSpec, ok := n.(*ast.TypeSpec); ok {
						if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							found = hasEmbededInterface(interfaceType, typ)
							return false
						}
					}
					return true
				}
			})

			require.Equal(t, test.want, found)
		})
	}
}

func TestCollectImports(t *testing.T) {
	tests := []struct {
		name string
		code string
		want []string
	}{
		{
			name: "no imports",
		},
		{
			name: "one import",
			code: `
			import "context"
			`,
			want: []string{
				`"context"`,
			},
		},
		{
			name: "many imports",
			code: `
			import (
				"context"
				"github.com/lexcao/genapi"
			)
			`,
			want: []string{
				`"context"`,
				`"github.com/lexcao/genapi"`,
			},
		},
		{
			name: "many imports with newline",
			code: `
			import (
				"context"

				"github.com/lexcao/genapi"
			)
			`,
			want: []string{
				`"context"`,
				`"github.com/lexcao/genapi"`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parseCodeNode(t, test.code, func(file *ast.File) func(ast.Node) bool {
				return func(n ast.Node) bool {
					got := collectImports(file.Imports)
					require.Equal(t, test.want, got)
					return false
				}
			})
		})
	}
}

func parseCodeNode(t *testing.T, code string, fn func(*ast.File) func(ast.Node) bool) {
	t.Helper()

	file, err := parser.ParseFile(token.NewFileSet(), "", fmt.Sprintf("package main\n%s", code), parser.ParseComments)
	require.NoError(t, err)

	ast.Inspect(file, fn(file))
}
