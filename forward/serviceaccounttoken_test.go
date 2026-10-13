package forward_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestServiceAccountTokenAPIContracts(t *testing.T) {
	for _, id := range []string{"createServiceAccountToken", "createIdentityServiceAccountToken", "revokeServiceAccountToken"} {
		t.Run(id, func(t *testing.T) {
			for _, status := range []int{200, 400} {
				client := contractClient(t, id, status)
				var res any
				var err error
				switch id {
				case "createServiceAccountToken":
					res, err = client.ServiceAccountTokens.New(context.Background(), contractParams[forward.ServiceAccountTokenNewParams](t, id))
				case "createIdentityServiceAccountToken":
					res, err = client.Identities.ServiceAccountTokens.New(context.Background(), pathSegment, contractParams[forward.ServiceAccountTokenNewParams](t, id))
				case "revokeServiceAccountToken":
					res, err = client.ServiceAccountTokens.Revoke(context.Background(), pathSegment, contractParams[forward.ServiceAccountTokenRevokeParams](t, id))
				}
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			}
		})
	}
}

type tokenManagementCredential struct{ calls int }

func (c *tokenManagementCredential) Token(context.Context) (string, error) {
	c.calls++
	return "business-token", nil
}

func TestServiceAccountTokenWithoutBusinessCredentials(t *testing.T) {
	t.Setenv("QODER_PAT", "")
	t.Setenv("QODER_ACCESS_TOKEN", "")
	client := forward.NewClient(option.WithAuthMode("sa"), option.WithBaseURL("https://api.test/api/v1/forward"), option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("Authorization") != "Bearer explicit-sa-key" {
			t.Fatal("incorrect Key authentication")
		}
		return reply(req, 201, string(contract(t, "createServiceAccountToken").Response)), nil
	})}))
	if _, err := client.ServiceAccountTokens.New(context.Background(), forward.ServiceAccountTokenNewParams{SAKey: "explicit-sa-key"}); err != nil {
		t.Fatal(err)
	}
}

func TestServiceAccountTokenKeyIsRequestScoped(t *testing.T) {
	t.Setenv("QODER_PAT", "ambient-pat")
	t.Setenv("QODER_ACCESS_TOKEN", "ambient-sat")
	for _, dynamic := range []bool{false, true} {
		credential := &tokenManagementCredential{}
		var requests []*http.Request
		var bodies []map[string]any
		opts := []option.RequestOption{option.WithPAT("business-token")}
		if dynamic {
			opts = []option.RequestOption{option.WithCredential(credential)}
		}
		opts = append(opts, option.WithBaseURL("https://api.test/api/v1/forward"), option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests = append(requests, req)
			if req.Method == "GET" {
				return reply(req, 200, `{"data":[]}`), nil
			}
			body := map[string]any{}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			bodies = append(bodies, body)
			if strings.HasSuffix(req.URL.Path, "/revoke") {
				return reply(req, 200, `{"revoked":true}`), nil
			}
			return reply(req, 201, string(contract(t, "createServiceAccountToken").Response)), nil
		})}))
		client := forward.NewClient(opts...)
		token, err := client.ServiceAccountTokens.New(context.Background(), forward.ServiceAccountTokenNewParams{SAKey: "explicit-sa-key"}, option.WithHeader("authorization", "Bearer wrong"))
		if err != nil {
			t.Fatal(err)
		}
		if credential.calls != 0 {
			t.Fatal("token management called the business credential")
		}
		if _, err = client.Identities.ServiceAccountTokens.New(context.Background(), "idn_one", forward.ServiceAccountTokenNewParams{SAKey: "explicit-sa-key"}); err != nil {
			t.Fatal(err)
		}
		if _, err = client.ServiceAccountTokens.Revoke(context.Background(), token.AuthTokenID, forward.ServiceAccountTokenRevokeParams{SAKey: "explicit-sa-key"}); err != nil {
			t.Fatal(err)
		}
		if _, err = client.Models.List(context.Background()); err != nil {
			t.Fatal(err)
		}
		for i, req := range requests {
			want := "Bearer explicit-sa-key"
			if i == 3 {
				want = "Bearer business-token"
			}
			if req.Header.Get("Authorization") != want {
				t.Fatal("incorrect request authentication")
			}
		}
		for _, body := range bodies {
			if len(body) != 0 {
				t.Fatalf("omitted parameters must produce {}: %v", body)
			}
		}
		if dynamic && credential.calls != 1 {
			t.Fatal("business credential was not retained")
		}
		if token.AccessToken != "sat-contract-secret" || token.ExpiresAt.IsZero() || token.JSON.IdentityID.Valid() {
			t.Fatal("admin token fields not decoded")
		}
		if strings.Contains(fmt.Sprintf("%+v %#v", token, token), token.AccessToken) {
			t.Fatal("formatting exposed SAT")
		}
	}
}

