package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HttpClientTester interface {
	HttpClient
	GetClient() *http.Client
}

func TestHttpClient(t *testing.T, createClient func() HttpClientTester) {
	t.Run("MockClient", func(t *testing.T) {
		setupMockClient := func(t *testing.T) HttpClient {
			client := createClient()
			httpmock.ActivateNonDefault(client.GetClient())
			t.Cleanup(httpmock.DeactivateAndReset)
			return client
		}

		t.Run("SetConfig", func(t *testing.T) {
			t.Run("BaseURL", func(t *testing.T) {
				client := setupMockClient(t)
				client.SetConfig(Config{BaseURL: "https://lexcao.com/genapi"})

				httpmock.RegisterResponder("GET", "https://lexcao.com/genapi",
					httpmock.NewStringResponder(200, "OK"))

				resp, err := client.Do(&Request{Method: "GET"})

				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})

			t.Run("GlobalHeaders", func(t *testing.T) {
				client := setupMockClient(t)
				client.SetConfig(Config{
					Header: http.Header{
						"User-Agent": []string{"genapi/test_suite"},
					},
				})

				httpmock.RegisterResponder("GET", "/genapi",
					func(req *http.Request) (*http.Response, error) {
						assert.Equal(t, "genapi/test_suite", req.Header.Get("User-Agent"))
						return httpmock.NewStringResponse(200, "OK"), nil
					})

				resp, err := client.Do(&Request{Method: "GET", Path: "/genapi"})

				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
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
						httpmock.RegisterResponder("GET", "/", httpmock.NewStringResponder(200, "OK"))
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
						httpmock.RegisterResponder("GET", "/test", httpmock.NewStringResponder(200, "OK"))
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
						httpmock.RegisterResponder("GET", "/?a=1&b=2", httpmock.NewStringResponder(200, "OK"))
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
						httpmock.RegisterResponder("GET", "/header", func(req *http.Request) (*http.Response, error) {
							if req.Header.Get("User-Agent") != "genapi/test_suite" {
								return httpmock.NewStringResponse(400, "Bad Request"), nil
							}
							return httpmock.NewStringResponse(200, "OK"), nil
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
						httpmock.RegisterResponder("POST", "/body", func(req *http.Request) (*http.Response, error) {
							var body map[string]any
							err := json.NewDecoder(req.Body).Decode(&body)
							if err != nil {
								return httpmock.NewStringResponse(500, "Internal Server Error"), err
							}
							if body["name"] != "John" {
								return httpmock.NewStringResponse(400, "Bad Request"), nil
							}
							return httpmock.NewStringResponse(200, "OK"), nil
						})
					},
				},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					client := setupMockClient(t)
					tc.registerMock()

					resp, err := client.Do(&tc.request)

					require.NoErrorf(t, err, "Failed for request: %v", tc.request)
					assert.Equal(t, http.StatusOK, resp.StatusCode)
				})
			}

			t.Run("Context", func(t *testing.T) {
				t.Run("Cancellation", func(t *testing.T) {
					client := setupMockClient(t)
					done := make(chan struct{})

					httpmock.RegisterResponder("GET", "/", func(req *http.Request) (*http.Response, error) {
						select {
						case <-req.Context().Done():
							return nil, req.Context().Err()
						case <-done:
							return httpmock.NewStringResponse(200, "OK"), nil
						}
					})

					ctx, cancel := context.WithCancel(context.Background())
					go func() {
						time.Sleep(10 * time.Millisecond)
						cancel()
					}()

					resp, err := client.Do(&Request{Method: "GET", Path: "/", Context: ctx})
					close(done)

					require.ErrorIs(t, err, context.Canceled)
					require.Nil(t, resp)
				})
			})
		})

		t.Run("HandleResponse", func(t *testing.T) {
			client := setupMockClient(t)

			given := map[string]string{"name": "John"}
			httpmock.RegisterResponder("GET", "/", httpmock.NewJsonResponderOrPanic(200, given))

			resp, err := client.Do(&Request{Method: "GET", Path: "/"})
			require.NoError(t, err)

			var result map[string]string
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.Equal(t, "John", result["name"])
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
						assert.Equal(t, "/test", r.URL.Path)
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
						assert.Equal(t, "/?a=1&b=2", r.URL.String())
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
						assert.Equal(t, "genapi/test_suite", r.Header.Get("User-Agent"))
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

				require.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			})
		}
	})
}
