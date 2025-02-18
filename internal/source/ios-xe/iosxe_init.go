package iosxe

import (
	"encoding/xml"
	"fmt"

	"github.com/scrapli/scrapligo/driver/netconf"
)

func (is *IOSXESource) initDeviceInfo(d *netconf.Driver) error {
	r, err := d.Get(systemFilter)
	if err != nil {
		return fmt.Errorf("error with system filter: %s", err)
	}
	err = xml.Unmarshal(r.RawResult, &is.SystemInfo)
	if err != nil {
		return fmt.Errorf("error with unmarshaling device info: %s", err)
	}
	return nil
}

func (is *IOSXESource) initDeviceHardwareInfo(d *netconf.Driver) error {
	r, err := d.Get(hwFilter)
	if err != nil {
		return fmt.Errorf("error with hardware filter: %s", err)
	}
	err = xml.Unmarshal(r.RawResult, &is.HardwareInfo)
	if err != nil {
		return fmt.Errorf("error with unmarshaling hardware info: %s", err)
	}
	return nil
}

func (is *IOSXESource) initInterfaces(d *netconf.Driver) error {
	var ifaceReply interfaceReply
	r, err := d.Get(interfaceFilter)
	if err != nil {
		return fmt.Errorf("error with interface filter: %s", err)
	}
	err = xml.Unmarshal(r.RawResult, &ifaceReply)
	if err != nil {
		return fmt.Errorf("error with unmarshaling interfaces: %s", err)
	}
	is.Interfaces = make(map[string]iface)
	for _, iface := range ifaceReply.Interfaces {
		is.Interfaces[iface.Name] = iface
	}
	return nil
}

func (is *IOSXESource) initArpData(d *netconf.Driver) error {
	var arpReply arpReply
	r, err := d.Get(arpFilter)
	if err != nil {
		return fmt.Errorf("error with arp filter: %s", err)
	}
	err = xml.Unmarshal(r.RawResult, &arpReply)
	if err != nil {
		return fmt.Errorf("error with unmarshaling arp reply: %s", err)
	}
	is.ArpEntries = make([]arpEntry, 0)
	for _, arpVrf := range arpReply.ArpVrf {
		is.ArpEntries = append(is.ArpEntries, arpVrf.ArpOper...)
	}
	return nil
}
