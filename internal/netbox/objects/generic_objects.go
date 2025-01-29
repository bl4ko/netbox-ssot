package objects

import "github.com/bl4ko/netbox-ssot/internal/constants"

type IDItem interface {
	GetID() int
	GetObjectType() constants.ContentType
	GetAPIPath() constants.APIPath
}

type OrphanItem interface {
	GetID() int
	GetObjectType() constants.ContentType
	GetAPIPath() constants.APIPath

	GetNetboxObject() *NetboxObject
}

type MACAddressOwner interface {
	GetID() int
	GetObjectType() constants.ContentType
	GetAPIPath() constants.APIPath

	GetPrimaryMACAddress() *MACAddress
	SetPrimaryMACAddress(mac *MACAddress)
}

type IPAddressOwner interface {
	GetID() int
	GetObjectType() constants.ContentType
	GetAPIPath() constants.APIPath

	GetPrimaryIPv4Address() *IPAddress
	GetPrimaryIPv6Address() *IPAddress
	SetPrimaryIPAddress(ip *IPAddress)
	SetPrimaryIPv6Address(ip *IPAddress)
}
