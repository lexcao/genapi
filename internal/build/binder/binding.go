package binder

import (
	"github.com/lexcao/genapi/internal/build/common"
	"github.com/lexcao/genapi/internal/build/model"
)

type binding interface {
	Name() string
	Bind(ctx *context) error
}

type context struct {
	Method       *model.Method
	ParamsByName map[string]model.Param
	BindedParams common.Set[string]
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
	method.Bindings = &model.MethodBindings{}

	paramsByName := map[string]model.Param{}
	for _, param := range method.Params {
		paramsByName[param.Name] = param
	}

	return &context{
		Method:       method,
		ParamsByName: paramsByName,
	}
}
