package dns

import (
	"fmt"
	"strings"
)

func validateDNS(host string) error {
	if len(host) > 255 {
		return fmt.Errorf("lookup host exceeds maximum length of 255 characters")
	}

	// Split hostname into labels and validate each
	labels := strings.Split(host, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return fmt.Errorf("invalid label length in hostname: %s", host)
		}

		// Check if label contains only allowed characters (RFC 1035)
		for _, c := range label {
			if !isValidHostnameChar(c) {
				return fmt.Errorf("invalid character in hostname: %c", c)
			}
		}
	}
	return nil
}

// Helper function to check valid hostname characters based on RFC 1035 standards
func isValidHostnameChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '-'
}
