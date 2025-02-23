package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"

	"github.com/lexcao/genapi/internal/model"
)

func GenerateFile(filename string, interfaces []model.Interface) error {
	if len(interfaces) == 0 {
		return nil
	}

	var buf bytes.Buffer

	var data = templateData{
		Imports:    interfaces[0].Imports,
		Package:    interfaces[0].Package,
		Interfaces: interfaces,
	}

	if err := templates.ExecuteTemplate(&buf, tmplMain, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fomatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format source: %w", err)
	}

	newFilename := strings.TrimSuffix(filename, ".go") + ".gen.go"
	return os.WriteFile(newFilename, fomatted, 0644)
}
