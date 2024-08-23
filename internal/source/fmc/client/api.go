package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	maxRetries     = 5
	initialBackoff = 500 * time.Millisecond
	backoffFactor  = 2.0
	maxBackoff     = 16 * time.Second
)

// exponentialBackoff calculates the backoff duration based on the number of attempts.
func exponentialBackoff(attempt int) time.Duration {
	backoff := time.Duration(float64(initialBackoff) * math.Pow(backoffFactor, float64(attempt)))
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	return backoff
}

// Authenticate performs authentication on FMC API. If successful it returns access and refresh tokens.
func (fmcc FMCClient) Authenticate() (string, string, error) {
	var (
		accessToken  string
		refreshToken string
		err          error
	)

	for attempt := 0; attempt < maxRetries; attempt++ {
		accessToken, refreshToken, err = fmcc.authenticateOnce()
		if err == nil {
			return accessToken, refreshToken, nil
		}

		fmcc.Logger.Debugf(fmcc.Ctx, "authentication attempt %d failed: %s", attempt, err)
		time.Sleep(exponentialBackoff(attempt))
	}

	return "", "", fmt.Errorf("authentication failed after %d attempts: %w", maxRetries, err)
}

// Helper function to Authenticate. Performs single attempt to authenticate to fmc api.
func (fmcc FMCClient) authenticateOnce() (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), fmcc.DefaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/fmc_platform/v1/auth/generatetoken", fmcc.BaseURL), nil)
	if err != nil {
		return "", "", fmt.Errorf("new request with context: %w", err)
	}

	// Add Basic authentication header
	auth := fmt.Sprintf("%s:%s", fmcc.Username, fmcc.Password)
	auth = base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))

	res, err := fmcc.HTTPClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("req err: %w", err)
	}
	defer res.Body.Close() // Close the response body

	// Extract access and refresh tokens from response
	accessToken := res.Header.Get("X-auth-access-token")
	refreshToken := res.Header.Get("X-auth-refresh-token")
	if accessToken == "" || refreshToken == "" {
		return "", "", fmt.Errorf("failed extracting access and refresh tokens from response") //nolint:goerr113
	}
	return accessToken, refreshToken, nil
}

// MakeRequest sends an HTTP request to the specified path using the given method and body.
// It retries the request with exponential backoff up to a maximum number of attempts.
// If the request fails after the maximum number of attempts, it returns an error.
func (fmcc *FMCClient) MakeRequest(ctx context.Context, method, path string, body io.Reader, result interface{}) error {
	var (
		resp           *http.Response
		err            error
		tokenRefreshed bool
	)

	for attempt := 0; attempt < maxRetries; attempt++ {
		if ctx.Err() != nil {
			fmcc.Logger.Debugf(ctx, "context canceled or expired: %s", ctx.Err())
			return ctx.Err()
		}

		resp, err = fmcc.makeRequestOnce(ctx, method, path, body)
		if err != nil {
			fmcc.Logger.Debugf(ctx, "request attempt %d failed: %s", attempt, err)
			time.Sleep(exponentialBackoff(attempt))
			continue
		}

		// Check if the status code is 401 Unauthorized
		if resp.StatusCode == http.StatusUnauthorized {
			if !tokenRefreshed {
				fmcc.Logger.Debugf(ctx, "received 401 Unauthorized, attempting to refresh token")

				accessToken, refreshToken, authErr := fmcc.Authenticate()
				if authErr != nil {
					return fmt.Errorf("failed to refresh token: %w", authErr)
				}

				// Update the FMCClient with the new tokens.
				fmcc.AccessToken = accessToken
				fmcc.RefreshToken = refreshToken

				tokenRefreshed = true // Mark that the token has been refreshed.
				continue              // Retry the request immediately after refreshing the token.
			}
			// If the token has already been refreshed, return the 401 error.
			return fmt.Errorf("request failed with 401 Unauthorized after token refresh")
		}

		// Process the response if it's not a 401
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		err = json.Unmarshal(bodyBytes, result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		return nil
	}

	return fmt.Errorf("request failed after %d attempts: %w", maxRetries, err)
}

// makeRequestOnce sends an HTTP request to the specified path using the given method and body.
// It is a helper function for MakeRequest that sends the request only once.
func (fmcc *FMCClient) makeRequestOnce(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, fmcc.DefaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctxWithTimeout, method, fmt.Sprintf("%s/%s", fmcc.BaseURL, path), body)
	if err != nil {
		return nil, err
	}
	// Set the Authorization header.
	req.Header.Set("X-auth-access-token", fmcc.AccessToken)
	return fmcc.HTTPClient.Do(req)
}

