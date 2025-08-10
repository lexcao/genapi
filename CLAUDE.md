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

## GitHub Workflow

### Remote Setup
```bash
# Add upstream remote (if working with a fork)
git remote add upstream https://github.com/original-repo/genapi.git

# Verify remotes
git remote -v
# origin    https://github.com/your-fork/genapi.git (fetch)
# origin    https://github.com/your-fork/genapi.git (push) 
# upstream  https://github.com/original-repo/genapi.git (fetch)
# upstream  https://github.com/original-repo/genapi.git (push)
```

### Branch Management
```bash
# Sync with upstream before creating feature branches
git checkout main
git fetch upstream
git merge upstream/main
git push origin main

# Create feature branch
git checkout -b feature/your-feature-name

# Push feature branch to origin
git push -u origin feature/your-feature-name
```

### Pull Request Workflow
```bash
# Before creating PR, ensure feature branch is up to date
git checkout main
git fetch upstream  
git merge upstream/main
git checkout feature/your-feature-name
git rebase main

# Run all checks before pushing
make all

# Push updated feature branch
git push origin feature/your-feature-name --force-with-lease

# Create PR targeting upstream/main
gh pr create --base main --head your-fork:feature/your-feature-name \
  --title "feat: your feature description" \
  --body "Description of changes"
```

### Sync with Upstream
```bash
# Regular sync to stay current with upstream
git checkout main
git fetch upstream
git merge upstream/main
git push origin main

# Update existing feature branches
git checkout feature/your-feature-name  
git rebase main
git push origin feature/your-feature-name --force-with-lease
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

## Task Management

For each task, you should propose a solution for review, DO NOT start coding without approval.

