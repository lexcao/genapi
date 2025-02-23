package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/model"
	"github.com/lexcao/genapi/internal/parser/annotation"
)

func TestBindPath(t *testing.T) {
	testBind(t, &pathBinding{}, []testCase{
		{
			name: "OK",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Path: "/path",
					},
				},
			},
			expectedBindings: model.Bindings{
				Path: "/path",
			},
		},
	})
}
