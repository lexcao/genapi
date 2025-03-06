package binder

import (
	"github.com/lexcao/genapi/internal/build/binder/printer"
	"github.com/lexcao/genapi/internal/build/model"
)

func BindInterface(iface *model.Interface) error {
	headers := map[string]bindedVariablesPrinter{}
	for _, header := range iface.Annotations.Headers {
		for _, value := range header.Values {
			headers[header.Key] = append(headers[header.Key], model.BindedVariable{
				Variable: value,
			})
		}
	}

	config := printer.PrintWith("genapi.Config", func(p *printer.Printer) {
		if iface.Annotations.BaseURL.Value != "" {
			p.KeyValueLine(func(p *printer.Printer) {
				p.Unquoted("BaseURL")
			}, func(p *printer.Printer) {
				p.Quote(iface.Annotations.BaseURL.Value)
			})
		}

		if len(headers) > 0 {
			iface.Imports.Add(`"net/http"`)

			p.KeyValueLine(func(p *printer.Printer) {
				p.Unquoted("Header")
			}, func(p *printer.Printer) {
				p.Item(bindedHeaderPrinter{
					orderBy: iface.Annotations.Headers,
					binded:  headers,
				})
			})
		}
	})

	if len(headers) == 0 && iface.Annotations.BaseURL.Value == "" {
		config = "genapi.Config{}"
	}

	iface.Bindings = &model.InterfaceBindings{
		Config: config,
	}

	return nil
}
