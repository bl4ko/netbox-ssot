package utils

import (
	"fmt"
	"net"
	"strings"
)

// Function that receives ipAddress and performs a reverse lookup
// to get the hostname. If the reverse lookup fails, it returns an empty string.
func ReverseLookup(ipAddress string) string {
	names, err := net.LookupAddr(ipAddress)
	if err != nil {
		return ""
	}

	if len(names) > 0 {
		domain := strings.TrimSuffix(names[0], ".")
		return domain
	}

	return ""
}

// Function that receives hostname and performs a forward lookup
// to get the IP address. If the forward lookup fails, it returns an empty string.
func Lookup(hostname string) string {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return ""
	}

	if len(ips) > 0 {
		return ips[0].String()
	}

	return ""
}

// Function that converts string representation of ipv4 mask (e.g. 255.255.255.128) to
// bit representation (e.g. 25).
func MaskToBits(mask string) (int, error) {
	ipMask := net.IPMask(net.ParseIP(mask).To4())
	if ipMask == nil {
		return 0, fmt.Errorf("invalid mask: %s", mask)
	}
	ones, _ := ipMask.Size()
	return ones, nil
}

// GetIPVersion returns the version of the IP address.
// It returns 4 for IPv4, 6 for IPv6, and 0 if the IP address is invalid.
func GetIPVersion(ipAddress string) int {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0
	}
	if ip.To4() != nil {
		return 4
	}
	return 6
}

// Function that checks if given IP address is part of the
// given subnet.
// e.g. ipAddress "172.31.4.129" and subnet "172.31.4.145/25"
// Return true.
func SubnetContainsIpAddress(ipAddress string, subnet string) bool {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false
	}
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}
	return ipnet.Contains(ip)
}
