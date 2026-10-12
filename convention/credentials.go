package convention

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
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
	return "qca: could not resolve authentication method. Configure an access token or credential, or QODER_ACCESS_TOKEN (QODER_PAT is supported in PAT mode)"
}

// RequireCredential resolves the selected mode per attempt. Explicit headers
// override injection; custom authentication middleware remains supported.
func RequireCredential() RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		if r.authInstalled {
			return nil
		}
		r.authInstalled = true
		r.authMiddlewareIndex = len(r.Middlewares)
		r.Middlewares = append(r.Middlewares, credentialMiddleware(r))
		return nil
	})
}

func credentialMiddleware(cfg *RequestConfig) middleware {
	return func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
		mode := cfg.AuthMode
		if sa, ok := cfg.Credential.(*ServiceAccountCredential); ok {
			if mode == "pat" || cfg.AuthToken != "" {
				return nil, fmt.Errorf("qca: ServiceAccountCredential requires SA mode without a static token")
			}
			mode = "sa"
			if err := sa.validateRegion(req.URL.Hostname()); err != nil {
				return nil, err
			}
		}
		if mode == "" {
			mode = "pat"
		}
		if mode != "pat" && mode != "sa" {
			return nil, fmt.Errorf("qca: auth mode must be pat or sa")
		}
		if req.Header.Get("Authorization") != "" {
			return next(req)
		}
		var token string
		if cfg.Credential != nil {
			var err error
			token, err = cfg.Credential.Token(req.Context())
			if err != nil {
				return nil, err
			}
			if strings.TrimSpace(token) == "" {
				return nil, fmt.Errorf("qca: credential returned an empty access token")
			}
		} else if !cfg.authDeleted {
			token = cfg.envAccessToken
			if token == "" && mode == "pat" {
				token = cfg.envPAT
			}
		}
		if token != "" {
			cloned := req.Clone(req.Context())
			cloned.Header.Set("Authorization", "Bearer "+token)
			return next(cloned)
		}
		if cfg.APIKey != "" || len(cfg.Middlewares) > 1 ||
			(cfg.authDeleted && (cfg.AuthToken != "" || cfg.envAccessToken != "" || (mode == "pat" && cfg.envPAT != ""))) {
			return next(req)
		}
		return nil, &NoCredentialsError{}
	}
}

// EnvironmentAuthentication captures defaults in the common layer. Explicit
// credentials win; SA mode never falls back to the personal PAT.
func EnvironmentAuthentication() RequestOption {
	access, pat := os.Getenv("QODER_ACCESS_TOKEN"), os.Getenv("QODER_PAT")
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.envAccessToken, r.envPAT = access, pat
		return RequireCredential().Apply(r)
	})
}

// MarkAuthorizationDeleted suppresses static/environment header restoration;
// dynamic credentials retain their existing header-deletion behavior.
func (r *RequestConfig) MarkAuthorizationDeleted() { r.authDeleted = true }
