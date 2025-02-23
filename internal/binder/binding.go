package binder

import (
	"strings"

	"github.com/lexcao/genapi/internal/model"
	"github.com/lexcao/genapi/internal/parser/annotation"
)

type binding interface {
	Name() string
	Bind(ctx *context) error
}

type context struct {
	Method       *model.Method
	ParamsByName map[string]model.Param
	BindedParams map[string]struct{}
}

var bindings = []binding{
	&methodBinding{},
	&pathBinding{},
	&pathParamBinding{},
	&queryBinding{},
	&headerBinding{},
	&contextBinding{},
	&bodyBinding{},
	&resultBinding{},
}

func newBindingContext(method *model.Method) *context {
	method.Bindings = &model.Bindings{}

	paramsByName := map[string]model.Param{}
	for _, param := range method.Params {
		paramsByName[param.Name] = param
	}

	return &context{
		Method:       method,
		ParamsByName: paramsByName,
		BindedParams: make(map[string]struct{}),
	}
}

// input - is a string value of a map
// value - is a variable param like {owner}
// example:
// input: `map[string]string{"owner": "{owner}", "repo": "{repo}"}`
// ouput: `map[string]string{"owner": owner, "repo": repo}`
func replaceVariable(input string, value annotation.Variable) string {
	if !value.IsVariable() {
		return input
	}

	v := value.Escape()
	return strings.ReplaceAll(input, `"{`+v+`}"`, v)
}