func TestServiceAccountTokenNeverRetries(t *testing.T) {
	for _, method := range []string{"admin", "identity", "revoke"} {
		for _, status := range []int{429, 503, 0} {
			t.Run(fmt.Sprintf("%s/%d", method, status), func(t *testing.T) {
				calls := 0
				cause := errors.New("connection interrupted")
				client := testClient(func(req *http.Request) (*http.Response, error) {
					calls++
					if status == 0 {
						return nil, cause
					}
					response := reply(req, status, `{"error":{"type":"api_error","message":"issuance uncertain"}}`)
					response.Header.Set("x-should-retry", "true")
					response.Header.Set("retry-after-ms", "1")
					return response, nil
				}, option.WithMaxRetries(3))
				opts := []option.RequestOption{option.WithMaxRetries(4), option.WithHeader("Idempotency-Key", "unsupported-key")}
				var err error
				switch method {
				case "admin":
					_, err = client.ServiceAccountTokens.New(context.Background(), forward.ServiceAccountTokenNewParams{SAKey: "sa-key"}, opts...)
				case "identity":
					_, err = client.Identities.ServiceAccountTokens.New(context.Background(), "idn_one", forward.ServiceAccountTokenNewParams{SAKey: "sa-key"}, opts...)
				case "revoke":
					_, err = client.ServiceAccountTokens.Revoke(context.Background(), "auth_one", forward.ServiceAccountTokenRevokeParams{SAKey: "sa-key"}, opts...)
				}
				if err == nil || calls != 1 {
					t.Fatalf("err=%v calls=%d", err, calls)
				}
				if status == 0 && !errors.Is(err, cause) {
					t.Fatal("transport cause lost")
				}
			})
		}
	}
}

func TestServiceAccountTokenMissingKeyFailsBeforeTransport(t *testing.T) {
	client := testClient(func(*http.Request) (*http.Response, error) { t.Fatal("invalid Key reached transport"); return nil, nil })
	for _, key := range []string{"", " ", "key\nInjected:value", "key\t", "key\x00"} {
		if _, err := client.ServiceAccountTokens.New(context.Background(), forward.ServiceAccountTokenNewParams{SAKey: key}); err == nil {
			t.Fatal("invalid Key accepted")
		}
		if _, err := client.Identities.ServiceAccountTokens.New(context.Background(), "idn_one", forward.ServiceAccountTokenNewParams{SAKey: key}); err == nil {
			t.Fatal("invalid Key accepted")
		}
		if _, err := client.ServiceAccountTokens.Revoke(context.Background(), "auth_one", forward.ServiceAccountTokenRevokeParams{SAKey: key}); err == nil {
			t.Fatal("invalid Key accepted")
		}
	}
}

func TestServiceAccountTokenJSONDoesNotContainKey(t *testing.T) {
	params := forward.ServiceAccountTokenNewParams{SAKey: "secret-sa-key", TTLSeconds: forward.Int(43200), Metadata: map[string]any{"purpose": "audit"}}
	b, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err = json.NewDecoder(strings.NewReader(string(b))).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(body, map[string]any{"ttl_seconds": float64(43200), "metadata": map[string]any{"purpose": "audit"}}) {
		t.Fatal("incorrect issuance JSON")
	}
	if strings.Contains(fmt.Sprintf("%+v %#v", params, params), params.SAKey) {
		t.Fatal("formatting exposed SA Key")
	}
	// The request body must remain {} when only authentication is configured.
	b, err = json.Marshal(forward.ServiceAccountTokenRevokeParams{SAKey: "secret-sa-key"})
	if err != nil || string(b) != "{}" {
		t.Fatalf("body=%s err=%v", b, err)
	}
}
