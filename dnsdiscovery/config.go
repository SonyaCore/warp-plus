package dns

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type DNSServers map[string]Country

//go:embed servers.json
var EServers []byte

// ScanSettings contains configurable parameters for dns discovery process
type ScanSettings struct {
	TimeoutSeconds int      `json:"timeout_seconds"`
	MaxConcurrent  int      `json:"max_concurrent"`
	RetryAttempts  int      `json:"retry_attempts"`
	TestDomains    []string `json:"test_domains"`
}

// loadServers loads server config from list of servers in json file
func loadServers(servers *DNSServers) error {
	if serverPath := os.Getenv("SERVERSPATH"); serverPath != "" {
		// Read file from path
		data, err := os.ReadFile(serverPath)
		if err != nil {
			return fmt.Errorf("error reading server file: %w", err)
		}

		if err := json.Unmarshal(data, &servers); err != nil {
			return fmt.Errorf("error parsing server file: %w", err)
		}
		return nil
	}

	if err := json.Unmarshal(EServers, &servers); err != nil {
		return fmt.Errorf("error parsing server file: %w", err)
	}
	return nil
}

// validateConfig ensures the configuration is valid
func validate(config *Discovery) error {
	if len(config.DNSServers) == 0 {
		return errors.New("no DNS servers configured")
	}

	for countryCode, countryDNS := range config.DNSServers {
		if countryDNS.Country == "" {
			return errors.New("missing country name for code " + countryCode)
		}
		if len(countryDNS.Servers) == 0 {
			return errors.New("no servers configured for country " + countryCode)
		}

		for _, server := range countryDNS.Servers {
			if server.Name == "" {
				return errors.New("missing server name in country " + countryCode)
			}
			if server.Primary == "" {
				return errors.New("missing primary DNS for " + server.Name + "in " + countryCode)
			}
			if server.Secondary == "" {
				return errors.New("missing secondary DNS for " + server.Name + "in " + countryCode)
			}
		}
	}

	return nil
}
