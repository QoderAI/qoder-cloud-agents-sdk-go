package convention

import (
	"net/http"
	"time"
)

// DefaultHTTPClient returns a client with a ten-minute response-header timeout.
// Response bodies, including event streams, have no default total deadline.
// A custom http.DefaultTransport wrapper is preserved without modification.
func DefaultHTTPClient() *http.Client {
	if transport, ok := http.DefaultTransport.(*http.Transport); ok {
		transport = transport.Clone()
		transport.ResponseHeaderTimeout = 10 * time.Minute
		return &http.Client{Transport: transport}
	}
	return &http.Client{Transport: http.DefaultTransport}
}
