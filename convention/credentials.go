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
