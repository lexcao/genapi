package it

import (
	"fmt"
	"os"
	"testing"

	"github.com/lexcao/genapi/internal/build"
)

func TestMain(m *testing.M) {
	if err := build.Run(build.Config{
		Filename: "api.go",
	}); err != nil {
		fmt.Printf("failed to build API: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
