package fortigate

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAPIClient(t *testing.T) {
	httpClient := &http.Client{}
	c := NewAPIClient("testtoken", "https://example.com", httpClient)

	if c.APIToken != "testtoken" {
		t.Errorf("APIToken = %q, want %q", c.APIToken, "testtoken")
	}
	if c.BaseURL != "https://example.com" {
		t.Errorf("BaseURL = %q, want %q", c.BaseURL, "https://example.com")
	}
	if c.HTTPClient != httpClient {
		t.Error("HTTPClient was not set correctly")
	}
}

func TestMakeRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer testtoken" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer testtoken")
		}
		if r.Method != http.MethodGet {
			t.Errorf("Method = %q, want %q", r.Method, http.MethodGet)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	c := NewAPIClient("testtoken", server.URL, server.Client())
	resp, err := c.MakeRequest(context.Background(), http.MethodGet, "api/v2/test", nil)
	if err != nil {
		t.Fatalf("MakeRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"status":"ok"}` {
		t.Errorf("Body = %q, want %q", string(body), `{"status":"ok"}`)
	}
}

func TestMakeRequest_InvalidHost(t *testing.T) {
	c := NewAPIClient("token", "http://invalid-host-that-does-not-exist:99999", &http.Client{})
	resp, err := c.MakeRequest(context.Background(), http.MethodGet, "test", nil)
	if err == nil {
		resp.Body.Close()
		t.Error("expected error for invalid host, got nil")
	}
}
