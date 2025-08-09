# Use Resty as genapi HttpClient


## Install

```bash
go get github.com/lexcao/genapi/pkg/clients/resty
```

## Usage

```go
package main

import (
	"github.com/lexcao/genapi"
	"github.com/lexcao/genapi/pkg/clients/resty"
    resty_client "github.com/go-resty/resty/v2"
)

func main() {
	client, err := genapi.New[api.GitHub](
		genapi.WithHttpClient(resty.DefaultClient),           // default Resty client
        genapi.WithHttpClient(resty.New(resty_client.New())), // customized Resty client
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
}
```

## Example

Please check [example](./example/main.go) for more details.

```bash
$ cd ./pkg/clients/resty/example
$ go run main.go
```