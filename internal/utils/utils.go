package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/logger"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Validates array of regex relations
// Regex relation is a string of format "regex = value".
func ValidateRegexRelations(regexRelations []string) error {
	for _, regexRelation := range regexRelations {
		relation := strings.Split(regexRelation, "=")
		if len(relation) != len([]string{"regex", "value"}) {
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
// where key is regex and value is value.
func ConvertStringsToRegexPairs(input []string) map[string]string {
	output := make(map[string]string, len(input))
	for _, s := range input {
		pair := strings.Split(s, "=")
		output[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}
	return output
}

// Matches input string to a regex from input map patterns,
// and returns the value. If there is no match, it returns an empty string.
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
// e.g. "  @Test@ " -> "test".
func Slugify(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")

	// Remove characters except lowercase letters, numbers, underscores, hyphens
	reg := regexp.MustCompile("[^a-z0-9_-]+")
	name = reg.ReplaceAllString(name, "")
	return name
}

// Converts string name to its alphanumeric representation
// with underscored only.
func Alphanumeric(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")

	// Remove characters except lowercase letters, numbers, underscores
	reg := regexp.MustCompile("[^a-z0-9_]+")
	name = reg.ReplaceAllString(name, "")
	return name
}

// Function that takes osType and osVersion and returns a
// an universal platform name that then can be shared between
// multiple objects.
func GeneratePlatformName(osType string, osVersion string) string {
	if osType == "" {
		osType = constants.DefaultOSName
	}
	if osVersion == "" {
		osVersion = constants.DefaultOSVersion
	}
	return fmt.Sprintf("%s %s", osType, osVersion)
}

// Function that returns true if the given string
// representing an virtual machine interface name is valid and false otherwise.
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

// Converts strings of format fieldName = value to map[fieldName] = value.
func ConvertStringsToPairs(input []string) map[string]string {
	output := make(map[string]string, len(input))
	for _, s := range input {
		pair := strings.Split(s, "=")
		output[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}
	return output
}

// Define a new type that implements the runes.Set interface.
type mnSet struct{}

// Contains implements the runes.Set interface for mnSet.
// It returns true if the rune is a nonspacing mark.
func (mnSet) Contains(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

// Function that removes diacritics and normalizes strings.
func removeDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(mnSet{}), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// Function that matches array of names to array of emails (which could be subset of it).
// It returns a map of name -> email.
//
// E.g. names = ["John Doe", "Jane Doe"], emails = ["jane.doe@example"]
// Output: map["Jane Doe"] = "jane.doe@example".
func MatchNamesWithEmails(names []string, emails []string, logger *logger.Logger) map[string]string {
	normalizedNames := make(map[string]string) // Map for easy lookup
	for _, name := range names {
		// Normalize name: remove diacritics, spaces, and convert to lowercase
		normalized := strings.ReplaceAll(strings.ToLower(removeDiacritics(name)), " ", "")
		normalizedNames[normalized] = name
	}

	matches := make(map[string]string)
	for _, email := range emails {
		username := strings.Split(strings.ToLower(email), "@")[0]
		username = strings.ReplaceAll(username, ".", "") // Remove common separators
		username = strings.ReplaceAll(username, "_", "")

		// Try to find a match
		if name, exists := normalizedNames[username]; exists {
			matches[name] = email
		} else {
			// Handle no match or implement additional matching logic
			logger.Warningf("No direct match found for email: %s\n", email)
		}
	}
	return matches
}
