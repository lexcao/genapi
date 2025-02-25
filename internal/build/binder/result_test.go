package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
)

func TestBindResult(t *testing.T) {
	testBind(t, &resultBinding{}, []testCase{
		{
			name:             "empty result",
			given:            model.Method{},
			expectedBindings: model.MethodBindings{},
		},
		{
			name: "single result with any",
			given: model.Method{
				Results: []model.Param{
					{Type: "Result"},
				},
			},
			expectedBindings: model.MethodBindings{
				Results: &model.BindingResults{
					Assignment: "resp, err",
					Statement:  "genapi.MustHandleResponse[Result](resp, err)",
				},
			},
		},
		{
			name: "single result with Response",
			given: model.Method{
				Results: []model.Param{
					{Type: "genapi.Response"},
				},
			},
			expectedBindings: model.MethodBindings{
				Results: &model.BindingResults{
					Assignment: "resp, _",
					Statement:  "*resp",
				},
			},
		},
		{
			name: "single result with *Response",
			given: model.Method{
				Results: []model.Param{
					{Type: "*genapi.Response"},
				},
			},
			expectedBindings: model.MethodBindings{
				Results: &model.BindingResults{
					Assignment: "resp, _",
					Statement:  "resp",
				},
			},
		},
		{
			name: "single result with error",
			given: model.Method{
				Results: []model.Param{
					{Type: "error"},
				},
			},
			expectedBindings: model.MethodBindings{
				Results: &model.BindingResults{
					Assignment: "resp, err",
					Statement:  "genapi.HandleResponse0(resp, err)",
				},
			},
		},
		{
			name: "multiple result with any",
			given: model.Method{
				Results: []model.Param{
					{Type: "Result"},
					{Type: "error"},
				},
			},
			expectedBindings: model.MethodBindings{
				Results: &model.BindingResults{
					Assignment: "resp, err",
					Statement:  "genapi.HandleResponse[Result](resp, err)",
				},
			},
		},
		{
			name: "multiple result with Response",
			given: model.Method{
				Results: []model.Param{
					{Type: "genapi.Response"},
					{Type: "error"},
				},
			},
			expectedBindings: model.MethodBindings{
				Results: &model.BindingResults{
					Assignment: "resp, err",
					Statement:  "*resp, err",
				},
			},
		},
		{
			name: "multiple result with *Response",
			given: model.Method{
				Results: []model.Param{
					{Type: "*genapi.Response"},
					{Type: "error"},
				},
			},
			expectedBindings: model.MethodBindings{
				Results: &model.BindingResults{
					Assignment: "resp, err",
					Statement:  "resp, err",
				},
			},
		},
		{
			name: "multiple result with unsupported primitive",
			given: model.Method{
				Results: []model.Param{
					{Type: "string"},
					{Type: "error"},
				},
			},
			expectedError: &ErrUnsupported{Message: "type [string] is not supported as a result type"},
		},
	})
}
