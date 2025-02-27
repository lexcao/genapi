.PHONY: all test lint coverage generate clean

all: lint test

test:
	go test -v -race ./...

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