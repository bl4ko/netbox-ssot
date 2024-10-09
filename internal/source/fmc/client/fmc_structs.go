package client

// VLANInterfaceInfo represents information about a VLAN interface.
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

// DeviceInfo represents information about a FMC device.
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

// VlanInterface represents a VLAN interface.
type VlanInterface struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// PhysicalInterface represents a physical interface.
type PhysicalInterface struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// PaginationResponse represents the paging information in the API response.
type PaginationResponse struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
	Pages  int `json:"pages"`
}

// LinksResponse represents the links in the API response.
type LinksResponse struct {
	Self string `json:"self"`
}

// APIResponse represents the API response.
type APIResponse[T any] struct {
	Links  LinksResponse      `json:"links"`
	Paging PaginationResponse `json:"paging"`
	Items  []T                `json:"items"`
}

// Domain represents a domain in FMC.
type Domain struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Device represents a device in FMC.
type Device struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// PhysicalInterfaceInfo represents information about a physical interface.
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