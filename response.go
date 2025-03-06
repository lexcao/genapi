package genapi

import (
	"encoding/json"
	"fmt"
)

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
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, fmt.Errorf("failed to decode response body: %w", err)
		}

		return result, nil
	}

	return result, &Error{Response: resp}
}

func MustHandleResponse[T any](resp *Response, err error) T {
	result, err := HandleResponse[T](resp, err)
	if err != nil {
		panic(err)
	}
	return result
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
