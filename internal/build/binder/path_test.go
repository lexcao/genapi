package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
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
			expectedBindings: model.MethodBindings{
				Path: "/path",
			},
		},
	})
}
