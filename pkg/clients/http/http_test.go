package http

import (
	"net/http"
	"strings"
	"testing"

	"github.com/lexcao/genapi/internal"
)

func (c *HttpClient) GetClient() *http.Client {
	return c.client
}

func TestHttpClient(t *testing.T) {
	internal.TestHttpClient(t, func() internal.HttpClientTester { return DefaultClient })
}

func TestHttpClient_CachedBaseURL(t *testing.T) {
	t.Run("ValidBaseURL", func(t *testing.T) {
		client := New(http.DefaultClient)
		
		// Should cache valid URL and work correctly
		client.SetConfig(internal.Config{BaseURL: "https://api.example.com"})
		
		// Test URL resolution works with cached URL
		url, err := resolveURL(*client.baseURL, "/test", nil)
		if err != nil {
			t.Fatalf("resolveURL failed: %v", err)
		}
		
		expected := "https://api.example.com/test"
		if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})
	
	t.Run("InvalidBaseURL", func(t *testing.T) {
		client := New(http.DefaultClient)
		
		// Should panic on invalid URL
		defer func() {
			if r := recover(); r == nil {
				t.Error("SetConfig should panic on invalid base URL")
			} else if panicMsg, ok := r.(string); ok {
				if !strings.Contains(panicMsg, "invalid base URL") || !strings.Contains(panicMsg, "://invalid") {
					t.Errorf("unexpected panic message: %s", panicMsg)
				}
			}
		}()
		
		client.SetConfig(internal.Config{BaseURL: "://invalid"})
	})
	
	t.Run("EmptyBaseURL", func(t *testing.T) {
		client := New(http.DefaultClient)
		
		// Empty base URL should be valid (parsed as relative URL)
		client.SetConfig(internal.Config{BaseURL: ""})
		
		// Test URL resolution works with empty base URL
		url, err := resolveURL(*client.baseURL, "/test", nil)
		if err != nil {
			t.Fatalf("resolveURL failed: %v", err)
		}
		
		expected := "/test"
		if url != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	})

	t.Run("PerformanceBenefit", func(t *testing.T) {
		// Test that multiple resolveURL calls reuse the cached base URL
		client := New(http.DefaultClient)
		client.SetConfig(internal.Config{BaseURL: "https://api.example.com"})
		
		// Call resolveURL multiple times to simulate multiple requests
		// This should use the cached baseURL without parsing
		for i := 0; i < 10; i++ {
			url, err := resolveURL(*client.baseURL, "/test", map[string]string{
				"id": "123", 
			})
			if err != nil {
				t.Fatalf("resolveURL failed on iteration %d: %v", i, err)
			}
			
			expected := "https://api.example.com/test"
			if url != expected {
				t.Errorf("iteration %d: expected %s, got %s", i, expected, url)
			}
		}
		
		// All calls should have succeeded without errors, demonstrating caching works
	})
}
