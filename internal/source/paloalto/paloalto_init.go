package paloalto

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/eth"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer3"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"
	"github.com/PaloAltoNetworks/pango/netw/zone"
	"github.com/PaloAltoNetworks/pango/vsys"
)

// Init system info collects system info from paloalto.
func (pas *PaloAltoSource) InitSystemInfo(c *pango.Firewall) error {
	pas.SystemInfo = c.Client.SystemInfo
	return nil
}

func (pas *PaloAltoSource) InitVirtualSystems(c *pango.Firewall) error {
	virtualSystems, err := c.Vsys.GetAll()
	if err != nil {
		return fmt.Errorf("get all virtual systems: %s", err)
	}
	pas.VirtualSystems = make(map[string]vsys.Entry)
	pas.SecurityZones = make(map[string]zone.Entry)
	pas.Iface2SecurityZone = make(map[string]string)
	for _, virtualSystem := range virtualSystems {
		securityZones, err := c.Network.Zone.GetAll(virtualSystem.Name)
		if err != nil {
			return fmt.Errorf("get zones for virtual system %s: %s", virtualSystem.Name, err)
		}
		pas.VirtualSystems[virtualSystem.Name] = virtualSystem
		for _, securityZone := range securityZones {
			pas.SecurityZones[securityZone.Name] = securityZone
			for _, iface := range securityZone.Interfaces {
				pas.Iface2SecurityZone[iface] = securityZone.Name
			}
		}
	}
	return nil
}

func (pas *PaloAltoSource) InitVirtualRouters(c *pango.Firewall) error {
	routers, err := c.Network.VirtualRouter.GetAll()
	if err != nil {
		return err
	}
	pas.VirtualRouters = make(map[string]router.Entry)
	pas.Iface2VirtualRouter = make(map[string]string)
	for _, router := range routers {
		pas.VirtualRouters[router.Name] = router
		for _, routerInterface := range router.Interfaces {
			pas.Iface2VirtualRouter[routerInterface] = router.Name
		}
	}
	return nil
}

func (pas *PaloAltoSource) InitInterfaces(c *pango.Firewall) error {
	ethInterfaces, err := c.Network.EthernetInterface.GetAll()
	if err != nil {
		return err
	}
	pas.Ifaces = make(map[string]eth.Entry)
	pas.Iface2SubIfaces = make(map[string][]layer3.Entry)
	for _, ethInterface := range ethInterfaces {
		pas.Ifaces[ethInterface.Name] = ethInterface
		subInterfaces, err := c.Network.Layer3Subinterface.GetAll(layer3.EthernetInterface, ethInterface.Name)
		if err != nil {
			return fmt.Errorf("layer 3 subinterfaces: %s", err)
		}
		pas.Iface2SubIfaces[ethInterface.Name] = make([]layer3.Entry, 0, len(subInterfaces))
		pas.Iface2SubIfaces[ethInterface.Name] = subInterfaces
	}
	return nil
}
