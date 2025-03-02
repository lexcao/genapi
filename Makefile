.PHONY: all test lint coverage generate clean website

all: lint test

test:
	go test -v -race ./...

test-e2e:
	E2E_TEST=true go test -v -race ./test/e2e

lint:
	golangci-lint run --out-format=colored-line-number --new=false --new-from-rev=

coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

generate:
	go generate ./...

clean:
	rm -f coverage.txt coverage.html
	go clean 

website:
	go run ./website/server.go -dir ./website
