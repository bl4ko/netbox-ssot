package dnac

import (
	"context"
	"fmt"
	"net/http"

	dnac "github.com/cisco-en-programmability/dnacenter-go-sdk/v5/sdk"
)

func (ds *DnacSource) InitSites(ctx context.Context, c *dnac.Client) error {
	sites, response, err := c.Sites.GetSite(nil)
	if err != nil {
		return fmt.Errorf("init sites: %s", err)
	}
	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("init sites response code: %s", response.String())
	}
	for _, site := range *sites.Response {
		ds.Sites[site.ID] = site
	}
	return nil
}

func (ds *DnacSource) InitDevices(ctx context.Context, c *dnac.Client) error {
	devices, response, err := c.Devices.Devices(nil)
	if err != nil {
		return fmt.Errorf("init devices: %s", err)
	}
	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("init devices response code: %s", response.String())
	}
	for _, device := range *devices.Response {
		fmt.Printf("%T %v\n", device, device)
	}
	return nil
}
