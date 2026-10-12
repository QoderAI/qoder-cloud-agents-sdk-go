package convention

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ServiceAccountCredentialConfig configures lazy SA Key exchange. Region and
// scopes must be explicit. The supplied HTTP client is copied, not owned.
type ServiceAccountCredentialConfig struct {
	SAKey      string
	Region     string
	Scopes     []string
	TTLSeconds int
	Timeout    time.Duration
	HTTPClient *http.Client
}
type tokenExchange struct {
	done      chan struct{}
	token     string
	err       error
	safeUntil time.Time
}

// ServiceAccountCredential shares an instance-owned SAT cache across clients.
// Key rotation uses a new instance; no refresh runs while idle.
type ServiceAccountCredential struct {
	key, region, scope            string
	ttl                           int
	timeout                       time.Duration
	client                        *http.Client
	mu                            sync.Mutex
	now                           func() time.Time
	token                         string
	refreshAt, safeUntil, retryAt time.Time
	failure                       error
	transient                     bool
	inflight                      *tokenExchange
}

func NewServiceAccountCredential(config ServiceAccountCredentialConfig) (*ServiceAccountCredential, error) {
	if strings.TrimSpace(config.SAKey) == "" {
		return nil, fmt.Errorf("qca: SA key must not be empty")
	}
	if config.Region != "cn" && config.Region != "global" {
		return nil, fmt.Errorf("qca: region must be cn or global")
	}
	if len(config.Scopes) == 0 {
		return nil, fmt.Errorf("qca: scopes must not be empty")
	}
	scopes := []string{}
	seen := map[string]bool{}
	for _, scope := range config.Scopes {
		if scope == "" || strings.ContainsAny(scope, " \t\r\n") {
			return nil, fmt.Errorf("qca: invalid scope name")
		}
		if !seen[scope] {
			scopes = append(scopes, scope)
			seen[scope] = true
		}
	}
	if config.TTLSeconds == 0 {
		config.TTLSeconds = 3600
	}
	if config.TTLSeconds < 1 || config.TTLSeconds > 43200 {
		return nil, fmt.Errorf("qca: TTLSeconds must be between 1 and 43200")
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.Timeout < 0 {
		return nil, fmt.Errorf("qca: exchange timeout must be positive")
	}
	client := http.Client{}
	if config.HTTPClient != nil {
		client = *config.HTTPClient
	}
	client.Jar = nil
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &ServiceAccountCredential{key: config.SAKey, region: config.Region, scope: strings.Join(scopes, " "), ttl: config.TTLSeconds,
		timeout: config.Timeout, client: &client, now: time.Now}, nil
}
func ServiceAccountCredentialFromEnv(config ServiceAccountCredentialConfig) (*ServiceAccountCredential, error) {
	config.SAKey = os.Getenv("QODER_SA_KEY")
	return NewServiceAccountCredential(config)
}
func (c *ServiceAccountCredential) String() string {
	return "ServiceAccountCredential(region=" + c.region + ", credentials=redacted)"
}
func (c *ServiceAccountCredential) GoString() string { return c.String() }
func (c *ServiceAccountCredential) validateRegion(host string) error {
	region := ""
	if host == "api.qoder.com" {
		region = "global"
	}
	if host == "api.qoder.com.cn" {
		region = "cn"
	}
	if region != "" && region != c.region {
		return fmt.Errorf("qca: SA credential region does not match base URL")
	}
	return nil
}

// Token permits each waiter to cancel independently; exchange has its own
// bounded context and may finish/cache after its first caller cancels.
func (c *ServiceAccountCredential) Token(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	c.mu.Lock()
	now := c.now()
	if now.Before(c.retryAt) {
		token, err := c.fallback()
		c.mu.Unlock()
		return token, err
	}
	if c.token != "" && now.Before(c.refreshAt) {
		token := c.token
		c.mu.Unlock()
		return token, nil
	}
	flight := c.inflight
	if flight == nil {
		flight = &tokenExchange{done: make(chan struct{})}
		c.inflight = flight
		go c.exchange(flight)
	}
	c.mu.Unlock()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-flight.done:
		if flight.token != "" && !c.now().Before(flight.safeUntil) {
			return "", fmt.Errorf("qca: exchanged token is no longer safe to use")
		}
		return flight.token, flight.err
	}
}
func (c *ServiceAccountCredential) fallback() (string, error) {
	if c.transient && c.token != "" && c.now().Before(c.safeUntil) {
		return c.token, nil
	}
	return "", c.failure
}
func (c *ServiceAccountCredential) exchange(flight *tokenExchange) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	started := c.now()
	c.mu.Lock()
	previousToken := c.token
	c.mu.Unlock()
	token, ttl, cooldown, transient, err := c.fetch(ctx, previousToken)
	c.mu.Lock()
	defer c.mu.Unlock()
	if err == nil {
		lifetime := time.Duration(ttl) * time.Second
		safeUntil := started.Add(lifetime - min(5*time.Second, lifetime/20))
		if !c.now().Before(safeUntil) {
			err = fmt.Errorf("qca: exchanged token has insufficient remaining lifetime")
		} else {
			c.token, c.safeUntil = token, safeUntil
			c.refreshAt = started.Add(lifetime - min(time.Minute, lifetime/10))
			c.retryAt, c.failure = time.Time{}, nil
			flight.token, flight.safeUntil = token, safeUntil
		}
	}
	if err != nil {
		c.failure, c.transient = err, transient
		c.retryAt = c.now().Add(cooldown)
		flight.token, flight.err = c.fallback()
		flight.safeUntil = c.safeUntil
	}
	c.inflight = nil
	close(flight.done)
}
func (c *ServiceAccountCredential) fetch(ctx context.Context, previousToken string) (token string, ttl int, cooldown time.Duration, transient bool, err error) {
	cooldown = 5 * time.Second
	host := "https://openapi.qoder.sh"
	if c.region == "cn" {
		host = "https://openapi.qoder.com.cn"
	}
	body, _ := json.Marshal(map[string]any{"grant_type": "client_credentials", "audience": "qoder", "scope": c.scope, "ttl_seconds": c.ttl})
	req, _ := http.NewRequestWithContext(ctx, "POST", host+"/api/v1/serviceToken/exchange", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return "", 0, cooldown, true, safeExchangeConnectionError(err)
	}
	defer res.Body.Close()
	content, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return "", 0, cooldown, true, safeExchangeConnectionError(err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		parsed := apierror.Error{}
		_ = json.Unmarshal(content, &parsed)
		redact := func(v string) string {
			for _, secret := range []string{c.key, previousToken} {
				if secret != "" {
					v = strings.ReplaceAll(v, secret, "[redacted]")
				}
			}
			if len(v) > 256 {
				v = v[:256]
			}
			return v
		}
		requestID := parsed.RequestID
		if requestID == "" {
			requestID = res.Header.Get("X-Request-ID")
		}
		if requestID == "" {
			requestID = res.Header.Get("Request-ID")
		}
		failure := &apierror.Error{StatusCode: res.StatusCode, Code: redact(parsed.Code), RequestID: redact(requestID), Message: "service account token exchange failed"}
		if res.StatusCode == 429 {
			if delay, ok := parseRetryAfterHeader(res); ok && delay > cooldown {
				cooldown = delay
			}
		}
		return "", 0, cooldown, res.StatusCode == 429 || res.StatusCode >= 500, failure
	}
	var data struct {
		AccessToken string          `json:"access_token"`
		TokenType   string          `json:"token_type"`
		ExpiresIn   json.RawMessage `json:"expires_in"`
		ExpiresAt   json.RawMessage `json:"expires_at"`
	}
	invalid := fmt.Errorf("qca: invalid service account token exchange response")
	if json.Unmarshal(content, &data) != nil || strings.TrimSpace(data.AccessToken) == "" || !strings.EqualFold(data.TokenType, "Bearer") {
		return "", 0, cooldown, false, invalid
	}
	value := string(data.ExpiresIn)
	if strings.HasPrefix(value, "\"") {
		if json.Unmarshal(data.ExpiresIn, &value) != nil {
			return "", 0, cooldown, false, invalid
		}
	}
	if value == "" || strings.IndexFunc(value, func(r rune) bool { return r < '0' || r > '9' }) >= 0 {
		return "", 0, cooldown, false, invalid
	}
	ttl, err = strconv.Atoi(value)
	if err != nil || ttl <= 0 || ttl > 43200 {
		return "", 0, cooldown, false, invalid
	}
	if len(data.ExpiresAt) > 0 {
		var expiresAt string
		if json.Unmarshal(data.ExpiresAt, &expiresAt) != nil {
			return "", 0, cooldown, false, invalid
		}
		if _, err := time.Parse(time.RFC3339Nano, expiresAt); err != nil {
			return "", 0, cooldown, false, invalid
		}
	}
	return data.AccessToken, ttl, cooldown, false, nil
}

func safeExchangeConnectionError(err error) error {
	// Do not retain a transport error that could contain the request's credentials.
	for _, cause := range []error{context.DeadlineExceeded, context.Canceled} {
		if errors.Is(err, cause) {
			return fmt.Errorf("qca: service account token exchange failed: %w", cause)
		}
	}
	return fmt.Errorf("qca: service account token exchange connection failed")
}
