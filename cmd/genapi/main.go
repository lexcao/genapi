package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lexcao/genapi/internal/build"
)

func main() {
	// 1. Parse command flags
	var config build.Config
	flag.StringVar(&config.Filename, "file", "", "Input file path")
	flag.Parse()

	if config.Filename == "" {
		fmt.Println("Error: input file is required")
		flag.Usage()
		os.Exit(1)
	}

	if err := build.Run(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
