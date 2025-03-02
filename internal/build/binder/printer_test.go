package binder

import (
	"testing"

	"github.com/lexcao/genapi/internal/build/binder/printer"
	"github.com/lexcao/genapi/internal/build/model"
	"github.com/stretchr/testify/assert"
)

func TestBindedVariablesPrinter(t *testing.T) {
	variables := bindedVariablesPrinter{
		model.BindedVariable{Variable: "iter"},
		model.BindedVariable{Variable: "{str}", Param: &model.Param{Name: "str", Type: "string"}},
		model.BindedVariable{Variable: "{num}", Param: &model.Param{Name: "num", Type: "int"}},
		model.BindedVariable{Variable: "{boolean}", Param: &model.Param{Name: "boolean", Type: "bool"}},
	}

	expect := `[]string{
	"iter",
	str,
	strconv.Itoa(int(num)),
	strconv.FormatBool(boolean),
}`

	actual := printer.Print(variables)
	assert.Equal(t, expect, actual)
}
