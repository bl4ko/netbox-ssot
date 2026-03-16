package objects

import (
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

type WirelessLANGroup struct {
	NetboxObject
	// Name is the name of the wireless lan group. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slug of the wireless lan group. This field is required.
	Slug string `json:"slug,omitempty"`
	// Parent is the parent wireless lan group.
	Parent *WirelessLANGroup `json:"parent,omitempty"`
}

func (wlg WirelessLANGroup) String() string {
	return fmt.Sprintf("WirelessLANGroup{Name: %s, Slug: %s}", wlg.Name, wlg.Slug)
}

// WirelessLANGroup implements IDItem interface.
func (wlg *WirelessLANGroup) GetID() int {
	return wlg.ID
}
func (wlg *WirelessLANGroup) GetObjectType() constants.ContentType {
	return constants.ContentTypeWirelessLANGroup
}
func (wlg *WirelessLANGroup) GetAPIPath() constants.APIPath {
	return constants.WirelessLANGroupsAPIPath
}

// WirelessLANGroup implements OrphanItem interface.
func (wlg *WirelessLANGroup) GetNetboxObject() *NetboxObject {
	return &wlg.NetboxObject
}

type WirelessLANStatus struct {
	Choice
}

var (
	WirelessLanStatusActive     = WirelessLANStatus{Choice{Value: "active", Label: "Active"}}
	WirelessLanStatusReserved   = WirelessLANStatus{Choice{Value: "reserved", Label: "Reserved"}}
	WirelessLanStatusDisabled   = WirelessLANStatus{Choice{Value: "disabled", Label: "Disabled"}}
	WirelessLanStatusDeprecated = WirelessLANStatus{
		Choice{Value: "deprecated", Label: "Deprecated"},
	}
)

type WirelessLANAuthType struct {
	Choice
}

var (
	WirelessLanAuthTypeOpen        = WirelessLANAuthType{Choice{Value: "open", Label: "Open"}}
	WirelessLanAuthTypeWep         = WirelessLANAuthType{Choice{Value: "wep", Label: "WEP"}}
	WirelessLanAuthTypeWpaPersonal = WirelessLANAuthType{
		Choice{Value: "wpa-personal", Label: "WPA Personal (PSK)"},
	}
	WirelessLanAuthTypeWpaEnterprise = WirelessLANAuthType{
		Choice{Value: "wpa-enterprise", Label: "WPA Enterprise"},
	}
)

type WirelessLANAuthCipher struct {
	Choice
}

var (
	WirelessLANAuthCipherAuto = WirelessLANAuthCipher{Choice{Value: "auto", Label: "Auto"}}
	WirelessLANAuthCipherTkip = WirelessLANAuthCipher{Choice{Value: "tkip", Label: "TKIP"}}
	WirelessLANAuthCipherAes  = WirelessLANAuthCipher{Choice{Value: "aes", Label: "AES"}}
)

type WirelessLAN struct {
	NetboxObject
	// SSID is the name of the wireless lan. This field is required.
	SSID string `json:"ssid,omitempty"`
	// Vlan that the wireless lan is associated with.
	Vlan *Vlan `json:"vlan,omitempty"`
	// Group is the group of the wireless lan.
	Group *WirelessLANGroup `json:"group,omitempty"`
	// Status is the status of the wireless lan. This field is required.
	Status *WirelessLANStatus `json:"status,omitempty"`
	// Tenant of the wireless lan.
	Tenant *Tenant `json:"tenant,omitempty"`
	// AuthType is the authentication type of the wireless lan.
	AuthType *WirelessLANAuthType `json:"auth_type,omitempty"`
	// AuthCipher is the authentication cipher of the wireless lan.
	AuthCipher *WirelessLANAuthCipher `json:"auth_cipher,omitempty"`
	// AuthPsk is the pre-shared key of the wireless lan.
	AuthPsk string `json:"auth_psk,omitempty"`
	// Comments is the comments about the wireless lan.
	Comments string `json:"comments,omitempty"`
}

func (wl WirelessLAN) String() string {
	return fmt.Sprintf("WirelessLAN{SSID: %s}", wl.SSID)
}

// WirelessLAN implements IDItem interface.
func (wl *WirelessLAN) GetID() int {
	return wl.ID
}
func (wl *WirelessLAN) GetObjectType() constants.ContentType {
	return constants.ContentTypeWirelessLAN
}
func (wl *WirelessLAN) GetAPIPath() constants.APIPath {
	return constants.WirelessLANsAPIPath
}

// WirelessLAN implements OrphanItem interface.
func (wl *WirelessLAN) GetNetboxObject() *NetboxObject {
	return &wl.NetboxObject
}
