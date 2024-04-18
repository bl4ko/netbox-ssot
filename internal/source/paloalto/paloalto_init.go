package paloalto

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/eth"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer3"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"
	"github.com/PaloAltoNetworks/pango/netw/zone"
	"github.com/PaloAltoNetworks/pango/vsys"
)

// Init system info collects system info from paloalto.
func (pas *PaloAltoSource) initSystemInfo(c *pango.Firewall) error {
	pas.SystemInfo = c.Client.SystemInfo
	return nil
}

func (pas *PaloAltoSource) initVirtualSystems(c *pango.Firewall) error {
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

func (pas *PaloAltoSource) initVirtualRouters(c *pango.Firewall) error {
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

// initInterfaces collects all ethernet interfaces and subinterfaces
// from paloalto API. It stores them as attribute of the paloalto source.
func (pas *PaloAltoSource) initInterfaces(c *pango.Firewall) error {
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

// Structs to parse xml arp data response.
type ArpData struct {
	XMLName xml.Name  `xml:"response"`    // This ensures the root element is correctly recognized
	Status  string    `xml:"status,attr"` // This captures the "status" attribute in the response tag
	Result  ArpResult `xml:"result"`      // This nests the result struct under the result tag
}

type ArpResult struct {
	Max     int        `xml:"max"`
	Total   int        `xml:"total"`
	Timeout int        `xml:"timeout"`
	DP      string     `xml:"dp"`
	Entries []ArpEntry `xml:"entries>entry"` // Correct path to entry elements
}

type ArpEntry struct {
	Status    string `xml:"status"`
	IP        string `xml:"ip"`
	MAC       string `xml:"mac"`
	TTL       string `xml:"ttl"`
	Interface string `xml:"interface"`
	Port      string `xml:"port"`
}

// initArpData collects all arp entries from the paloalto source.
// It stores them as attribute of the paloalto source.
func (pas *PaloAltoSource) initArpData(c *pango.Firewall) error {
	if pas.SourceConfig.CollectArpData {
		var arpData ArpData
		arpXMLString := "<show><arp><entry name='all'/></arp></show>"
		arpXMLResponse, err := c.Op(arpXMLString, "", nil, nil)
		if err != nil {
			return fmt.Errorf("init arp data: %s", err)
		}
		err = xml.Unmarshal(arpXMLResponse, &arpData)
		if err != nil {
			return fmt.Errorf("init arp data: %s", err)
		}
		if arpData.Result.Entries != nil {
			pas.ArpData = arpData.Result.Entries
		}
	}
	return nil
}
