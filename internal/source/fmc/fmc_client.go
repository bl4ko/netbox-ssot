package fmc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

type fmcClient struct {
	HTTPClient   *http.Client
	BaseURL      string
	Username     string
	Password     string
	AccessToken  string
	RefreshToken string
	Timeout      time.Duration
}

func newFMCClient(username string, password string, httpScheme string, hostname string, port int, httpClient *http.Client) (*fmcClient, error) {
	// First we obtain access and refresh token
	c := &fmcClient{
		HTTPClient: httpClient,
		BaseURL:    fmt.Sprintf("%s://%s:%d/api", httpScheme, hostname, port),
		Username:   username,
		Password:   password,
		Timeout:    time.Second * constants.DefaultAPITimeout,
	}

	aToken, rToken, err := c.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("authentication: %w", err)
	}

	c.AccessToken = aToken
	c.RefreshToken = rToken

	return c, nil
}

// Authenticate performs authentication on FMC API. If successful it returns access and refresh tokens.
func (fmcc fmcClient) Authenticate() (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), fmcc.Timeout)
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

type PagingResponse struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
	Pages  int `json:"pages"`
}

type LinksResponse struct {
	Self string `json:"self"`
}

type APIResponse[T any] struct {
	Links  LinksResponse  `json:"links"`
	Paging PagingResponse `json:"paging"`
	Items  []T            `json:"items"`
}

type Domain struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Device struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func (fmcc *fmcClient) MakeRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, fmcc.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", fmcc.BaseURL, path), body)
	if err != nil {
		return nil, err
	}
	// Set the Authorization header.
	req.Header.Set("X-auth-access-token", fmcc.AccessToken)
	return fmcc.HTTPClient.Do(req)
}

