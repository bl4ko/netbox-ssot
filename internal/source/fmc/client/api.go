package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
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
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/fmc_platform/v1/auth/generatetoken", fmcc.BaseURL),
		nil,
	)
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
		return "", "", fmt.Errorf(
			"failed extracting access and refresh tokens from response",
		) //nolint:goerr113
	}
	return accessToken, refreshToken, nil
}

// makeRequestOnce sends an HTTP request to the specified path using the given method and body.
// It reads the response body and returns it along with the response object.
func (fmcc *FMCClient) makeRequestOnce(
	ctx context.Context,
	method, path string,
	body io.Reader,
) (*http.Response, []byte, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, fmcc.DefaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctxWithTimeout,
		method,
		fmt.Sprintf("%s/%s", fmcc.BaseURL, path),
		body,
	)
	if err != nil {
		return nil, nil, err
	}

	// Set the Authorization header.
	req.Header.Set("X-auth-access-token", fmcc.AccessToken)

	resp, err := fmcc.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	// Ensure the response body is fully read before the context is canceled
	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close() // Close the response body immediately after reading
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, bodyBytes, nil
}

// MakeRequest sends an HTTP request to the specified path using the given method and body.
// It retries the request with exponential backoff up to a maximum number of attempts.
// If the request fails after the maximum number of attempts, it returns an error.
func (fmcc *FMCClient) MakeRequest(
	ctx context.Context,
	method, path string,
	body io.Reader,
	result interface{},
) error {
	var (
		resp           *http.Response
		bodyBytes      []byte
		err            error
		tokenRefreshed bool
	)

	for attempt := 0; attempt < maxRetries; attempt++ {
		if ctx.Err() != nil {
			fmcc.Logger.Debugf(ctx, "context canceled or expired: %s", ctx.Err())
			return ctx.Err()
		}

		// Recreate context on each attempt
		reqCtx, cancel := context.WithTimeout(ctx, fmcc.DefaultTimeout)
		defer cancel()

		fmcc.Logger.Debugf(
			fmcc.Ctx,
			"Making %s request to %s with body=%v (attempt=%d)",
			method,
			path,
			body,
			attempt,
		)

		resp, bodyBytes, err = fmcc.makeRequestOnce(reqCtx, method, path, body)
		if err != nil {
			if reqCtx.Err() == context.Canceled || reqCtx.Err() == context.DeadlineExceeded {
				fmcc.Logger.Debugf(
					fmcc.Ctx,
					"request attempt %d failed due to context timeout/cancellation: %s",
					attempt,
					err,
				)
				time.Sleep(exponentialBackoff(attempt))
				continue
			}
			fmcc.Logger.Debugf(fmcc.Ctx, "request attempt %d failed: %s", attempt, err)
			time.Sleep(exponentialBackoff(attempt))
			continue
		}

		// Check if the status code is 401 Unauthorized
		if resp.StatusCode == http.StatusUnauthorized {
			if !tokenRefreshed {
				fmcc.Logger.Debugf(
					fmcc.Ctx,
					"received 401 Unauthorized, attempting to refresh token",
				)

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
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			fmcc.Logger.Info(
				fmcc.Ctx,
				"Too many requests performed. FMC allows only 120 requests per miniute. Sleeping for a minute",
			)
			time.Sleep(1 * time.Minute)
			continue
		}

		// Process the response if it's not a 401
		if resp.StatusCode != http.StatusOK {
			// Extract and include the response body in the error message
			respBody := "<empty>"
			if len(bodyBytes) > 0 {
				respBody = string(bodyBytes)
			}
			return fmt.Errorf(
				"unexpected status code: %d, response body: %s",
				resp.StatusCode,
				respBody,
			)
		}

		err = json.Unmarshal(bodyBytes, result)
		if err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		return nil
	}

	return fmt.Errorf("request failed after %d attempts: %w", maxRetries, err)
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
		domainsURL := fmt.Sprintf("fmc_platform/v1/info/domain?offset=%d&limit=%d", offset, limit)
		err := fmcc.MakeRequest(ctx, http.MethodGet, domainsURL, nil, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("make request for domains (%s): %w", domainsURL, err)
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
		devicesURL := fmt.Sprintf(
			"fmc_config/v1/domain/%s/devices/devicerecords?offset=%d&limit=%d",
			domainUUID,
			offset,
			limit,
		)
		var marshaledResponse APIResponse[Device]
		err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("make request for devices (%s): %w", devicesURL, err)
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
func (fmcc *FMCClient) GetDevicePhysicalInterfaces(
	domainUUID string,
	deviceID string,
) ([]PhysicalInterface, error) {
	offset := 0
	limit := 25
	pIfaces := []PhysicalInterface{}
	ctx := context.Background()

	for {
		pInterfacesURL := fmt.Sprintf(
			"fmc_config/v1/domain/%s/devices/devicerecords/%s/physicalinterfaces?offset=%d&limit=%d",
			domainUUID,
			deviceID,
			offset,
			limit,
		)
		var marshaledResponse APIResponse[PhysicalInterface]
		err := fmcc.MakeRequest(ctx, http.MethodGet, pInterfacesURL, nil, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf(
				"make request for physical interfaces (%s): %w",
				pInterfacesURL,
				err,
			)
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

func (fmcc *FMCClient) GetDeviceVLANInterfaces(
	domainUUID string,
	deviceID string,
) ([]VlanInterface, error) {
	offset := 0
	limit := 25
	vlanIfaces := []VlanInterface{}
	ctx := context.Background()

	for {
		vInterfacesURL := fmt.Sprintf(
			"fmc_config/v1/domain/%s/devices/devicerecords/%s/vlaninterfaces?offset=%d&limit=%d",
			domainUUID,
			deviceID,
			offset,
			limit,
		)
		var marshaledResponse APIResponse[VlanInterface]
		err := fmcc.MakeRequest(ctx, http.MethodGet, vInterfacesURL, nil, &marshaledResponse)
		if err != nil {
			if strings.Contains(
				err.Error(),
				"VLAN Interface type is not supported on this device model",
			) {
				fmcc.Logger.Debugf(
					fmcc.Ctx,
					"VLAN Interface type is not supported on this device model",
				)
				return nil, nil
			}
			return nil, fmt.Errorf(
				"make request for VLAN interfaces with (%s): %w",
				vInterfacesURL,
				err,
			)
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

func (fmcc *FMCClient) GetDeviceEtherChannelInterfaces(
	domainUUID string,
	deviceID string,
) ([]EtherChannelInterface, error) {
	offset := 0
	limit := 25
	etherChannelIfaces := []EtherChannelInterface{}
	ctx := context.Background()

	for {
		vInterfacesURL := fmt.Sprintf(
			"fmc_config/v1/domain/%s/devices/devicerecords/%s/vlaninterfaces?offset=%d&limit=%d",
			domainUUID,
			deviceID,
			offset,
			limit,
		)
		var marshaledResponse APIResponse[EtherChannelInterface]
		err := fmcc.MakeRequest(ctx, http.MethodGet, vInterfacesURL, nil, &marshaledResponse)
		if err != nil {
			if strings.Contains(
				err.Error(),
				"VLAN Interface type is not supported on this device model",
			) {
				fmcc.Logger.Debugf(
					fmcc.Ctx,
					"VLAN Interface type is not supported on this device model",
				)
				return nil, nil
			}
			return nil, fmt.Errorf(
				"make request for VLAN interfaces with (%s): %w",
				vInterfacesURL,
				err,
			)
		}

		if len(marshaledResponse.Items) > 0 {
			etherChannelIfaces = append(etherChannelIfaces, marshaledResponse.Items...)
		}

		if len(marshaledResponse.Items) < limit {
			break
		}
		offset += limit
	}

	return etherChannelIfaces, nil
}

func (fmcc *FMCClient) GetDeviceSubInterfaces(
	domainUUID string,
	deviceID string,
) ([]SubInterface, error) {
	offset := 0
	limit := 25
	subIfaces := []SubInterface{}
	ctx := context.Background()

	for {
		subInterfacesURL := fmt.Sprintf(
			"fmc_config/v1/domain/%s/devices/devicerecords/%s/subinterfaces?offset=%d&limit=%d",
			domainUUID,
			deviceID,
			offset,
			limit,
		)
		var marshaledResponse APIResponse[SubInterface]
		err := fmcc.MakeRequest(ctx, http.MethodGet, subInterfacesURL, nil, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf(
				"make request for sub interfaces with (%s): %w",
				subInterfacesURL,
				err,
			)
		}

		if len(marshaledResponse.Items) > 0 {
			subIfaces = append(subIfaces, marshaledResponse.Items...)
		}

		if len(marshaledResponse.Items) < limit {
			break
		}
		offset += limit
	}
	return subIfaces, nil
}

func (fmcc *FMCClient) GetPhysicalInterfaceInfo(
	domainUUID string,
	deviceID string,
	interfaceID string,
) (*PhysicalInterfaceInfo, error) {
	var pInterfaceInfo PhysicalInterfaceInfo
	ctx := context.Background()

	devicesURL := fmt.Sprintf(
		"fmc_config/v1/domain/%s/devices/devicerecords/%s/physicalinterfaces/%s",
		domainUUID,
		deviceID,
		interfaceID,
	)
	err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &pInterfaceInfo)
	if err != nil {
		return nil, fmt.Errorf("make request for physical interface info (%s): %w", devicesURL, err)
	}

	return &pInterfaceInfo, nil
}

func (fmcc *FMCClient) GetVLANInterfaceInfo(
	domainUUID string,
	deviceID string,
	interfaceID string,
) (*VLANInterfaceInfo, error) {
	var vlanInterfaceInfo VLANInterfaceInfo
	ctx := context.Background()

	devicesURL := fmt.Sprintf(
		"fmc_config/v1/domain/%s/devices/devicerecords/%s/vlaninterfaces/%s",
		domainUUID,
		deviceID,
		interfaceID,
	)
	err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &vlanInterfaceInfo)
	if err != nil {
		return nil, fmt.Errorf(
			"make request for VLAN interface info with (%s): %w",
			devicesURL,
			err,
		)
	}

	return &vlanInterfaceInfo, nil
}

func (fmcc *FMCClient) GetEtherChannelInterfaceInfo(
	domainUUID string,
	deviceID string,
	interfaceID string,
) (*EtherChannelInterfaceInfo, error) {
	var etherChannelInterfaceInfo EtherChannelInterfaceInfo
	ctx := context.Background()

	devicesURL := fmt.Sprintf(
		"fmc_config/v1/domain/%s/devices/devicerecords/%s/etherchannelinterfaces/%s",
		domainUUID,
		deviceID,
		interfaceID,
	)
	err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &etherChannelInterfaceInfo)
	if err != nil {
		return nil, fmt.Errorf(
			"make request for EtherChannel interface info with (%s): %w",
			devicesURL,
			err,
		)
	}
	return &etherChannelInterfaceInfo, nil
}

func (fmcc *FMCClient) GetSubInterfaceInfo(
	domainUUID string,
	deviceID string,
	interfaceID string,
) (*SubInterfaceInfo, error) {
	var subInterfaceInfo SubInterfaceInfo
	ctx := context.Background()

	devicesURL := fmt.Sprintf(
		"fmc_config/v1/domain/%s/devices/devicerecords/%s/subinterfaces/%s",
		domainUUID,
		deviceID,
		interfaceID,
	)
	err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &subInterfaceInfo)
	if err != nil {
		return nil, fmt.Errorf("make request for sub interface info with (%s): %w", devicesURL, err)
	}
	return &subInterfaceInfo, nil
}

func (fmcc *FMCClient) GetDeviceInfo(domainUUID string, deviceID string) (*DeviceInfo, error) {
	var deviceInfo DeviceInfo
	ctx := context.Background()

	devicesURL := fmt.Sprintf(
		"fmc_config/v1/domain/%s/devices/devicerecords/%s",
		domainUUID,
		deviceID,
	)
	err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil, &deviceInfo)
	if err != nil {
		return nil, fmt.Errorf("make request for device info with (%s): %w", devicesURL, err)
	}

	return &deviceInfo, nil
}
