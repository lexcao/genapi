# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

genapi is a declarative HTTP client generator for Go that creates HTTP clients from interface definitions with annotations. It follows a code generation approach where interfaces are parsed, annotated with HTTP method details, and compiled into client implementations.

## Key Architecture Components

- **Code Generation Pipeline**: `internal/build/` contains the three-stage pipeline:
  - `parser/` - Parses Go interface definitions and annotations
  - `binder/` - Binds parsed data to internal representations 
  - `generator/` - Generates concrete client implementations
- **Runtime System**: `internal/runtime/` provides the runtime client factory and registry
- **HTTP Client Abstraction**: `pkg/clients/` contains pluggable HTTP client implementations (default `http`, optional `resty`)
- **Interface Contract**: All generated clients implement `genapi.Interface` which requires `SetHttpClient(HttpClient)`

## Development Commands

### Core Development
```bash
# Run all tests across all Go modules
make test

# Run linters across all Go modules  
make lint

# Generate client code from interface definitions
make generate
# or directly: go generate ./...

# Run all checks (lint + test)
make all
```

### Single Test Execution
```bash
# Run tests in specific directory
go test -v -race -count=1 ./internal/...

# Run specific test
go test -v -run TestSpecificFunction ./path/to/package
```

### Code Generation
```bash
# Generate client from interface file
go run github.com/lexcao/genapi/cmd/genapi -file path/to/api.go

# This creates api.gen.go with the concrete client implementation
```

### Website/Examples
```bash
# Run development website
make website
# or: go run ./website/server.go -dir ./website
```

## Module Structure

This is a multi-module repository with separate `go.mod` files in:
- Root module: Core genapi library
- `pkg/clients/resty/`: Optional Resty HTTP client implementation  
- `website/`: Development website and examples

The Makefile automatically discovers and operates on all modules.

## Generated Code Pattern

Interfaces with `//go:generate` directives and genapi annotations become concrete clients:
- Input: Interface with `@BaseURL`, `@GET`, `@POST` etc. annotations
- Output: `*.gen.go` files with `Register` calls and client implementations
- Runtime: `genapi.New[InterfaceType]()` creates configured client instances

## HTTP Client Architecture

The system supports pluggable HTTP clients through the `HttpClient` interface. Default is `pkg/clients/http`, with `pkg/clients/resty` as an alternative. Custom implementations must implement `SetConfig(Config)` and `Do(*Request) (*Response, error)`.