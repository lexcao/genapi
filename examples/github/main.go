package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/lexcao/genapi"
	"github.com/lexcao/genapi/examples/github/api"
)

func main() {
	client, err := genapi.New[api.GitHub]()
	if err != nil {
		fmt.Printf("failed to create client: %v\n", err)
		return
	}

	contributors, err := client.Contributors(context.Background(), "lexcao", "genapi")
	if err != nil {
		fmt.Printf("failed to get contributors: %v\n", err)
		var apiErr *genapi.Error
		if errors.As(err, &apiErr) {
			body, _ := io.ReadAll(apiErr.Response.Body)
			fmt.Printf("API error: %v\n", string(body))
		}
		return
	}

	for _, contributor := range contributors {
		fmt.Printf("%s (%d)\n", contributor.Login, contributor.Contributions)
	}
}
