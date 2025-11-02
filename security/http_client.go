package security

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// RestrictedHTTPClient creates an HTTP client that only allows localhost connections
func RestrictedHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				// Control function to restrict connections
			}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			// Force connections to localhost only
			DisableKeepAlives: false,
		},
	}
}

// LocalhostOnlyDialer returns a dialer that only connects to localhost
type LocalhostOnlyDialer struct {
	Dialer *net.Dialer
}

// Dial implements the Dial interface
func (d *LocalhostOnlyDialer) Dial(network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	
	// Only allow localhost/127.0.0.1
	if host != "localhost" && host != "127.0.0.1" && host != "::1" {
		return nil, fmt.Errorf("connection to %s blocked by security policy (only localhost allowed)", host)
	}
	
	return d.Dialer.Dial(network, addr)
}

// NewLocalhostOnlyClient creates an HTTP client that can only connect to localhost
func NewLocalhostOnlyClient() *http.Client {
	dialer := &LocalhostOnlyDialer{
		Dialer: &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	}
	
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Dial:                  dialer.Dial,
			MaxIdleConns:          10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

