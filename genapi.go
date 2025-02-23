package genapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Interface interface {
	setHttpClient(HttpClient)
}

func New[T Interface](opts ...Option) T {
	panic("not implemented")
}

type Request struct {
	Body       any
	Method     string
	Path       string
	Queries    url.Values
	Headers    http.Header
	PathParams map[string]string
	Context    context.Context
}

type Response = http.Response

type HttpClient interface {
	Do(req *Request) (*Response, error)
}

type Error struct {
	Response *Response
}

func (e *Error) Error() string {
	return fmt.Sprintf("handle response error: %s", e.Response.Status)
}

func HandleResponse[T any](resp *Response, err error) (T, error) {
	var result T
	if err != nil {
		return result, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return result, fmt.Errorf("failed to read response body: %w", err)
		}
		defer resp.Body.Close()

		return result, json.Unmarshal(body, &result)
	}

	return result, &Error{Response: resp}
}

func HandleResponse0(resp *Response, err error) error {
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return &Error{Response: resp}
}
