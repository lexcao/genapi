package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

type HttpClientTester interface {
	HttpClient
	GetClient() *http.Client
}

func TestHttpClient(t *testing.T, createClient func() HttpClientTester) {
	t.Run("MockClient", func(t *testing.T) {
		setupMockClient := func(t *testing.T) HttpClient {
			client := createClient()
			t.Cleanup(mockTransport.Activate(client.GetClient()))
			return client
		}

		t.Run("SetConfig", func(t *testing.T) {
			t.Run("BaseURL", func(t *testing.T) {
				client := setupMockClient(t)
				client.SetConfig(Config{BaseURL: "https://lexcao.com/genapi"})

				mockTransport.RegisterResponder("GET", "https://lexcao.com/genapi",
					mockTransport.stringResponder(200, "OK"))

				resp, err := client.Do(&Request{Method: "GET"})

				RequireNoError(t, err)
				AssertEqual(t, http.StatusOK, resp.StatusCode)
			})

			t.Run("GlobalHeaders", func(t *testing.T) {
				client := setupMockClient(t)
				client.SetConfig(Config{
					Header: http.Header{
						"User-Agent": []string{"genapi/test_suite"},
					},
				})

				mockTransport.RegisterResponder("GET", "/genapi",
					func(req *http.Request) (*http.Response, error) {
						AssertEqual(t, "genapi/test_suite", req.Header.Get("User-Agent"))
						return mockTransport.stringResponse(200, "OK"), nil
					})

				resp, err := client.Do(&Request{Method: "GET", Path: "/genapi"})

				RequireNoError(t, err)
				AssertEqual(t, http.StatusOK, resp.StatusCode)
			})
		})

		t.Run("Do", func(t *testing.T) {
			testCases := []struct {
				name         string
				request      Request
				registerMock func()
			}{
				{
					name: "GET /",
					request: Request{
						Method: "GET",
						Path:   "/",
					},
					registerMock: func() {
						mockTransport.RegisterResponder("GET", "/", mockTransport.stringResponder(200, "OK"))
					},
				},
				{
					name: "GET /{slug}",
					request: Request{
						Method: "GET",
						Path:   "/{slug}",
						PathParams: map[string]string{
							"slug": "test",
						},
					},
					registerMock: func() {
						mockTransport.RegisterResponder("GET", "/test", mockTransport.stringResponder(200, "OK"))
					},
				},
				{
					name: "GET /?a=1&b=2",
					request: Request{
						Method: "GET",
						Path:   "/",
						Queries: url.Values{
							"a": []string{"1"},
							"b": []string{"2"},
						},
					},
					registerMock: func() {
						mockTransport.RegisterResponder("GET", "/?a=1&b=2", mockTransport.stringResponder(200, "OK"))
					},
				},
				{
					name: "GET /header",
					request: Request{
						Path: "/header",
						Header: http.Header{
							"User-Agent": []string{"genapi/test_suite"},
						},
					},
					registerMock: func() {
						mockTransport.RegisterResponder("GET", "/header", func(req *http.Request) (*http.Response, error) {
							if req.Header.Get("User-Agent") != "genapi/test_suite" {
								return mockTransport.stringResponse(400, "Bad Request"), nil
							}
							return mockTransport.stringResponse(200, "OK"), nil
						})
					},
				},
				{
					name: "POST /body",
					request: Request{
						Method: "POST",
						Path:   "/body",
						Body:   map[string]any{"name": "John"},
					},
					registerMock: func() {
						mockTransport.RegisterResponder("POST", "/body", func(req *http.Request) (*http.Response, error) {
							var body map[string]any
							err := json.NewDecoder(req.Body).Decode(&body)
							if err != nil {
								return mockTransport.stringResponse(500, "Internal Server Error"), err
							}
							if body["name"] != "John" {
								return mockTransport.stringResponse(400, "Bad Request"), nil
							}
							return mockTransport.stringResponse(200, "OK"), nil
						})
					},
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					client := setupMockClient(t)
					tc.registerMock()

					resp, err := client.Do(&tc.request)

					RequireNoError(t, err)
					AssertEqual(t, http.StatusOK, resp.StatusCode)
				})
			}

			t.Run("Context", func(t *testing.T) {
				t.Run("Cancellation", func(t *testing.T) {
					client := setupMockClient(t)
					done := make(chan struct{})

					mockTransport.RegisterResponder("GET", "/", func(req *http.Request) (*http.Response, error) {
						select {
						case <-req.Context().Done():
							return nil, req.Context().Err()
						case <-done:
							return mockTransport.stringResponse(200, "OK"), nil
						}
					})

					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						time.Sleep(10 * time.Millisecond)
						cancel()
					}()

					_, err := client.Do(&Request{Method: "GET", Path: "/", Context: ctx})
					close(done)

					RequireErrorIs(t, err, context.Canceled)
				})
			})
		})

		t.Run("HandleResponse", func(t *testing.T) {
			client := setupMockClient(t)

			mockTransport.RegisterResponder("GET", "/", mockTransport.stringResponder(200, `{"name": "John"}`))

			resp, err := client.Do(&Request{Method: "GET", Path: "/"})
			RequireNoError(t, err)

			var result map[string]string
			err = json.NewDecoder(resp.Body).Decode(&result)
			RequireNoError(t, err)
			AssertEqual(t, "John", result["name"])
		})
	})

	t.Run("MockServer", func(t *testing.T) {
		testCases := []struct {
			name       string
			request    Request
			mockServer func(t *testing.T) http.HandlerFunc
		}{
			{
				name: "GET /",
				request: Request{
					Method: "GET",
					Path:   "/",
				},
				mockServer: func(t *testing.T) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {}
				},
			},
			{
				name: "GET /{slug}",
				request: Request{
					Method: "GET",
					Path:   "/{slug}",
					PathParams: map[string]string{
						"slug": "test",
					},
				},
				mockServer: func(t *testing.T) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						AssertEqual(t, "/test", r.URL.Path)
					}
				},
			},
			{
				name: "GET /?a=1&b=2",
				request: Request{
					Method: "GET",
					Path:   "/",
					Queries: url.Values{
						"a": []string{"1"},
						"b": []string{"2"},
					},
				},
				mockServer: func(t *testing.T) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						AssertEqual(t, "/?a=1&b=2", r.URL.String())
					}
				},
			},
			{
				name: "GET /header",
				request: Request{
					Path: "/header",
					Header: http.Header{
						"User-Agent": []string{"genapi/test_suite"},
					},
				},
				mockServer: func(t *testing.T) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						AssertEqual(t, "genapi/test_suite", r.Header.Get("User-Agent"))
					}
				},
			},
			{
				name: "POST /body",
				request: Request{
					Method: "POST",
					Path:   "/body",
					Body:   map[string]any{"name": "John"},
				},
				mockServer: func(t *testing.T) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						var body map[string]any
						err := json.NewDecoder(r.Body).Decode(&body)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
						if body["name"] != "John" {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
					}
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				server := httptest.NewServer(tc.mockServer(t))
				defer server.Close()

				client := createClient()
				client.SetConfig(Config{BaseURL: server.URL})

				resp, err := client.Do(&tc.request)

				RequireNoError(t, err)
				AssertEqual(t, http.StatusOK, resp.StatusCode)
			})
		}
	})
}

func RequireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatal(err)
	}
}

func RequireErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("expected error %v, got %v", target, err)
	}
}

func AssertEqual(t *testing.T, expected, actual any) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func AssertNotEqual(t *testing.T, expected, actual any) {
	t.Helper()

	if reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
