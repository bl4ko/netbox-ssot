package fortigate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIResponse[T any] struct {
	HTTPStatus int    `json:"http_status"`
	Serial     string `json:"serial"`
	Version    string `json:"version"`
	Results    T      `json:"results"`
}

type DeviceResponse struct {
	Hostname string `json:"hostname"`
}

type InterfaceResponse struct {
	Name        string        `json:"name"`
	Vdom        string        `json:"vdom"`
	IP          string        `json:"Ip"`
	Type        string        `json:"type"`
	Status      string        `json:"status"`
	Speed       string        `json:"speed"`
	Description string        `json:"description"`
	MTU         int           `json:"mtu"`
	MAC         string        `json:"macaddr"`
	VlanID      int           `json:"vlanid"`
	SecondaryIP []SecondaryIP `json:"secondaryip"`
	VRRPIP      []VRRPIP      `json:"vrrp"`
}

type SecondaryIP struct {
	IP string `json:"ip"`
}
type VRRPIP struct {
	VRIP string `json:"vrip"`
}

// Init system info collects system info from paloalto.
func (fs *FortigateSource) initSystemInfo(ctx context.Context, c *FortiClient) error {
	res, err := c.MakeRequest(ctx, http.MethodGet, "cmdb/system/global/", nil)
	if err != nil {
		return fmt.Errorf("request error: %s", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("body read error: %s", err)
	}

	var deviceResponse APIResponse[DeviceResponse]
	err = json.Unmarshal(body, &deviceResponse)
	if err != nil {
		return fmt.Errorf("body unmarshal error: %s", err)
	}

	if deviceResponse.HTTPStatus != http.StatusOK {
		return fmt.Errorf("got http status: %d", deviceResponse.HTTPStatus)
	}

	fs.SystemInfo = FortiSystemInfo{
		Hostname: deviceResponse.Results.Hostname,
		Version:  deviceResponse.Version,
		Serial:   deviceResponse.Serial,
	}

	return nil
}

// Fetches all information about interfaces from fortigate api.
func (fs *FortigateSource) initInterfaces(ctx context.Context, c *FortiClient) error {
	// Interfaces
	res, err := c.MakeRequest(ctx, http.MethodGet, "cmdb/system/interface/", nil)
	if err != nil {
		return fmt.Errorf("request error: %s", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("body read error: %s", err)
	}
	var interfaceResponse APIResponse[[]InterfaceResponse]
	err = json.Unmarshal(body, &interfaceResponse)
	if err != nil {
		return fmt.Errorf("body unamrshal error: %s", err)
	}

	if interfaceResponse.HTTPStatus != http.StatusOK {
		return fmt.Errorf("got http status: %d", interfaceResponse.HTTPStatus)
	}

	fs.Ifaces = make(map[string]InterfaceResponse, len(interfaceResponse.Results))
	for _, iface := range interfaceResponse.Results {
		fs.Ifaces[iface.Name] = iface
	}

	return nil
}
