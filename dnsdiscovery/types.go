package dns

// Server represents a single DNS server
type Server struct {
	Name      string `json:"name"`
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}

// Country represents DNS servers for a country
type Country struct {
	Country string   `json:"country"`
	Servers []Server `json:"servers"`
}

// Result represents the scanning result for a single DNS server
type Result struct {
	CountryCode string
	Country     string
	Provider    string
	IP          string
	ResponseMs  float64
	IsReachable bool
}