// GetDomains returns a list of domains from the FMC API.
// It sends a GET request to the /fmc_platform/v1/info/domain endpoint.
func (fmcc *FMCClient) GetDomains() ([]Domain, error) {
	offset := 0
	limit := 25
	domains := []Domain{}
	ctx := context.Background()

	for {
		var marshaledResponse APIResponse[Domain]
		err := fmcc.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("fmc_platform/v1/info/domain?offset=%d&limit=%d", offset, limit), nil, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("make request for domains: %w", err)
		}

		if len(marshaledResponse.Items) > 0 {
			domains = append(domains, marshaledResponse.Items...)
		}

		if len(marshaledResponse.Items) < limit {
			break
		}
		offset += limit
	}

	return domains, nil
}

// GetDevices returns a list of devices from the FMC API for the specified domain.
func (fmcc *FMCClient) GetDevices(domainUUID string) ([]Device, error) {
	offset := 0
	limit := 25
	devices := []Device{}
	ctx := context.Background()

	for {
		devicesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords?offset=%d&limit=%d", domainUUID, offset, limit)
		var marshaledResponse APIResponse[Device]
		err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("make request for devices: %w", err)
		}

		if len(marshaledResponse.Items) > 0 {
			devices = append(devices, marshaledResponse.Items...)
		}

		if len(marshaledResponse.Items) < limit {
			break
		}
		offset += limit
	}

	return devices, nil
}

// GetDevicePhysicalInterfaces returns a list of physical interfaces for the specified device in the specified domain.
func (fmcc *FMCClient) GetDevicePhysicalInterfaces(domainUUID string, deviceID string) ([]PhysicalInterface, error) {
	offset := 0
	limit := 25
	pIfaces := []PhysicalInterface{}
	ctx := context.Background()

	for {
		pInterfacesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s/physicalinterfaces?offset=%d&limit=%d", domainUUID, deviceID, offset, limit)
		var marshaledResponse APIResponse[PhysicalInterface]
		err := fmcc.MakeRequest(ctx, http.MethodGet, pInterfacesURL, nil, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("make request for physical interfaces: %w", err)
		}

		if len(marshaledResponse.Items) > 0 {
			pIfaces = append(pIfaces, marshaledResponse.Items...)
		}

		if len(marshaledResponse.Items) < limit {
			break
		}
		offset += limit
	}

	return pIfaces, nil
}

func (fmcc *FMCClient) GetDeviceVLANInterfaces(domainUUID string, deviceID string) ([]VlanInterface, error) {
	offset := 0
	limit := 25
	vlanIfaces := []VlanInterface{}
	ctx := context.Background()

	for {
		pInterfacesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s/vlaninterfaces?offset=%d&limit=%d", domainUUID, deviceID, offset, limit)
		var marshaledResponse APIResponse[VlanInterface]
		err := fmcc.MakeRequest(ctx, http.MethodGet, pInterfacesURL, nil, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("make request for VLAN interfaces: %w", err)
		}

		if len(marshaledResponse.Items) > 0 {
			vlanIfaces = append(vlanIfaces, marshaledResponse.Items...)
		}

		if len(marshaledResponse.Items) < limit {
			break
		}
		offset += limit
	}

	return vlanIfaces, nil
}

func (fmcc *FMCClient) GetPhysicalInterfaceInfo(domainUUID string, deviceID string, interfaceID string) (*PhysicalInterfaceInfo, error) {
	var pInterfaceInfo PhysicalInterfaceInfo
	ctx := context.Background()

	devicesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s/physicalinterfaces/%s", domainUUID, deviceID, interfaceID)
	err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &pInterfaceInfo)
	if err != nil {
		return nil, fmt.Errorf("make request for physical interface info: %w", err)
	}

	return &pInterfaceInfo, nil
}

func (fmcc *FMCClient) GetVLANInterfaceInfo(domainUUID string, deviceID string, interfaceID string) (*VLANInterfaceInfo, error) {
	var vlanInterfaceInfo VLANInterfaceInfo
	ctx := context.Background()

	devicesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s/vlaninterfaces/%s", domainUUID, deviceID, interfaceID)
	err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &vlanInterfaceInfo)
	if err != nil {
		return nil, fmt.Errorf("make request for VLAN interface info: %w", err)
	}

	return &vlanInterfaceInfo, nil
}

func (fmcc *FMCClient) GetDeviceInfo(domainUUID string, deviceID string) (*DeviceInfo, error) {
	var deviceInfo DeviceInfo
	ctx := context.Background()

	devicesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s", domainUUID, deviceID)
	err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &deviceInfo)
	if err != nil {
		return nil, fmt.Errorf("make request for device info: %w", err)
	}

	return &deviceInfo, nil
}
