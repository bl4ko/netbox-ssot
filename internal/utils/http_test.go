package utils

import (
	"crypto/tls"
	"net/http"
	"reflect"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	wantInsecureClient := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	gotInsecureClient, err := NewHTTPClient(false, "")
	if err != nil {
		t.Errorf("not expecting error, but got: %s", err)
	}
	if !reflect.DeepEqual(wantInsecureClient, gotInsecureClient) {
		t.Errorf("want: %v, got: %v", wantInsecureClient, gotInsecureClient)
	}

	// wrong path
	_, err = NewHTTPClient(true, "\\//")
	if err == nil {
		t.Error("want error but got none")
	}

	customCertPool, _ := LoadExtraCert("")
	wantHTTPClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: customCertPool,
			},
		},
	}
	gotHTTPClient, err := NewHTTPClient(true, "")
	if err != nil {
		t.Errorf("don't want error but got: %s", err)
	}

	if !reflect.DeepEqual(wantHTTPClient, gotHTTPClient) {
		t.Errorf("want: %v, got %v", wantHTTPClient, gotHTTPClient)
	}
}
