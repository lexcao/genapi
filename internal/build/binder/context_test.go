package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
)

func TestBindContext(t *testing.T) {
	testBind(t, &contextBinding{}, []testCase{
		{
			name:             "empty context",
			given:            model.Method{},
			expectedBindings: model.Bindings{},
		},
		{
			name: "context with value",
			given: model.Method{
				Params: []model.Param{
					{Name: "ctx", Type: "context.Context"},
				},
			},
			expectedBindings: model.Bindings{
				Context: "ctx",
			},
			expectedBinded: []string{"ctx"},
		},
		{
			name: "duplicated context",
			given: model.Method{
				Params: []model.Param{
					{Name: "ctx1", Type: "context.Context"},
					{Name: "ctx2", Type: "context.Context"},
				},
			},
			expectedError:  &ErrDuplicated{Type: "context.Context", Value: "ctx2"},
			expectedBinded: []string{"ctx1"},
		},
	})
}
