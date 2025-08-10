package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/lexcao/genapi/internal"
)

var DefaultClient = New(http.DefaultClient)

func New(client *http.Client) *HttpClient {
	return &HttpClient{client: client}
}

type HttpClient struct {
	config  internal.Config
	client  *http.Client
	baseURL *url.URL
}

func (c *HttpClient) SetConfig(config internal.Config) {
	if c.client == nil {
		c.client = http.DefaultClient
	}
	
	// Parse and cache base URL - panic on invalid since build-time validation should catch these
	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		panic(fmt.Sprintf("genapi: invalid base URL '%s': %v (this should have been caught at build time)", config.BaseURL, err))
	}
	
	c.config = config
	c.baseURL = baseURL
}

func (c *HttpClient) Do(req *internal.Request) (*internal.Response, error) {
	ctx := req.Context
	if ctx == nil {
		ctx = context.TODO()
	}

	var body io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(bodyBytes)
	}

	url, err := resolveURL(*c.baseURL, req.Path, req.PathParams)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header = make(http.Header)
	for k, v := range c.config.Header {
		httpReq.Header[k] = v
	}
	for k, v := range req.Header {
		httpReq.Header[k] = v
	}

	httpReq.URL.RawQuery = req.Queries.Encode()

	return c.client.Do(httpReq)
}

func resolveURL(baseURL url.URL, path string, pathParams map[string]string) (string, error) {
	var params []string
	for key, value := range pathParams {
		params = append(params, "{"+key+"}", url.PathEscape(value))
	}
	replacer := strings.NewReplacer(params...)

	part, err := url.Parse(replacer.Replace(path))
	if err != nil {
		return "", err
	}

	return baseURL.ResolveReference(part).String(), nil
}
