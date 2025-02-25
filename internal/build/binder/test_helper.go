package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/model"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name             string
	given            model.Method
	expectedError    error
	expectedBindings model.MethodBindings
	expectedBinded   []string
	expectedImports  []string
}

func testBind(t *testing.T, binding binding, cases []testCase) {
	t.Helper()

	for _, c := range cases {
		t.Run(binding.Name()+":"+c.name, func(t *testing.T) {
			c.given.Interface = &model.Interface{}
			ctx := newBindingContext(&c.given)
			err := binding.Bind(ctx)
			if c.expectedError != nil {
				require.Error(t, err)
				require.EqualError(t, err, c.expectedError.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, c.expectedBindings, *c.given.Bindings)

			require.ElementsMatch(t, c.expectedBinded, ctx.BindedParams.Slices())
			require.ElementsMatch(t, c.expectedImports, ctx.Method.Interface.Imports.Slices())
		})
	}
}
