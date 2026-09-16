package convention

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type Credential interface {
	Token(context.Context) (string, error)
}

type PATCredential struct{ token string }

func NewPATCredential(token string) (*PATCredential, error) {
	if token == "" {
		return nil, fmt.Errorf("PAT must not be empty")
	}
	return &PATCredential{token: token}, nil
}

func PATCredentialFromEnv(name string) (*PATCredential, error) {
	if name == "" {
		name = "QODER_PAT"
	}
	return NewPATCredential(os.Getenv(name))
}

func (c *PATCredential) Token(context.Context) (string, error) { return c.token, nil }

// NewCredentialTransport adds Bearer authentication without mutating the caller's request.
// An explicit Authorization header takes precedence. Tokens are fetched per attempt.
func NewCredentialTransport(credential Credential, next http.RoundTripper) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &credentialTransport{credential: credential, next: next}
}

type credentialTransport struct {
	credential Credential
	next       http.RoundTripper
}

func (transport *credentialTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	cloned := request.Clone(request.Context())
	cloned.Header = request.Header.Clone()
	if cloned.Header.Get("Authorization") == "" {
		if transport.credential == nil {
			return nil, fmt.Errorf("credential must not be nil")
		}
		token, err := transport.credential.Token(request.Context())
		if err != nil {
			return nil, err
		}
		cloned.Header.Set("Authorization", "Bearer "+token)
	}
	return transport.next.RoundTrip(cloned)
}

// NoCredentialsError is returned before a request is sent when no authentication
// could be resolved. It mirrors the request-time validation in the Python and
// TypeScript SDKs and is matchable with errors.As.
type NoCredentialsError struct{}

func (NoCredentialsError) Error() string {
	return "qca: could not resolve authentication method. Expected one of pat or credential to be set, or the QODER_PAT environment variable to be configured"
}

// RequireCredential fails a request that reaches the wire without any authentication.
// It is a defense-in-depth footgun-check on top of the server-side gateway; if the caller
// installed a credential via option.WithPAT / option.WithCredential / a custom middleware,
// or set the Authorization header directly, the sentinel steps aside.
func RequireCredential() RequestOption {
	preBuiltErr := &NoCredentialsError{}
	return RequestOptionFunc(func(r *RequestConfig) error {
		cfg := r
		check := func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
			if cfg.AuthToken != "" || cfg.APIKey != "" {
				return next(req)
			}
			if len(cfg.Middlewares) > 1 {
				return next(req)
			}
			if req.Header.Get("Authorization") != "" {
				return next(req)
			}
			return nil, preBuiltErr
		}
		r.Middlewares = append(r.Middlewares, check)
		return nil
	})
}
