package forward_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestForwardExpansionAPIContracts(t *testing.T) {
	cases := []struct {
		operation string
		call      func(forward.Client) (any, error)
	}{
		{"updateVaultCredential", func(c forward.Client) (any, error) {
			return c.Vaults.Credentials.Update(context.Background(), pathSegment, pathSegment, contractParams[forward.VaultCredentialUpdateParams](t, "updateVaultCredential"))
		}},
		{"listIdentityUsage", func(c forward.Client) (any, error) {
			return c.Usage.ListIdentities(context.Background(), contractParams[forward.UsageListParams](t, "listIdentityUsage"))
		}},
		{"listTemplateUsage", func(c forward.Client) (any, error) {
			return c.Usage.ListTemplates(context.Background(), contractParams[forward.UsageListParams](t, "listTemplateUsage"))
		}},
	}
	for _, c := range cases {
		t.Run(c.operation, func(t *testing.T) {
			for _, status := range []int{200, 400} {
				client := contractClient(t, c.operation, status)
				res, err := c.call(client)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			}
		})
	}
}

func TestForwardUsageHourlyPagination(t *testing.T) {
	for _, backward := range []bool{false, true} {
		t.Run(map[bool]string{false: "forward", true: "backward"}[backward], func(t *testing.T) {
			calls := 0
			client := testClient(func(req *http.Request) (*http.Response, error) {
				calls++
				q := req.URL.Query()
				if q.Get("start_at") != "2026-09-14T09:00:00" || q.Get("end_at") != "2026-09-14T12:00:00" || !reflect.DeepEqual(q["identity_ids"], []string{"idn_one", "idn_two"}) || !reflect.DeepEqual(q["template_ids"], []string{"tmpl_one", "tmpl_two"}) {
					t.Fatalf("query=%v", q)
				}
				if q.Has("start_time") || q.Has("end_time") || q.Has("identity_ids[]") {
					t.Fatal("unexpected legacy parameters")
				}
				if calls == 2 {
					key, next := "after_id", "last"
					if backward {
						key, next = "before_id", "first"
					}
					if q.Get(key) != next {
						t.Fatalf("cursor=%v", q)
					}
				}
				return reply(req, 200, `{"type":"identity_usage.list","start_at":"2026-09-14T09:00:00","end_at":"2026-09-14T12:00:00","first_id":"first","last_id":"last","has_more":true,"data":[{"type":"identity_usage","identity_id":"idn_one","session_count":3,"active_seconds":720.222,"credits":3.1}]}`), nil
			})
			params := forward.UsageListParams{StartAt: "2026-09-14T09:00:00", EndAt: "2026-09-14T12:00:00", IdentityIDs: []string{"idn_one", "idn_two"}, TemplateIDs: []string{"tmpl_one", "tmpl_two"}}
			if backward {
				params.BeforeID = forward.String("initial")
			}
			page, err := client.Usage.ListIdentities(context.Background(), params)
			if err != nil {
				t.Fatal(err)
			}
			if page.Data[0].ActiveSeconds != 720.222 {
				t.Fatal("fractional seconds lost")
			}
			if page.JSON.ExtraFields["start_at"].Raw() != `"2026-09-14T09:00:00"` {
				t.Fatal("window metadata lost")
			}
			if _, err = page.GetNextPage(); err != nil {
				t.Fatal(err)
			}
			if calls != 2 {
				t.Fatal(calls)
			}
		})
	}
}

func TestForwardCredentialUpdateMergePatch(t *testing.T) {
	for _, clearAll := range []bool{false, true} {
		t.Run(map[bool]string{false: "remove_key", true: "clear_metadata"}[clearAll], func(t *testing.T) {
			params := forward.VaultCredentialUpdateParams{IdentityID: forward.String("idn_one"), Auth: forward.VaultCredentialUpdateAuthUnionParam{OfMCPOAuth: &forward.MCPOAuthUpdateParam{Type: "mcp_oauth", ExpiresAt: param.Null[string](), Refresh: forward.MCPOAuthRefreshUpdateParam{Scope: param.Null[string](), RefreshToken: forward.String("test-refresh")}}}, Metadata: map[string]any{"remove": nil}}
			wantMetadata := any(map[string]any{"remove": nil})
			if clearAll {
				params.Metadata = param.NullMap[map[string]any]()
				wantMetadata = nil
			}
			client := testClient(func(req *http.Request) (*http.Response, error) {
				if req.URL.Query().Get("identity_id") != "idn_one" {
					t.Fatal(req.URL)
				}
				b, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}
				var got map[string]any
				if err = json.Unmarshal(b, &got); err != nil {
					t.Fatal(err)
				}
				want := map[string]any{"auth": map[string]any{"type": "mcp_oauth", "expires_at": nil, "refresh": map[string]any{"scope": nil, "refresh_token": "test-refresh"}}, "metadata": wantMetadata}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("body=%s", b)
				}
				return reply(req, 200, string(contract(t, "updateVaultCredential").Response)), nil
			})
			if _, err := client.Vaults.Credentials.Update(context.Background(), "vault_one", "cred_one", params); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestForwardCredentialUpdateDoesNotRetry(t *testing.T) {
	for _, status := range []int{0, 429, 503} {
		t.Run(statusName(status), func(t *testing.T) {
			calls := 0
			client := testClient(func(req *http.Request) (*http.Response, error) {
				calls++
				if status == 0 {
					return nil, errors.New("connection failed")
				}
				res := reply(req, status, `{"error":{"type":"api_error","message":"rotation uncertain"}}`)
				res.Header.Set("Retry-After-Ms", "1")
				res.Header.Set("X-Should-Retry", "true")
				return res, nil
			}, option.WithMaxRetries(3))
			_, err := client.Vaults.Credentials.Update(context.Background(), "vault_one", "cred_one", forward.VaultCredentialUpdateParams{Auth: forward.VaultCredentialUpdateAuthUnionParam{OfStaticBearer: &forward.StaticBearerUpdateParam{Type: "static_bearer", Token: forward.String("test-secret")}}}, option.WithHeader("Idempotency-Key", "caller-key"), option.WithMaxRetries(3))
			if err == nil || calls != 1 {
				t.Fatalf("calls=%d err=%v", calls, err)
			}
		})
	}
}
