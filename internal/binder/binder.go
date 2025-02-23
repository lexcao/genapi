package binder

import (
	"github.com/lexcao/genapi/internal/model"
)

func Bind(interfaceModel *model.Interface) error {
	for i := range interfaceModel.Methods {
		if err := BindMethod(&interfaceModel.Methods[i]); err != nil {
			return err
		}
	}

	return nil
}

func BindMethod(method *model.Method) error {
	bindCtx := newBindingContext(method)

	for _, binding := range bindings {
		if err := binding.Bind(bindCtx); err != nil {
			return err
		}
	}

	return nil
}