func (fmcc *fmcClient) GetDomains() ([]Domain, error) {
	offset := 0
	limit := 25
	domains := []Domain{}
	ctx := context.Background()
	for {
		apiResponse, err := fmcc.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("fmc_platform/v1/info/domain?offset=%d&limit=%d", offset, limit), nil)
		if err != nil {
			return nil, fmt.Errorf("make request for domains: %w", err)
		}
		defer apiResponse.Body.Close()
		if apiResponse.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("wrong status code: %d", apiResponse.StatusCode)
		}
		var marshaledResponse APIResponse[Domain]
		bodyBytes, err := io.ReadAll(apiResponse.Body)
		if err != nil {
			return nil, fmt.Errorf("response body readAll: %w", err)
		}
		err = json.Unmarshal(bodyBytes, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal response body: %w", err)
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

func (fmcc *fmcClient) GetDevices(domainUUID string) ([]Device, error) {
	offset := 0
	limit := 25
	devices := []Device{}
	ctx := context.Background()
	devicesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords?offset=%d&limit=%d", domainUUID, offset, limit)
	for {
		apiResponse, err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil)
		if err != nil {
			return nil, fmt.Errorf("make request for domains: %w", err)
		}
		defer apiResponse.Body.Close()
		if apiResponse.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("wrong status code: %d", apiResponse.StatusCode)
		}
		var marshaledResponse APIResponse[Device]
		bodyBytes, err := io.ReadAll(apiResponse.Body)
		if err != nil {
			return nil, fmt.Errorf("response body readAll: %w", err)
		}
		err = json.Unmarshal(bodyBytes, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal response body: %w", err)
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

type PhysicalInterface struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func (fmcc *fmcClient) GetDevicePhysicalInterfaces(domainUUID string, deviceID string) ([]PhysicalInterface, error) {
	offset := 0
	limit := 25
	pIfaces := []PhysicalInterface{}
	ctx := context.Background()
	pInterfacesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s/physicalinterfaces?offset=%d&limit=%d", domainUUID, deviceID, offset, limit)
	for {
		apiResponse, err := fmcc.MakeRequest(ctx, http.MethodGet, pInterfacesURL, nil)
		if err != nil {
			return nil, fmt.Errorf("make request for domains: %w", err)
		}
		defer apiResponse.Body.Close()
		if apiResponse.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("wrong status code: %d", apiResponse.StatusCode)
		}
		var marshaledResponse APIResponse[PhysicalInterface]
		bodyBytes, err := io.ReadAll(apiResponse.Body)
		if err != nil {
			return nil, fmt.Errorf("response body readAll: %w", err)
		}
		err = json.Unmarshal(bodyBytes, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal response body: %w", err)
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

type VlanInterface struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func (fmcc *fmcClient) GetDeviceVLANInterfaces(domainUUID string, deviceID string) ([]VlanInterface, error) {
	offset := 0
	limit := 25
	vlanIfaces := []VlanInterface{}
	ctx := context.Background()
	pInterfacesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s/vlaninterfaces?offset=%d&limit=%d", domainUUID, deviceID, offset, limit)
	for {
		apiResponse, err := fmcc.MakeRequest(ctx, http.MethodGet, pInterfacesURL, nil)
		if err != nil {
			return nil, fmt.Errorf("make request for domains: %w", err)
		}
		defer apiResponse.Body.Close()
		if apiResponse.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("wrong status code: %d", apiResponse.StatusCode)
		}
		var marshaledResponse APIResponse[VlanInterface]
		bodyBytes, err := io.ReadAll(apiResponse.Body)
		if err != nil {
			return nil, fmt.Errorf("response body readAll: %w", err)
		}
		err = json.Unmarshal(bodyBytes, &marshaledResponse)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal response body: %w", err)
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

type PhysicalInterfaceInfo struct {
	Type        string `json:"type"`
	MTU         int    `json:"MTU"`
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Mode        string `json:"mode"`
	Description string `json:"description"`
	Hardware    *struct {
		Speed  string `json:"speed"`
		Duplex string `json:"duplex"`
	} `json:"hardware"`
	SecurityZone *struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"securityZone"`
	IPv4 *struct {
		Static *struct {
			Address string `json:"address"`
			Netmask string `json:"netmask"`
		} `json:"static"`
	} `json:"ipv4"`
	IPv6 *struct {
		EnableIPv6 bool `json:"enableIPV6"`
	} `json:"ipv6"`
}

func (fmcc *fmcClient) GetPhysicalInterfaceInfo(domainUUID string, deviceID string, interfaceID string) (*PhysicalInterfaceInfo, error) {
	var pInterfaceInfo PhysicalInterfaceInfo
	ctx := context.Background()
	devicesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s/physicalinterfaces/%s", domainUUID, deviceID, interfaceID)
	apiResponse, err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("make request for domains: %w", err)
	}
	defer apiResponse.Body.Close()
	if apiResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", apiResponse.StatusCode)
	}
	bodyBytes, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("response body readAll: %w", err)
	}
	err = json.Unmarshal(bodyBytes, &pInterfaceInfo)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal response body: %w", err)
	}

	return &pInterfaceInfo, nil
}

func (fmcc *fmcClient) GetVLANInterfaceInfo(domainUUID string, deviceID string, interfaceID string) (*VLANInterfaceInfo, error) {
	var vlanInterfaceInfo VLANInterfaceInfo
	ctx := context.Background()
	devicesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s/vlaninterfaces/%s", domainUUID, deviceID, interfaceID)
	apiResponse, err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("make request for domains: %w", err)
	}
	defer apiResponse.Body.Close()
	if apiResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", apiResponse.StatusCode)
	}
	bodyBytes, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("response body readAll: %w", err)
	}
	err = json.Unmarshal(bodyBytes, &vlanInterfaceInfo)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal response body: %w", err)
	}

	return &vlanInterfaceInfo, nil
}

type VLANInterfaceInfo struct {
	Type        string `json:"type"`
	Mode        string `json:"mode"`
	VID         int    `json:"vlanId"`
	MTU         int    `json:"MTU"`
	Enabled     bool   `json:"enabled"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Hardware    *struct {
		Speed  string `json:"speed"`
		Duplex string `json:"duplex"`
	} `json:"hardware"`
	SecurityZone *struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"securityZone"`
	IPv4 *struct {
		Static *struct {
			Address string `json:"address"`
			Netmask string `json:"netmask"`
		} `json:"static"`
	} `json:"ipv4"`
	IPv6 *struct {
		EnableIPv6 bool `json:"enableIPV6"`
	} `json:"ipv6"`
}

type DeviceInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Model       string `json:"model"`
	ModelID     string `json:"modelId"`
	ModelNumber string `json:"modelNumber"`
	SWVersion   string `json:"sw_version"`
	Hostname    string `json:"hostName"`
	Metadata    struct {
		SerialNumber  string `json:"deviceSerialNumber"`
		InventoryData struct {
			CPUCores   string `json:"cpuCores"`
			CPUType    string `json:"cpuType"`
			MemoryInMB string `json:"memoryInMB"`
		} `json:"inventoryData"`
	} `json:"metadata"`
}

func (fmcc *fmcClient) GetDeviceInfo(domainUUID string, deviceID string) (*DeviceInfo, error) {
	var deviceInfo DeviceInfo
	ctx := context.Background()
	devicesURL := fmt.Sprintf("fmc_config/v1/domain/%s/devices/devicerecords/%s", domainUUID, deviceID)
	apiResponse, err := fmcc.MakeRequest(ctx, http.MethodGet, devicesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("make request for domains: %w", err)
	}
	defer apiResponse.Body.Close()
	if apiResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", apiResponse.StatusCode)
	}
	bodyBytes, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("response body readAll: %w", err)
	}
	err = json.Unmarshal(bodyBytes, &deviceInfo)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal response body: %w", err)
	}

	return &deviceInfo, nil
}
