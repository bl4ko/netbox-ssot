package f5

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIResponse[T any] struct {
	Kind  string `json:"kind"`
	Items T      `json:"items"`
}

type VirtualServerResponse struct {
	Name        string `json:"name"`
	FullPath    string `json:"fullPath"`
	Partition   string `json:"partition"`
	Destination string `json:"destination"`
	Description string `json:"description"`
	Pool        string `json:"pool"`
	IPProtocol  string `json:"ipProtocol"`
	Source      string `json:"source"`
	Mask        string `json:"mask"`
	Enabled     bool   `json:"enabled"`
	Disabled    bool   `json:"disabled"`
}

func (fs *F5Source) initVirtualServers(ctx context.Context, c *Client) error {
	res, err := c.MakeRequest(ctx, http.MethodGet, "ltm/virtual", nil)
	if err != nil {
		return fmt.Errorf("request error: %s", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("body read error: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("got http status: %d, body: %s", res.StatusCode, string(body))
	}

	var response APIResponse[[]VirtualServerResponse]
	if err = json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("body unmarshal error: %s", err)
	}

	fs.VirtualServers = response.Items
	fs.Logger.Debugf(fs.Ctx, "fetched %d virtual servers from F5 BIG-IP", len(fs.VirtualServers))
	return nil
}
