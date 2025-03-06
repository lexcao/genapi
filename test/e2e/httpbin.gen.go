// CODE GENERATED BY genapi. DO NOT EDIT.
package e2e

import (
	"github.com/lexcao/genapi"
	"net/http"
	"net/url"
)

type implHttpBin struct {
	client genapi.HttpClient
}

// SetHttpClient implments genapi.Interface
func (i *implHttpBin) SetHttpClient(client genapi.HttpClient) {
	i.client = client
}

func (i *implHttpBin) Get(value string) (*Response, error) {
	resp, err := i.client.Do(&genapi.Request{
		Method: "GET",
		Path:   "/get",
		Queries: url.Values{
			"key": []string{
				"value",
			},
			"value": []string{
				value,
			},
		},
		Header: http.Header{
			"X-Hello": []string{
				"world",
			},
		},
	})
	return genapi.HandleResponse[*Response](resp, err)
}

func (i *implHttpBin) Post(body *Body) (*Response, error) {
	resp, err := i.client.Do(&genapi.Request{
		Method: "POST",
		Path:   "/post",
		Body:   body,
	})
	return genapi.HandleResponse[*Response](resp, err)
}

func init() {
	genapi.Register[HttpBin, *implHttpBin](
		genapi.Config{
			BaseURL: "https://httpbin.org",
			Header: http.Header{
				"Global-Header": []string{
					"global-value",
				},
			},
		},
	)
}
