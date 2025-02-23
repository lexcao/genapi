package build

import (
	"os"
	"strings"

	"github.com/lexcao/genapi/internal/build/binder"
	"github.com/lexcao/genapi/internal/build/generator"
	"github.com/lexcao/genapi/internal/build/parser"
)

type Config struct {
	Filename string
	Output   string
}

func Run(config Config) error {
	interfaces, err := parser.ParseFile(config.Filename)
	if err != nil {
		return err
	}

	for _, iface := range interfaces {
		if err := binder.Bind(&iface); err != nil {
			return err
		}
	}

	content, err := generator.GenerateFile(config.Filename, interfaces)
	if err != nil {
		return err
	}

	output := strings.TrimSuffix(config.Filename, ".go") + ".gen.go"
	if config.Output != "" {
		output = config.Output
	}

	return os.WriteFile(output, content, 0644)
}
