package it

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func newServer(t *testing.T) http.HandlerFunc {
	// Define response structure
	type echoResponse struct {
		Path    string      `json:"path"`
		Body    any         `json:"body"`
		Headers http.Header `json:"headers"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Prepare response
		resp := echoResponse{
			Path:    r.URL.Path + "?" + r.URL.Query().Encode(),
			Headers: r.Header,
		}

		// Read and set body if present
		if r.Body != nil && r.ContentLength > 0 {
			defer r.Body.Close()
			var bodyData map[string]any
			if err := json.NewDecoder(r.Body).Decode(&bodyData); err != nil {
				http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
				return
			}
			resp.Body = bodyData
		}

		// Marshal the response body
		responseBody, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Failed to marshal response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))

		// Copy request headers, excluding Content-Length
		for key, values := range r.Header {
			if key == "Content-Length" {
				continue
			}
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Set status code from query parameter or default to 200
		if status := r.URL.Query().Get("status_code"); status != "" {
			code, err := strconv.Atoi(status)
			if err != nil {
				http.Error(w, "Invalid status code: "+err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(code)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		// Write response body
		if _, err := w.Write(responseBody); err != nil {
			t.Errorf("Failed to write response: %v", err)
			return
		}
	}
}
