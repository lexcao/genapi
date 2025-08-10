package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/lexcao/genapi/internal/build"
)

func main() {
	// 1. Parse command flags
	var config build.Config
	var fileModeStr string
	flag.StringVar(&config.Filename, "file", "", "Input file path")
	flag.StringVar(&fileModeStr, "filemode", "0600", "File permissions for generated files (octal, e.g., 0644)")
	flag.Parse()

	if config.Filename == "" {
		fmt.Println("Error: input file is required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse file mode
	if fileModeStr != "" {
		mode, err := strconv.ParseUint(fileModeStr, 8, 32)
		if err != nil {
			fmt.Printf("Error: invalid file mode '%s': %v\n", fileModeStr, err)
			os.Exit(1)
		}
		config.FileMode = os.FileMode(mode)
	}

	if err := build.Run(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
