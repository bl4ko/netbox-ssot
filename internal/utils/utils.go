package utils

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// Validates array of regex relations
// Regex relation is a string of format "regex = value"
func ValidateRegexRelations(regexRelations []string) error {
	for _, regexRelation := range regexRelations {
		relation := strings.Split(regexRelation, "=")
		if len(relation) != 2 {
			return fmt.Errorf("invalid regex relation: %s. Should be of format: regex = value", regexRelation)
		}
		regexStr := strings.TrimSpace(relation[0])
		_, err := regexp.Compile(regexStr)
		if err != nil {
			return fmt.Errorf("invalid regex: %s, in relation: %s", regexStr, regexRelation)
		}
	}
	return nil
}

// Converts array of strings, that are of form "regex = value", to a map
// where key is regex and value is value
func ConvertStringsToRegexPairs(input []string) map[string]string {
	output := make(map[string]string, len(input))
	for _, s := range input {
		pair := strings.Split(s, "=")
		output[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}
	return output
}

// Matches input string to a regex from input map patterns,
// and returns the value. If there is no match, it returns an empty string
func MatchStringToValue(input string, patterns map[string]string) (string, error) {
	for regex, value := range patterns {
		matched, err := regexp.MatchString(regex, input)
		if err != nil {
			return "", err
		}
		if matched {
			return value, nil
		}
	}
	return "", nil
}

// Converts string name to its slugified version.
// Slugified version can only contain: lowercase letters, numbers,
// underscores or hyphens.
// e.g. "My Name" -> "my-name"
// e.g. "  @Test@ " -> "test"
func Slugify(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")

	// Remove characters except lowercase letters, numbers, underscores, hyphens
	reg, _ := regexp.Compile("[^a-z0-9_-]+")
	name = reg.ReplaceAllString(name, "")
	return name
}

// Function that takes osType and osVersion and returns a
// an universal platform name that then can be shared between
// multiple objects.
func GeneratePlatformName(osType string, osVersion string) string {
	if osType == "" {
		osType = "Generic OS"
	}
	if osVersion == "" {
		osVersion = "Generic Version"
	}
	return fmt.Sprintf("%s %s", osType, osVersion)
}

// Function that receives ipAddress and performs a reverse lookup
// to get the hostname. If the reverse lookup fails, it returns an empty string.
func ReverseLookup(ipAddress string) string {
	names, err := net.LookupAddr(ipAddress)
	if err != nil {
		return ""
	}

	if len(names) > 0 {
		return names[0]
	}

	return ""
}

// Function that returns true if the given string
// representing an vm's interface name is valid and false otherwise.
// Valid interface names are the ones that pass regex filtering.
func IsVMInterfaceNameValid(vmIfaceName string) (bool, error) {
	ifaceFilter := map[string]string{
		"^(docker|cali|flannel|veth|br-|cni|tun|tap|lo|virbr|vxlan|wg|kube-bridge|kube-ipvs)\\w*": "yes",
	}

	ifaceName, err := MatchStringToValue(vmIfaceName, ifaceFilter)
	if err != nil {
		return false, err
	}

	if ifaceName == "yes" {
		return false, nil
	}

	return true, nil
}

// Function that converts string representation of ipv4 mask (e.g. 255.255.255.128) to
// bit representation (e.g. 25)
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
// Return true
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

// ExtractFunctionName attempts to extract the name of a function regardless of its signature.
// Note: This function sacrifices type safety and assumes the caller ensures the correct usage.
func ExtractFunctionName(i interface{}) string {
	// Ensure the provided interface is actually a function
	if reflect.TypeOf(i).Kind() != reflect.Func {
		panic("Argument to extractFunctionName is not a function!")
	}

	fullFuncName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	funcNameParts := strings.Split(fullFuncName, ".")
	return funcNameParts[len(funcNameParts)-1]
}

// Netbox description field
