package utils

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

// NewHTTPClient creates an http client with tls config depending on validateCert
// and caFile parameter.
func NewHTTPClient(validateCert bool, caFile string) (*http.Client, error) {
	httpClient := &http.Client{}
	if validateCert {
		customCertPool, err := LoadExtraCert(caFile)
		if err != nil {
			return nil, fmt.Errorf("load extra cert: %s", err)
		}
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: customCertPool,
			},
		}
	} else {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	return httpClient, nil
}
