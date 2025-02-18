package model

import "github.com/lexcao/genapi/internal/parser/annotation"

type Interface struct {
	Name        string
	Package     string
	Methods     []Method
	Annotations annotation.InterfaceAnnotations
}

type Method struct {
	Name        string
	Interface   string
	Params      []Param
	Results     []Param
	Annotations annotation.MethodAnnotations
}

type Param struct {
	Name string
	Type string
}
