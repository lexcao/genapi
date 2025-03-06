package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
)

func TestBindHeader(t *testing.T) {
	testBind(t, &headerBinding{}, []testCase{
		{
			name: "empty header",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{},
				},
			},
			expectedBindings: model.MethodBindings{},
		},
		{
			name: "one header",
			given: model.Method{
				Params: []model.Param{
					{Name: "token", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Header: `http.Header{
	"Authorization": []string{
		token,
	},
}`,
			},
			expectedBinded:  []string{"token"},
			expectedImports: []string{`"net/http"`},
		},
		{
			name: "multiple headers",
			given: model.Method{
				Params: []model.Param{
					{Name: "token", Type: "string"},
					{Name: "type", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
						{Key: "Content-Type", Values: []annotation.Variable{"{type}"}},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Header: `http.Header{
	"Authorization": []string{
		token,
	},
	"Content-Type": []string{
		type,
	},
}`,
			},
			expectedBinded:  []string{"token", "type"},
			expectedImports: []string{`"net/http"`},
		},
		{
			name: "header with constant value",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Accept", Values: []annotation.Variable{"application/json"}},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Header: `http.Header{
	"Accept": []string{
		"application/json",
	},
}`,
			},
			expectedImports: []string{`"net/http"`},
		},
		{
			name: "header param not found",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
					},
				},
			},
			expectedError: &ErrNotFound{Type: "header", Value: "token"},
		},
		{
			name: "multiple values for same header",
			given: model.Method{
				Params: []model.Param{
					{Name: "type1", Type: "string"},
					{Name: "type2", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Accept", Values: []annotation.Variable{"{type1}", "{type2}"}},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Header: `http.Header{
	"Accept": []string{
		type1,
		type2,
	},
}`,
			},
			expectedBinded:  []string{"type1", "type2"},
			expectedImports: []string{`"net/http"`},
		},
		{
			name: "mixed constant and variable values",
			given: model.Method{
				Params: []model.Param{
					{Name: "token", Type: "string"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Accept", Values: []annotation.Variable{"application/json"}},
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Header: `http.Header{
	"Accept": []string{
		"application/json",
	},
	"Authorization": []string{
		token,
	},
}`,
			},
			expectedBinded:  []string{"token"},
			expectedImports: []string{`"net/http"`},
		},
		{
			name: "integer header param",
			given: model.Method{
				Params: []model.Param{
					{Name: "apiVersion", Type: "int"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Api-Version", Values: []annotation.Variable{"{apiVersion}"}},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Header: `http.Header{
	"Api-Version": []string{
		strconv.Itoa(int(apiVersion)),
	},
}`,
			},
			expectedBinded:  []string{"apiVersion"},
			expectedImports: []string{`"net/http"`, `"strconv"`},
		},
		{
			name: "float header param",
			given: model.Method{
				Params: []model.Param{
					{Name: "rate", Type: "float64"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "X-Rate-Limit", Values: []annotation.Variable{"{rate}"}},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Header: `http.Header{
	"X-Rate-Limit": []string{
		strconv.FormatFloat(float64(rate), 'f', -1, 64),
	},
}`,
			},
			expectedBinded:  []string{"rate"},
			expectedImports: []string{`"net/http"`, `"strconv"`},
		},
		{
			name: "boolean header param",
			given: model.Method{
				Params: []model.Param{
					{Name: "debug", Type: "bool"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "X-Debug-Mode", Values: []annotation.Variable{"{debug}"}},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Header: `http.Header{
	"X-Debug-Mode": []string{
		strconv.FormatBool(debug),
	},
}`,
			},
			expectedBinded:  []string{"debug"},
			expectedImports: []string{`"net/http"`, `"strconv"`},
		},
		{
			name: "mixed type header params",
			given: model.Method{
				Params: []model.Param{
					{Name: "token", Type: "string"},
					{Name: "version", Type: "int"},
					{Name: "debug", Type: "bool"},
				},
				Annotations: annotation.MethodAnnotations{
					Headers: []annotation.Header{
						{Key: "Authorization", Values: []annotation.Variable{"{token}"}},
						{Key: "X-Api-Version", Values: []annotation.Variable{"{version}"}},
						{Key: "X-Debug-Mode", Values: []annotation.Variable{"{debug}"}},
					},
				},
			},
			expectedBindings: model.MethodBindings{
				Header: `http.Header{
	"Authorization": []string{
		token,
	},
	"X-Api-Version": []string{
		strconv.Itoa(int(version)),
	},
	"X-Debug-Mode": []string{
		strconv.FormatBool(debug),
	},
}`,
			},
			expectedBinded:  []string{"token", "version", "debug"},
			expectedImports: []string{`"net/http"`, `"strconv"`},
		},
	})
}
