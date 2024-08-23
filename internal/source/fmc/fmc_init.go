package fmc

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/source/fmc/client"
)

// Init initializes the FMC source.
func (fmcs *FMCSource) initObjects(c *client.FMCClient) error {
	domains, err := fmcs.initDomains(c)
	if err != nil {
		return fmt.Errorf("init domains: %s", err)
	}

	for _, domain := range domains {
		if err := fmcs.initDevices(c, domain); err != nil {
			return fmt.Errorf("init devices: %s", err)
		}
	}
	return nil
}

func (fmcs *FMCSource) initDomains(c *client.FMCClient) ([]client.Domain, error) {
	fmcs.Logger.Debug(fmcs.Ctx, "Getting domains from fmc...")
	domains, err := c.GetDomains()
	if err != nil {
		return nil, fmt.Errorf("get domains: %s", err)
	}
	for _, domain := range domains {
		fmcs.Domains[domain.UUID] = domain
	}
	fmcs.Logger.Debugf(fmcs.Ctx, "Received domains %v", domains)
	return domains, nil
}

func (fmcs *FMCSource) initDevices(c *client.FMCClient, domain client.Domain) error {
	fmcs.Logger.Debugf(fmcs.Ctx, "Getting devices for %s domain...", domain.Name)
	devices, err := c.GetDevices(domain.UUID)
	if err != nil {
		return fmt.Errorf("get devices: %s", err)
	}
	fmcs.Logger.Debugf(fmcs.Ctx, "Received devices %v", devices)

	for _, device := range devices {
		deviceInfo, err := c.GetDeviceInfo(domain.UUID, device.ID)
		if err != nil {
			return fmt.Errorf("error extracting device info: %s", err)
		}
		fmcs.Devices[device.ID] = deviceInfo

		// Initialize Device physical interfaces
		fmcs.Logger.Debugf(fmcs.Ctx, "Getting physical interfaces for device %s", deviceInfo.Name)
		err = fmcs.initDevicePhysicalInterfaces(c, domain, device)
		if err != nil {
			return fmt.Errorf("error initializing physical interfaces: %s", err)
		}

		// Initialize device VLAN interfaces
		fmcs.Logger.Debugf(fmcs.Ctx, "Getting vlan interfaces for device %s", deviceInfo.Name)
		err = fmcs.initDeviceVLANInterfaces(c, domain, device)
		if err != nil {
			return fmt.Errorf("error initializing vlan interfaces: %s", err)
		}
	}
	return nil
}

func (fmcs *FMCSource) initDevicePhysicalInterfaces(c *client.FMCClient, domain client.Domain, device client.Device) error {
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
	return nil
}

func (fmcs *FMCSource) initDeviceVLANInterfaces(c *client.FMCClient, domain client.Domain, device client.Device) error {
	vlanIfaces, err := c.GetDeviceVLANInterfaces(domain.UUID, device.ID)
	if err != nil {
		return fmt.Errorf("error getting vlan interfaces: %s", err)
	}
	for _, vlanIface := range vlanIfaces {
		vlanIfaceInfo, err := c.GetVLANInterfaceInfo(domain.UUID, device.ID, vlanIface.ID)
		if err != nil {
			return fmt.Errorf("error vlan interface info: %s", err)
		}
		fmcs.DeviceVlanIfaces[device.ID] = append(fmcs.DeviceVlanIfaces[device.ID], vlanIfaceInfo)
	}
	return nil
}
