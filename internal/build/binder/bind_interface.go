package binder

import (
	"fmt"
	"strings"

	"github.com/lexcao/genapi/internal/build/model"
)

func BindInterface(iface *model.Interface) error {
	var builder strings.Builder

	builder.WriteString("genapi.Config{")
	{
		// BaseURL
		if iface.Annotations.BaseURL.Value != "" {
			builder.WriteRune('\n')
			builder.WriteString(fmt.Sprintf("\tBaseURL: %q,\n", iface.Annotations.BaseURL.Value))
		}
	}
	{
		// Headers
		if len(iface.Annotations.Headers) > 0 {
			iface.Imports.Add(`"net/http"`)

			builder.WriteString("\tHeaders: http.Header{\n")
			for _, header := range iface.Annotations.Headers {
				builder.WriteString(fmt.Sprintf("\t\t%q: []string{\n", header.Key))
				for _, value := range header.Values {
					builder.WriteString(fmt.Sprintf("\t\t\t%q,\n", value.String()))
				}
				builder.WriteString("\t\t},\n")
			}
			builder.WriteString("\t},\n")
		}
	}
	builder.WriteString("}")

	iface.Bindings = &model.InterfaceBindings{
		Config: builder.String(),
	}

	return nil
}
