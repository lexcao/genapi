.PHONY: all test lint generate website

GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort)

define run_in_modules
	@set -e; for dir in $(GO_MOD_DIRS); do \
	  echo "\033[1;34m>> Found go.mod in \033[1;32m$${dir}\033[0m"; \
	  (cd "$${dir}" && $(1)); \
	done
endef

all: lint test

test:
	$(call run_in_modules,go mod tidy && go test -v -race -count=1 ./...)

lint:
	$(call run_in_modules,go mod tidy && golangci-lint run --out-format=colored-line-number --new=false --new-from-rev=)

generate:
	go generate ./...

website:
	go run ./website/server.go -dir ./website
