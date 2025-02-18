package iosxe

import "encoding/xml"

// hardwareReply is the top-level structure for the hardware response.
type hardwareReply struct {
	XMLName   xml.Name      `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 rpc-reply"`
	MessageID string        `xml:"message-id,attr"`
	Inventory []HWInventory `xml:"data>device-hardware-data>device-hardware>device-inventory"`
}

// HWInventory holds the part and serial numbers from each inventory entry.
type HWInventory struct {
	Type         string `xml:"hw-type"`
	DevIndex     string `xml:"hw-dev-index"`
	Version      string `xml:"version"`
	PartNumber   string `xml:"part-number"`
	SerialNumber string `xml:"serial-number"`
	Description  string `xml:"hw-description"`
	DevName      string `xml:"dev-name"`
	Class        string `xml:"hw-class"`
}

type systemReply struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 rpc-reply"`
	MessageID  string   `xml:"message-id,attr"`
	Hostname   string   `xml:"data>system>state>hostname"`
	DomainName string   `xml:"data>system>state>domain-name"`
}

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
	Name        string `xml:"name"`
	Type        string `xml:"type,attr"`
	Enabled     bool   `xml:"enabled"`
	Description string `xml:"description"`
}

// EthernetState provides details about the Ethernet settings.
type ethernetState struct {
	MACAddress    string `xml:"mac-address"`
	AutoNegotiate bool   `xml:"auto-negotiate"`
	PortSpeed     string `xml:"port-speed"`
}

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
