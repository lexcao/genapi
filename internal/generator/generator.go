package generator

import (
	"bytes"
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
	buf.WriteString("// Code generated by genapi. DO NOT EDIT.\n\n")
	buf.WriteString("package " + interfaces[0].Package + "\n\n")

	for _, iface := range interfaces {
		if err := templates.ExecuteTemplate(&buf, tmplInterface, iface); err != nil {
			return err
		}
	}

	fomatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	newFilename := strings.TrimSuffix(filename, ".go") + ".gen.go"
	return os.WriteFile(newFilename, fomatted, 0644)
}
