package dns

import (
	"context"
	"net"
	"time"
)

// NewResolver returns new *net.Resolver for lookup process
func NewResolver(ip string) *net.Resolver {
	return &net.Resolver{
		// using built-in DNS resolver netdns=go
		PreferGo: true,

		// setup net Dialer
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			// dns port
			const p = ":53"

			// dialer setup
			d := net.Dialer{
				Timeout: time.Second * 2,
			}
			return d.DialContext(ctx, "udp", ip+p)
		},
	}
}
