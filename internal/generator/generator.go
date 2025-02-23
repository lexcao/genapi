package generator

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/lexcao/genapi/internal/model"
)

func GenerateFile(filename string, interfaces []model.Interface) ([]byte, error) {
	if len(interfaces) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer

	var data = templateData{
		Imports:    interfaces[0].Imports,
		Package:    interfaces[0].Package,
		Interfaces: interfaces,
	}

	if err := templates.ExecuteTemplate(&buf, tmplMain, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	fomatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format source: %w", err)
	}

	return fomatted, nil
}
