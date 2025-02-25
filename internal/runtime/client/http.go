package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/lexcao/genapi/internal"
)

var DefaultClient = &HttpClient{
	client: http.DefaultClient,
}

type HttpClient struct {
	config internal.Config
	client *http.Client
}

func (c *HttpClient) SetConfig(config internal.Config) {
	c.client = http.DefaultClient
	c.config = config
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

	url, err := resolveURL(c.config.BaseURL, req.Path, req.PathParams)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header = make(http.Header)
	for k, v := range c.config.Headers {
		httpReq.Header[k] = v
	}
	for k, v := range req.Headers {
		httpReq.Header[k] = v
	}

	httpReq.URL.RawQuery = req.Queries.Encode()

	return c.client.Do(httpReq)
}

func resolveURL(baseURL string, path string, pathParams map[string]string) (string, error) {
	var params []string
	for key, value := range pathParams {
		params = append(params, "{"+key+"}", url.PathEscape(value))
	}
	replacer := strings.NewReplacer(params...)

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	part, err := url.Parse(replacer.Replace(path))
	if err != nil {
		return "", err
	}

	return base.ResolveReference(part).String(), nil
}
