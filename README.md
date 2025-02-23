# genapi

genapi is a API client like OpenFeign(https://github.com/OpenFeign/feign) in Go.


## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/lexcao/genapi"
)

//go:generate genapi

// @BaseURL("https://api.github.com")
// @Header("Accept", "application/vnd.github.v3+json")
type GitHub interface {
    genapi.Interface

	// @GET("/repos/{owner}")
    // @Query("sort", "{sort}")
	// ListRepositories list repositories for specific user
	ListRepositories(ctx context.Context, owner string, sort string) ([]Repository, error)
}

type Repository struct {
	Name string `json:"name"`
}

func main() {
    client := genapi.New[GitHub](
        genapi.WithHeader("X-Auth-Token", "GITHUB_TOKEN"),
    )

    repositories, err := client.ListRepositories(context.Background(), "octocat", "desc")
    if err != nil {
		fmt.Errorf("failed to list repositories: %w", err)
    }
    fmt.Println(repositories)
}
```

## Annotations

- case-insensitive
- @BaseURL("https://api.github.com")
- @Header("Accept", "application/vnd.github.v3+json")
- @GET, @POST, @PUT, @DELETE, @PATCH, @OPTIONS, @HEAD http methods
- @Query("sort", "{sort}")


Dependencies:
- App -> Generator -> Binder -> Parser -> Model -> Annotation
