# DNS Discovery

A Go package for discovering and benchmarking DNS servers across multiple countries. This tool helps measure DNS response times and identify the fastest available DNS servers.

## Features

- Support for DNS servers across 30+ countries
- Concurrent DNS server testing
- Response time measurement
- Configurable timeout and retry settings
- Environment variable configuration
- Embedded default server list
- RFC 1035 compliant hostname validation

## Usage

```go
package main

import (
    "github.com/bepass-org/warp-plus/dnsdiscovery"
    "log/slog"
    "fmt"
    "dns"
)

func main() {
    logger := slog.Default()
    
    // Initialize discovery
    d, err := dns.NewDiscovery(logger).Init(logger)
    if err != nil {
        logger.Error("Failed to initialize DNS discovery", "error", err)
        return
    }
    
    // Run the scan
    results := d.Scan()
    
    // Find the fastest DNS server
    fastest := d.Fastest(results)
    if fastest != nil {
        fmt.Printf("Fastest DNS: %s (%s) - %.2fms\n", 
            fastest.IP, fastest.Country, fastest.ResponseMs)
    }
}
```

## Configuration

### Environment Variables

- `SERVERSPATH`: Path to custom DNS servers JSON file
- `LOOKUPHOST`: Custom domain to test DNS servers (default: "google.com")

### Default Settings

```go
Settings {
    TimeoutSeconds: 2,
    MaxConcurrent: 100,
    RetryAttempts: 2,
    TestDomains: []string{"google.com", "cloudflare.com"},
}
```

### Server Configuration Format

The DNS servers are configured in JSON format:

```json
{
  "US": {
    "country": "United States",
    "servers": [
      {
        "name": "Provider Name",
        "primary": "1.1.1.1",
        "secondary": "1.0.0.1"
      }
    ]
  }
}
```

## Features in Detail

1. **Concurrent Scanning**: Uses worker pools to test multiple DNS servers simultaneously
2. **Response Time Measurement**: Measures DNS query response times in milliseconds
3. **Server Validation**: Validates DNS server configurations and hostnames
4. **Fallback Mechanism**: Uses embedded server list if no custom configuration is provided
5. **Custom Test Domains**: Supports configurable test domains via environment variables

## Contributing
Contributions to DNS Discovery are welcome. Please ensure to follow the project's coding standards and submit detailed pull requests.