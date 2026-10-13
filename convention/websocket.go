package convention

import (
	"fmt"
	"net/http"
)

// WebSocketHandshake prepares a single upgrade using the same headers,
// authentication and middleware as HTTP requests. RequestTimeout and the HTTP
// client's Timeout do not limit the lifetime of the upgraded connection.
// Custom HTTP doers cannot guarantee WebSocket upgrades; use a native http.Client.
func (cfg *RequestConfig) WebSocketHandshake() (*http.Client, *http.Request, error) {
	if cfg.CustomHTTPDoer != nil {
		return nil, nil, fmt.Errorf("qca: Realtime requires an *http.Client, not a custom HTTP doer")
	}
	base := cfg.BaseURL
	if base == nil {
		base = cfg.DefaultBaseURL
	}
	if base == nil || base.Host == "" || base.User != nil || base.RawQuery != "" || base.Fragment != "" || (base.Scheme != "http" && base.Scheme != "https") {
		return nil, nil, fmt.Errorf("qca: Realtime base URL must be an HTTP(S) URL without user info, query or fragment")
	}
	request := cfg.Request.Clone(cfg.Request.Context())
	var err error
	request.URL, err = base.Parse(request.URL.String())
	if err != nil {
		return nil, nil, err
	}
	client := *cfg.HTTPClient
	client.Timeout = 0
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	handler := middlewareNext(transport.RoundTrip)
	for i := len(cfg.Middlewares) - 1; i >= 0; i-- {
		handler = applyMiddleware(cfg.Middlewares[i], handler)
	}
	client.Transport = websocketRoundTrip(handler)
	return &client, request, nil
}

type websocketRoundTrip middlewareNext

func (f websocketRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
