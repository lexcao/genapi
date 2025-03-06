# genapi

[![Go Reference](https://pkg.go.dev/badge/github.com/lexcao/genapi.svg)](https://pkg.go.dev/github.com/lexcao/genapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/lexcao/genapi)](https://goreportcard.com/report/github.com/lexcao/genapi)
[![CI](https://github.com/lexcao/genapi/actions/workflows/ci.yml/badge.svg)](https://github.com/lexcao/genapi/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/lexcao/genapi)](LICENSE)

genapi is a declarative HTTP client generator for Go, inspired by OpenFeign(https://github.com/OpenFeign/feign). It allows you to write HTTP clients using just interface definitions with annotations, eliminating the need for manual HTTP request handling.

## Features

- Declarative HTTP API client using Go interfaces
- Automatic JSON marshaling/unmarshaling
- Support for path/query parameters
- Custom header support
- Flexible response handling
- Context support for cancellation/timeouts
- **Use your favorate http client in runtime**

## Installation

```bash
go get github.com/lexcao/genapi
```

## Quick Start

1. Define your API interface:

```go
package api

import "github.com/lexcao/genapi"

//go:generate go run github.com/lexcao/genapi/cmd/genapi  -file $GOFILE

// @BaseURL("https://api.github.com")
// @Header("Accept", "application/vnd.github.v3+json")
type GitHub interface {
    genapi.Interface

    // @GET("/repos/{owner}/{repo}/contributors")
    Contributors(ctx context.Context, owner, repo string) ([]Contributor, error)
}

type Contributor struct {
    Login         string `json:"login"`
    Contributions int    `json:"contributions"`
}
```

2. Generate the client code:

```bash
$ go generate ./api
```

3. Use the client:

```go
import (
    "fmt"
    "github.com/lexcao/genapi"
    "your/package/to/api"
)

func main() {
    client := genapi.New[api.GitHub]()

    contributors, err := client.Contributors(context.Background(), "lexcao", "genapi")
    if err != nil {
        log.Fatalf("failed to get contributors: %v", err)
    }

    for _, c := range contributors {
        fmt.Printf("%s: %d contributions\n", c.Login, c.Contributions)
    }
}
```

## Core Feature: replace HttpClient in runtime

### Replace default HttpClient

```go
package main

import (
    "net/http"
    "github.com/lexcao/genapi"
)

func main() {
    httpClient := &http.Client{}

    client := genapi.New[api.GitHub](
        genapi.WithHttpClient(http.New(httpClient))
    )

    // or set client in the runtime
    client.SetHttpClient(httpClient)
}
```

### Use Resty as genapi HttpClient

Resty is as seperate pkg, you need to install first

```bash
$ go get github.com/lexcao/genapi/pkg/clients/resty
```

Then use as following,

```go
package main

import (
	"github.com/lexcao/genapi"
	"github.com/lexcao/genapi/pkg/clients/resty"
    resty_client "github.com/go-resty/resty/v2"
)

func main() {
	client := genapi.New[api.GitHub](
		genapi.WithHttpClient(resty.DefaultClient),           // default Resty client
        genapi.WithHttpClient(resty.New(resty_client.New())), // customized Resty client
	)
}
```

### Provide your own implmentation of HttpClient

You should implmentate the interface `genapi.HttpClient`
(You can follow the [resty](./pkg/clients/resty/resty.go) for reference)
```go
type HttpClient interface {
	SetConfig(Config)
	Do(req *Request) (*Response, error)
}
```

After implmentation, you should test your code to verify the base cases provided by the genapi
(You can follow the [resty_test](./pkg/clients/resty/resty_test.go) for reference)
```go
package client

func (c *HttpClient) GetClient() *http.Client {
	return c.client.GetClient()
}

func TestHttpClient(t *testing.T) {
	genapi.TestHttpClient(t, func() genapi.HttpClientTester { return DefaultClient })
}
```

## Configuration

The client can be configured with various options:

```go
client := genapi.New[api.GitHub](
    // Set runtime client
    genapi.WithHttpClient(resty.DefaultClient),

    // Set dynamic BaseURL
    genapi.WithBaseURL(os.GetEnv("API_ENDPOINT")),
    
    // Add global headers
    genapi.WithHeader(map[string]string{
        "Authorization": "Bearer " + token,
    }),
)
```

## Documentation

### Interface Requirements

Every interface must embed `genapi.Interface`:

```go
type MyAPI interface {
    genapi.Interface
    // your API methods here
}
```

### Annotations

#### Interface Level Annotations
| Annotation | Description                            | Example                                 |
| ---------- | -------------------------------------- | --------------------------------------- |
| @BaseURL   | Base URL for all API requests          | `@BaseURL("https://api.github.com")`    |
| @Header    | Global headers applied to all requests | `@Header("Accept", "application/json")` |

#### Method Level Annotations
| Annotation       | Description             | Example                                      |
| ---------------- | ----------------------- | -------------------------------------------- |
| @GET, @POST, etc | HTTP method and path    | `@GET("/users/{id}")`                        |
| @Query           | URL query parameters    | `@Query("sort", "{sort}")`                   |
| @Header          | Method-specific headers | `@Header("Authorization", "Bearer {token}")` |

### Response Types

genapi supports multiple response formats to fit your needs:

```go
// No response body
func DeleteUser(ctx context.Context, id string) error

// Typed response with error handling
func GetUser(ctx context.Context, id string) (User, error)

// Raw response access
func GetRawResponse(ctx context.Context) (*genapi.Response, error)

// Must-style response (panics on error)
func MustGetUser(ctx context.Context, id string) User
```

### Error Handling

You can access the Response from error

```go
err := client.DeleteUser(ctx, id)
var apiErr *genapi.Error
if errors.As(err, &apiErr) {
    // handle error with apiErr.Response
}
```

## Development

### Prerequisites

- Go 1.18 or higher
- Make (optional, for using Makefile commands)

### Setup

1. Clone the repository
```bash
$ git clone https://github.com/lexcao/genapi.git
$ cd genapi
```

2. Install dependencies
```bash
$ go mod download
```

3. Run tests
```bash
$ make test
```

### Available Make Commands

- `make test` - Run tests for all go modules
- `make lint` - Run linter for all go modules
- `make generate` - Run go generate
- `make clean` - Clean up generated files

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests.

## License

This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE) file for details.
