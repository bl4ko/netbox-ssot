package paloalto

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/eth"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer3"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"
)

func (pas *PaloAltoSource) InitDevices(c *pango.Firewall) error {
	routers, err := c.Network.VirtualRouter.GetAll()
	if err != nil {
		return err
	}
	pas.Devices = make(map[string]router.Entry)
	pas.Interface2Router = make(map[string]string)
	for _, router := range routers {
		pas.Devices[router.Name] = router
		for _, routerInterface := range router.Interfaces {
			pas.Interface2Router[routerInterface] = router.Name
		}
	}
	return nil
}

func (pas *PaloAltoSource) InitInterfaces(c *pango.Firewall) error {
	ethInterfaces, err := c.Network.EthernetInterface.GetAll()
	if err != nil {
		return err
	}
	pas.Interfaces = make(map[string]eth.Entry)
	pas.Interface2Subinterfaces = make(map[string][]layer3.Entry)
	for _, ethInterface := range ethInterfaces {
		pas.Interfaces[ethInterface.Name] = ethInterface
		subInterfaces, err := c.Network.Layer3Subinterface.GetAll(layer3.EthernetInterface, ethInterface.Name)
		if err != nil {
			return fmt.Errorf("layer 3 subinterfaces: %s", err)
		}
		pas.Interface2Subinterfaces[ethInterface.Name] = make([]layer3.Entry, 0, len(subInterfaces))
		pas.Interface2Subinterfaces[ethInterface.Name] = subInterfaces
	}
	return nil
}
