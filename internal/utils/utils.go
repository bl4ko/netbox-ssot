package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/bl4ko/netbox-ssot/internal/logger"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Deref will dereference a generic pointer.
// It will fail at compile time if pt is not a pointer.
// This is useful for functions that require pointers to structs and not structs directly.
func Deref[T any](ptr *T) T {
	return *ptr
}

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
	name = strings.ReplaceAll(name, " ", "-")

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

// Function that returns true if the given interface name should be
// filtered out, or false if it shouldn't be.
func FilterInterfaceName(ifaceName string, ifaceFilter string) bool {
	if ifaceFilter == "" {
		return false
	}
	compiledFilter := regexp.MustCompile(ifaceFilter)
	return compiledFilter.MatchString(ifaceName)
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
	funcNameFull := funcNameParts[len(funcNameParts)-1]
	funcNameFull = strings.TrimSuffix(funcNameFull, "-fm")
	return funcNameFull
}

// ExtractFunctionNameWithTrimPrefix extracts the function name and trims the prefix.
// Note: This function sacrifices type safety and assumes the caller ensures the correct usage.
func ExtractFunctionNameWithTrimPrefix(i interface{}, prefix string) string {
	funcNameFull := ExtractFunctionName(i)
	return strings.TrimPrefix(funcNameFull, prefix)
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
func MatchNamesWithEmails(ctx context.Context, names []string, emails []string, logger *logger.Logger) map[string]string {
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
			logger.Debugf(ctx, "No direct match found for email: %s", email)
		}
	}
	return matches
}

// Function that loads additional certs from the given certDirPath.
// In case empty certPath is provided, deefault cert pool is returned.
func LoadExtraCert(certPath string) (*x509.CertPool, error) {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("system cert pool: %s", err)
	}
	if certPath != "" {
		// Read in the cert file
		certFile, err := os.ReadFile(certPath)
		if err != nil {
			return nil, fmt.Errorf("read cert file: %s", err)
		}
		ok := rootCAs.AppendCertsFromPEM((certFile))
		if !ok {
			return nil, fmt.Errorf("failed to append cert from PEM file: %s", certFile)
		}
	}
	return rootCAs, nil
}

// Function that returns Transport config with extra cert from
// cert path loaded.
func LoadExtraCertInTransportConfig(certPath string) (*http.Transport, error) {
	rootCAs, err := LoadExtraCert(certPath)
	if err != nil {
		return nil, fmt.Errorf("load extra cert: %s", err)
	}
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: rootCAs,
		},
	}, nil
}
