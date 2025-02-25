package binder

import (
	"github.com/lexcao/genapi/internal/build/model"
)

func Bind(iface *model.Interface) error {
	for i := range iface.Methods {
		if err := BindMethod(&iface.Methods[i]); err != nil {
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
