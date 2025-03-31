package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var mockTransport = &mockTransportImpl{
	responders: make(map[string]roundTripFunc),
}

type roundTripFunc = func(*http.Request) (*http.Response, error)

type mockTransportImpl struct {
	responders map[string]roundTripFunc
}

func (m *mockTransportImpl) Activate(client *http.Client) func() {
	client.Transport = m
	return func() {
		m.Reset()
		client.Transport = nil
	}
}

func (m *mockTransportImpl) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.String()
	if responder, ok := m.responders[key]; ok {
		return responder(req)
	}
	return nil, fmt.Errorf("no mock for %s", key)
}

func (m *mockTransportImpl) RegisterResponder(method, path string, responder roundTripFunc) {
	m.responders[method+" "+path] = responder
}

func (m *mockTransportImpl) Reset() {
	m.responders = make(map[string]roundTripFunc)
}

func (m *mockTransportImpl) stringResponder(status int, body string) roundTripFunc {
	return func(req *http.Request) (*http.Response, error) {
		return m.stringResponse(status, body), nil
	}
}

func (m *mockTransportImpl) stringResponse(status int, body string) *http.Response {
	return &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
