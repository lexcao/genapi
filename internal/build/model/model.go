package model

import (
	"github.com/lexcao/genapi/internal/build/parser/annotation"
)

type Interface struct {
	Name        string
	Package     string
	Imports     []string
	Methods     []Method
	Annotations annotation.InterfaceAnnotations
}

type Method struct {
	Name        string
	Interface   string
	Params      []Param
	Results     []Param
	Annotations annotation.MethodAnnotations
	Bindings    *Bindings
}

// Bindings bind Method and Annotation to genapi.Request
type Bindings struct {
	Method     string
	Path       string
	Body       string
	Context    string
	Queries    string
	Headers    string
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
