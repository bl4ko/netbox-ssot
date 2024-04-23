package fmc

import "fmt"

func (fmcs *FMCSource) initDevices(c *fmcClient) error {
	fmcs.Domains = make(map[string]*Domain)
	domains, err := c.GetDomains()
	if err != nil {
		return fmt.Errorf("get domains: %s", err)
	}

	fmcs.Devices = make(map[string]*DeviceInfo)
	fmcs.DevicePhysicalIfaces = make(map[string][]*PhysicalInterfaceInfo)
	fmcs.DeviceVlanIfaces = make(map[string][]*VLANInterfaceInfo)
	for _, domain := range domains {
		devices, err := c.GetDevices(domain.UUID)
		if err != nil {
			return fmt.Errorf("get devices: %s", err)
		}
		for _, device := range devices {
			deviceInfo, err := c.GetDeviceInfo(domain.UUID, device.ID)
			if err != nil {
				return fmt.Errorf("error extracting device info: %s", err)
			}
			fmcs.Devices[device.ID] = deviceInfo

			fmcs.DevicePhysicalIfaces[device.ID] = make([]*PhysicalInterfaceInfo, 0)
			pIfaces, err := c.GetDevicePhysicalInterfaces(domain.UUID, device.ID)
			if err != nil {
				return fmt.Errorf("error getting physical interfaces: %s", err)
			}
			for _, pInterface := range pIfaces {
				pIfaceInfo, err := c.GetPhysicalInterfaceInfo(domain.UUID, device.ID, pInterface.ID)
				if err != nil {
					return fmt.Errorf("get physical interface info: %s", err)
				}
				fmcs.DevicePhysicalIfaces[device.ID] = append(fmcs.DevicePhysicalIfaces[device.ID], pIfaceInfo)
			}

			fmcs.DeviceVlanIfaces[device.ID] = make([]*VLANInterfaceInfo, 0)
			vlanIfaces, err := c.GetDeviceVLANInterfaces(domain.UUID, device.ID)
			if err != nil {
				return fmt.Errorf("error getting vlan interfaces")
			}
			for _, vlanIface := range vlanIfaces {
				vlanIfaceInfo, err := c.GetVLANInterfaceInfo(domain.UUID, device.ID, vlanIface.ID)
				if err != nil {
					return fmt.Errorf("error vlan interface info")
				}
				fmcs.DeviceVlanIfaces[device.ID] = append(fmcs.DeviceVlanIfaces[device.ID], vlanIfaceInfo)
			}
		}
	}
	return nil
}
