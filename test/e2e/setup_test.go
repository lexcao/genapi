package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("E2E_TEST") != "true" {
		os.Exit(0)
	}

	// get current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get current directory: %v", err)
		os.Exit(1)
	}

	// build genapi
	if err := runCommand("go", "build", "-o", "genapi", path.Join(dir, "../../cmd/genapi/main.go")); err != nil {
		fmt.Printf("failed to build genapi: %v", err)
		os.Exit(1)
	}

	// run go generate
	if err := runCommand("go", "generate", dir); err != nil {
		fmt.Printf("failed to run go generate: %v", err)
		os.Exit(1)
	}

	// run tests
	m.Run()

	// clear binary
	os.Remove("genapi")
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
