package dns

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

const wcount = 100

// Discovery represents the main configuration structure
type Discovery struct {
	DNSServers
	Settings ScanSettings
	l        *slog.Logger
}

func NewDiscovery(logger *slog.Logger) *Discovery {
	return &Discovery{
		DNSServers: make(DNSServers),
		Settings: ScanSettings{
			TimeoutSeconds: 2,
			MaxConcurrent:  100,
			RetryAttempts:  2,
			TestDomains:    []string{"google.com", "cloudflare.com"},
		},
		l: logger,
	}
}

// Init loads the DNS configuration from multiple possible locations
func (d *Discovery) Init(l *slog.Logger) (*Discovery, error) {

	config := NewDiscovery(l)

	// load servers
	if err := loadServers(&config.DNSServers); err != nil {
		return nil, err
	}

	// validate config
	if err := validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

// Scan performs the DNS scanning and timing measurement
func (d *Discovery) Scan() []Result {
	var dispatcher = 1
	var totalOps = 0
	var results = make([]Result, 0, totalOps)

	for _, countryData := range d.DNSServers {
		totalOps += len(countryData.Servers) * 2 // Primary and Secondary IPs
	}

	resultsChan := make(chan Result, totalOps)

	workerCount := wcount
	jobs := make(chan job, totalOps)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(d, &wg, jobs, resultsChan)
	}

	// queue jobs
	go func() {
		for countryCode, countryData := range d.DNSServers {
			for _, server := range countryData.Servers {
				for _, ip := range []string{server.Primary, server.Secondary} {
					d.l.Debug("DNS lookup", "dispatcher", dispatcher, "ip", ip, "subsystem", "dnsdiscovery")
					dispatcher++
					jobs <- job{
						countryCode: countryCode,
						country:     countryData.Country,
						provider:    server.Name,
						ip:          ip,
					}
				}
			}
		}
		close(jobs)
	}()

	// Start collector goroutine
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// append results
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func (d *Discovery) Fastest(results []Result) string {
	var fastest *Result

	for i, result := range results {
		// Skip unreachable servers
		if !result.isReachable {
			continue
		}

		// Initialize fastest with first reachable server
		if fastest == nil {
			fastest = &results[i]
			continue
		}

		// Update if we find a faster server
		if result.responseMs < fastest.responseMs {
			fastest = &results[i]
		}
	}

	return fastest.countryCode
}

// lookup tests a single DNS server and measures response time
func (d *Discovery) lookup(countryCode, country, provider, ip string) Result {

	result := Result{
		countryCode: countryCode,
		country:     country,
		provider:    provider,
		ip:          ip,
		isReachable: false,
	}

	resolver := NewResolver(ip)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	startTime := time.Now()
	_, err := resolver.LookupHost(ctx, d.Settings.TestDomains[0])
	duration := time.Since(startTime)

	if err == nil {
		result.isReachable = true
		result.responseMs = float64(duration.Milliseconds())
	} else {
		result.responseMs = -1
	}
	return result
}
