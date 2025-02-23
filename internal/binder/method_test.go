package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/model"
	"github.com/lexcao/genapi/internal/parser/annotation"
)

func TestBindMethod(t *testing.T) {
	testBind(t, &methodBinding{}, []testCase{
		{
			name: "OK",
			given: model.Method{
				Annotations: annotation.MethodAnnotations{
					RequestLine: annotation.RequestLine{
						Method: "GET",
					},
				},
			},
			expectedBindings: model.Bindings{
				Method: "GET",
			},
		},
	})
}
