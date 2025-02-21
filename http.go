package genapi // TODO: check if we need to update package

import (
	"net/http"
	"net/url"
)

type DefaultHttpClient struct {
	*http.Client
}

func (c *DefaultHttpClient) Do(req *Request) (*Response, error) {
	// mapping to http.Request
	httpReq := &http.Request{
		Method: req.Method,
		URL:    &url.URL{Path: req.Path},
	}
	return c.Client.Do(httpReq)
}
