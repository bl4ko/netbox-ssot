package utils

import (
	"net/http"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	_, err := NewHTTPClient(false, "")
	if err != nil {
		t.Errorf("not expecting error, but got: %s", err)
	}

	// wrong path
	_, err = NewHTTPClient(true, "\\//")
	if err == nil {
		t.Error("expected error but got none")
	}

	// Check if `InsecureSkipVerify` is set correctly
	insecureClient, err := NewHTTPClient(false, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	transport := insecureClient.Transport.(*http.Transport) //nolint:forcetypeassert
	if transport.TLSClientConfig.InsecureSkipVerify != true {
		t.Errorf("expected InsecureSkipVerify to be true, got false")
	}

	// Check if RootCAs is set when expected
	certClient, err := NewHTTPClient(true, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	transport = certClient.Transport.(*http.Transport) //nolint:forcetypeassert
	if transport.TLSClientConfig.RootCAs == nil {
		t.Errorf("expected RootCAs to be set, got nil")
	}
}
