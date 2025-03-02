package binder

import (
	"fmt"
	"net/textproto"

	"github.com/lexcao/genapi/internal/build/binder/printer"
	"github.com/lexcao/genapi/internal/build/model"
	"github.com/lexcao/genapi/internal/build/parser/annotation"
)

type bindedVariablesPrinter []model.BindedVariable

func (v bindedVariablesPrinter) Name() string {
	return "[]string"
}

func (v bindedVariablesPrinter) Print(p *printer.Printer) {
	for _, variable := range v {
		p.Indent()
		printVariable(p, variable)
		p.Unquoted(",")
		p.NewLine()
	}
}

func printVariable(p *printer.Printer, variable model.BindedVariable) {
	if variable.IsVariable() && variable.Param != nil {
		p.Unquoted(tryConvert(variable.Param))
	} else {
		p.Quote(variable.String())
	}
}

func tryConvert(param *model.Param) string {
	switch param.Type {
	case "int", "int8", "int16", "int32", "int64":
		return fmt.Sprintf("strconv.Itoa(int(%s))", param.Name)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("strconv.FormatUint(uint64(%s), 10)", param.Name)
	case "float32", "float64":
		return fmt.Sprintf("strconv.FormatFloat(float64(%s), 'f', -1, 64)", param.Name)
	case "bool":
		return fmt.Sprintf("strconv.FormatBool(%s)", param.Name)
	default:
		return param.Name
	}
}

type bindedHeaderPrinter struct {
	orderBy []annotation.Header
	binded  map[string]bindedVariablesPrinter
}

func (v bindedHeaderPrinter) Name() string {
	return "http.Header"
}

func (v bindedHeaderPrinter) Print(p *printer.Printer) {
	seen := map[string]bool{}
	for _, header := range v.orderBy {
		key := textproto.CanonicalMIMEHeaderKey(header.Key)
		values, ok := v.binded[key]
		if !ok {
			continue
		}
		if seen[key] {
			continue
		}
		seen[key] = true

		p.KeyValueLine(func(p *printer.Printer) {
			p.Quote(key)
		}, func(p *printer.Printer) {
			p.Item(values)
		})
	}
}
