package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/model"
)

func TestBindBody(t *testing.T) {
	testBind(t, &bodyBinding{}, []testCase{
		{
			name:             "empty",
			given:            model.Method{},
			expectedBindings: model.Bindings{},
		},
		{
			name: "empty with disallowed primitive",
			given: model.Method{
				Params: []model.Param{
					{Name: "num", Type: "int"},
				},
			},
			expectedBindings: model.Bindings{},
		},
		{
			name: "one body param",
			given: model.Method{
				Params: []model.Param{
					{Name: "req", Type: "Request"},
				},
			},
			expectedBindings: model.Bindings{
				Body: "req",
			},
			expectedBinded: []string{"req"},
		},
		{
			name: "duplicated body param",
			given: model.Method{
				Params: []model.Param{
					{Name: "req1", Type: "Request"},
					{Name: "req2", Type: "Request"},
				},
			},
			expectedError:  &ErrDuplicated{Type: "body", Value: "req2"},
			expectedBinded: []string{"req1"},
		},
	})
}
