package iosxe

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/options"
)

//nolint:revive
type IOSXESource struct {
	common.Config

	// IOSXE fetched data. Initialized in init functions.
	SystemInfo systemReply
	Interfaces map[string]iface
	ArpEntries []arpEntry

	// IOSXE synced data. Created in sync functions.
	NBDevice     *objects.Device
	NBInterfaces map[string]*objects.Interface // interfaceName -> netboxInterface

	// User defined relations
	HostTenantRelations map[string]string
	HostSiteRelations   map[string]string
	VlanGroupRelations  map[string]string
	VlanTenantRelations map[string]string
}

const systemFilter = `<system xmlns="http://openconfig.net/yang/system">
 <config>
 </config>
  <state>
	  <hostname/>
		<domain-name/>
  </state>
</system>
`

type systemReply struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 rpc-reply"`
	MessageID  string   `xml:"message-id,attr"`
	Hostname   string   `xml:"data>system>state>hostname"`
	DomainName string   `xml:"data>system>state>domain-name"`
}

const interfaceFilter = `<interfaces xmlns="http://openconfig.net/yang/interfaces">
    <interface>
      <name/>
      <state>
        <name/>
        <type/>
        <enabled/>
      </state>
      <ethernet xmlns="http://openconfig.net/yang/interfaces/ethernet">
        <state>
          <mac-address/>
          <auto-negotiate/>
          <port-speed/>
        </state>
      </ethernet>
    </interface>
  </interfaces>`

// InterfacesReply holds the entire response structure with the message ID and a slice of interfaces.
type interfaceReply struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 rpc-reply"`
	MessageID  string   `xml:"message-id,attr"`
	Interfaces []iface  `xml:"data>interfaces>interface"`
}

type iface struct {
	Name     string         `xml:"name"`
	State    interfaceState `xml:"state"`
	Ethernet ethernetState  `xml:"ethernet>state"`
}

// InterfaceState captures the state of the interface, including its operational status.
type interfaceState struct {
	Name    string `xml:"name"`
	Type    string `xml:"type,attr"`
	Enabled bool   `xml:"enabled"`
}

// EthernetState provides details about the Ethernet settings.
type ethernetState struct {
	MACAddress    string `xml:"mac-address"`
	AutoNegotiate bool   `xml:"auto-negotiate"`
	PortSpeed     string `xml:"port-speed"`
}

const arpFilter = `<arp-data xmlns="http://cisco.com/ns/yang/Cisco-IOS-XE-arp-oper"/>`

type arpReply struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 rpc-reply"`
	MessageID string   `xml:"message-id,attr"`
	ArpVrf    []arpVrf `xml:"data>arp-data>arp-vrf"`
}

type arpVrf struct {
	Vrf     string     `xml:"vrf"`
	ArpOper []arpEntry `xml:"arp-oper"`
}

type arpEntry struct {
	Address   string `xml:"address"`
	Interface string `xml:"interface"`
	Type      string `xml:"type"`
	Mode      string `xml:"mode"`
	HWType    string `xml:"hwtype"`
	MAC       string `xml:"hardware"`
}

func (is *IOSXESource) Init() error {
	d, err := netconf.NewDriver(
		is.SourceConfig.Hostname,
		options.WithAuthUsername(is.SourceConfig.Username),
		options.WithAuthPassword(is.SourceConfig.Password),
		options.WithPort(is.SourceConfig.Port),
		options.WithAuthNoStrictKey(), // inside container we can't confirm ssh key
	)
	if err != nil {
		return fmt.Errorf("failed to create driver: %s", err)
	}
	err = d.Open()
	if err != nil {
		return fmt.Errorf("failed to open driver: %s", err)
	}
	defer d.Close()

	// Initialize regex relations for this source
	is.VlanGroupRelations = utils.ConvertStringsToRegexPairs(is.SourceConfig.VlanGroupRelations)
	is.Logger.Debugf(is.Ctx, "VlanGroupRelations: %s", is.VlanGroupRelations)
	is.VlanTenantRelations = utils.ConvertStringsToRegexPairs(is.SourceConfig.VlanTenantRelations)
	is.Logger.Debugf(is.Ctx, "VlanTenantRelations: %s", is.VlanTenantRelations)
	is.HostTenantRelations = utils.ConvertStringsToRegexPairs(is.SourceConfig.HostTenantRelations)
	is.Logger.Debugf(is.Ctx, "HostTenantRelations: %s", is.HostTenantRelations)
	is.HostSiteRelations = utils.ConvertStringsToRegexPairs(is.SourceConfig.HostSiteRelations)
	is.Logger.Debugf(is.Ctx, "HostSiteRelations: %s", is.HostSiteRelations)

	// Initialize items from vsphere API to local storage
	initFunctions := []func(*netconf.Driver) error{
		is.initDeviceInfo,
		is.initInterfaces,
		is.initArpData,
	}

	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(d); err != nil {
			return fmt.Errorf("iosxe initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		is.Logger.Infof(is.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"), duration.Seconds())
	}
	return nil
}

func (is *IOSXESource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		is.syncDevice,
		is.syncInterfaces,
		is.syncArpTable,
	}

	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		is.Logger.Infof(is.Ctx, "Successfully synced %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync"), duration.Seconds())
	}
	return nil
}
