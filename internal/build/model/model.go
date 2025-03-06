package model

import (
	"github.com/lexcao/genapi/internal/build/common"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
)

type Interface struct {
	Name        string
	Package     string
	Imports     common.Set[string]
	Methods     []Method
	Annotations annotation.InterfaceAnnotations
	Bindings    *InterfaceBindings
}

type InterfaceBindings struct {
	Config string
}

type Method struct {
	Name        string
	Interface   *Interface
	Params      []Param
	Results     []Param
	Annotations annotation.MethodAnnotations
	Bindings    *MethodBindings
}

type BindedVariable struct {
	annotation.Variable
	Param *Param
}

// MethodBindings bind Method and Annotation to genapi.Request
type MethodBindings struct {
	Method     string
	Path       string
	Body       string
	Context    string
	Queries    string
	Header     string
	PathParams string
	Imports    string
	Results    *BindingResults
}

type BindingResults struct {
	Assignment string
	Statement  string
}

type Param struct {
	Name string
	Type string
}
