package objects

import "fmt"

// https://github.com/netbox-community/netbox/blob/2e74952ac6cc68348284dee8b9517fe0a36179a2/netbox/wireless/choices.py#L16
type WirelessLanStatus struct {
	Choice
}

var (
	WirelessLanStatusActive     = WirelessLanStatus{Choice{Value: "active", Label: "Active"}}
	WirelessLanStatusReserved   = WirelessLanStatus{Choice{Value: "reserved", Label: "Reserved"}}
	WirelessLanStatusDisabled   = WirelessLanStatus{Choice{Value: "disabled", Label: "Disabled"}}
	WirelessLanStatusDeprecated = WirelessLanStatus{Choice{Value: "deprecated", Label: "Deprecated"}}
)

type WirelessLanAuthType struct {
	Choice
}

var (
	WirelessLanAuthTypeOpen          = WirelessLanAuthType{Choice{Value: "open", Label: "Open"}}
	WirelessLanAuthTypeWep           = WirelessLanAuthType{Choice{Value: "wep", Label: "WEP"}}
	WirelessLanAuthTypeWpaPersonal   = WirelessLanAuthType{Choice{Value: "wpa-personal", Label: "WPA Personal (PSK)"}}
	WirelessLanAuthTypeWpaEnterprise = WirelessLanAuthType{Choice{Value: "wpa-enterprise", Label: "WPA Enterprise"}}
)

type WirelessLanAuthCipher struct {
	Choice
}

var (
	WirelessLanAuthCipherAuto = WirelessLanAuthCipher{Choice{Value: "auto", Label: "Auto"}}
	WirelessLanAuthCipherTkip = WirelessLanAuthCipher{Choice{Value: "tkip", Label: "TKIP"}}
	WirelessLanAuthCipherAes  = WirelessLanAuthCipher{Choice{Value: "aes", Label: "AES"}}
)

type WirelessLan struct {
	NetboxObject
	// SSID is the name of the wireless lan. This field is required.
	SSID string `json:"ssid,omitempty"`
	// Vlan that the wireless lan is associated with.
	Vlan *Vlan `json:"vlan,omitempty"`
	// WirelessLanStatus is the status of the wireless lan. This field is required.
	Status *WirelessLanStatus `json:"status,omitempty"`
	// Tenant of the wireless lan.
	Tenant *Tenant `json:"tenant,omitempty"`
	// AuthType is the authentication type of the wireless lan.
	AuthType *WirelessLanAuthType `json:"auth_type,omitempty"`
	// AuthCipher is the authentication cipher of the wireless lan.
	AuthCipher *WirelessLanAuthCipher `json:"auth_cipher,omitempty"`
	// AuthPsk is the pre-shared key of the wireless lan.
	AuthPsk string `json:"auth_psk,omitempty"`
	// Comments is the comments about the wireless lan.
	Comments string `json:"comments,omitempty"`
}

func (wl WirelessLan) String() string {
	return fmt.Sprintf("WirelessLan{SSID: %s, Vlan: %s}", wl.SSID, wl.Vlan)
}
