package forward_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

// Official Forward ownership contracts (2026-10-20), independent of inventories.
func TestIdentityOwnershipContracts(t *testing.T) {
	ctx := context.Background()
	for _, scope := range []struct {
		name     string
		identity param.Opt[string]
	}{{"PAT identity", forward.String("idn_ /?&")}, {"explicit empty", forward.String("")}, {"omitted", param.Opt[string]{}}} {
		t.Run(scope.name, func(t *testing.T) {
			identity := scope.identity
			var opts []option.RequestOption
			if identity.Valid() {
				opts = append(opts, forward.WithIdentityID(identity.Value))
			}
			cases := []struct {
				method, path string
				call         func(forward.Client) error
			}{
				{"GET", "environments", func(c forward.Client) error {
					_, err := c.Environments.List(ctx, forward.EnvironmentListParams{IdentityID: identity})
					return err
				}},
				{"POST", "environments", func(c forward.Client) error {
					_, err := c.Environments.New(ctx, forward.EnvironmentNewParams{Name: "env", IdentityID: identity})
					return err
				}},
				{"GET", "environments/env", func(c forward.Client) error { _, err := c.Environments.Get(ctx, "env", opts...); return err }},
				{"POST", "environments/env", func(c forward.Client) error {
					_, err := c.Environments.Update(ctx, "env", forward.EnvironmentUpdateParams{Name: forward.String("renamed"), IdentityID: identity})
					return err
				}},
				{"POST", "environments/env/archive", func(c forward.Client) error { _, err := c.Environments.Archive(ctx, "env", opts...); return err }},
				{"DELETE", "environments/env", func(c forward.Client) error { return c.Environments.Delete(ctx, "env", opts...) }},
				{"GET", "skills", func(c forward.Client) error {
					_, err := c.Skills.List(ctx, forward.SkillListParams{IdentityID: identity})
					return err
				}},
				{"POST", "skills", func(c forward.Client) error {
					_, err := c.Skills.New(ctx, forward.SkillNewParams{Files: []io.Reader{strings.NewReader("# skill")}, IdentityID: identity})
					return err
				}},
				{"GET", "skills/skill", func(c forward.Client) error {
					_, err := c.Skills.Get(ctx, "skill", forward.SkillGetParams{IdentityID: identity})
					return err
				}},
				{"PUT", "skills/skill", func(c forward.Client) error {
					_, err := c.Skills.Update(ctx, "skill", forward.SkillUpdateParams{Description: forward.String("desc"), IdentityID: identity})
					return err
				}},
				{"DELETE", "skills/skill", func(c forward.Client) error { return c.Skills.Delete(ctx, "skill", opts...) }},
				{"GET", "skills/skill/versions", func(c forward.Client) error {
					_, err := c.Skills.Versions.List(ctx, "skill", forward.SkillVersionListParams{IdentityID: identity})
					return err
				}},
				{"POST", "skills/skill/versions", func(c forward.Client) error {
					_, err := c.Skills.Versions.New(ctx, "skill", forward.SkillVersionNewParams{Files: []io.Reader{strings.NewReader("# skill")}, IdentityID: identity})
					return err
				}},
				{"GET", "skills/skill/versions/1", func(c forward.Client) error { _, err := c.Skills.Versions.Get(ctx, "skill", "1", opts...); return err }},
				{"DELETE", "skills/skill/versions/1", func(c forward.Client) error {
					_, err := c.Skills.Versions.Delete(ctx, "skill", "1", opts...)
					return err
				}},
				{"GET", "skills/skill/versions/1/content", func(c forward.Client) error {
					res, err := c.Skills.Versions.Download(ctx, "skill", "1", opts...)
					if err == nil {
						res.Body.Close()
					}
					return err
				}},
				{"GET", "vaults", func(c forward.Client) error {
					_, err := c.Vaults.List(ctx, forward.VaultListParams{IdentityID: identity})
					return err
				}},
				{"POST", "vaults", func(c forward.Client) error {
					_, err := c.Vaults.New(ctx, forward.VaultNewParams{DisplayName: "vault", IdentityID: identity})
					return err
				}},
				{"GET", "vaults/vault", func(c forward.Client) error { _, err := c.Vaults.Get(ctx, "vault", opts...); return err }},
				{"DELETE", "vaults/vault", func(c forward.Client) error { return c.Vaults.Delete(ctx, "vault", opts...) }},
				{"GET", "vaults/vault/credentials", func(c forward.Client) error {
					_, err := c.Vaults.Credentials.List(ctx, "vault", forward.VaultCredentialListParams{IdentityID: identity})
					return err
				}},
				{"POST", "vaults/vault/credentials", func(c forward.Client) error {
					_, err := c.Vaults.Credentials.New(ctx, "vault", forward.VaultCredentialNewParams{Auth: map[string]any{"type": "static_bearer", "token": "secret"}, IdentityID: identity})
					return err
				}},
				{"GET", "vaults/vault/credentials/cred", func(c forward.Client) error {
					_, err := c.Vaults.Credentials.Get(ctx, "vault", "cred", opts...)
					return err
				}},
				{"POST", "vaults/vault/credentials/cred", func(c forward.Client) error {
					_, err := c.Vaults.Credentials.Update(ctx, "vault", "cred", forward.VaultCredentialUpdateParams{Metadata: map[string]any{"note": "rotation"}, IdentityID: identity})
					return err
				}},
				{"DELETE", "vaults/vault/credentials/cred", func(c forward.Client) error { return c.Vaults.Credentials.Delete(ctx, "vault", "cred", opts...) }},
				{"GET", "files", func(c forward.Client) error {
					_, err := c.Files.List(ctx, forward.FileListParams{IdentityID: identity})
					return err
				}},
				{"POST", "files", func(c forward.Client) error {
					_, err := c.Files.Upload(ctx, forward.FileUploadParams{File: strings.NewReader("content"), IdentityID: identity})
					return err
				}},
				{"GET", "files/file", func(c forward.Client) error { _, err := c.Files.GetMetadata(ctx, "file", opts...); return err }},
				{"DELETE", "files/file", func(c forward.Client) error { return c.Files.Delete(ctx, "file", opts...) }},
				{"GET", "files/file/content", func(c forward.Client) error {
					res, err := c.Files.Download(ctx, "file", opts...)
					if err == nil {
						res.Body.Close()
					}
					return err
				}},
			}
			for _, tc := range cases {
				t.Run(tc.method+" "+tc.path, func(t *testing.T) {
					calls := 0
					c := testClient(func(r *http.Request) (*http.Response, error) {
						if r.URL.Host == "storage.test" {
							if r.URL.RawQuery != "signature=ok" || r.Header.Get("Authorization") != "" {
								t.Fatal("scope or auth leaked to storage")
							}
							return reply(r, 200, "content"), nil
						}
						calls++
						if r.Method != tc.method || r.URL.Path != "/api/v1/forward/"+tc.path {
							t.Fatalf("route: %s %s", r.Method, r.URL)
						}
						query := r.URL.Query()
						if query.Has("identity_id") != identity.Valid() || query.Get("identity_id") != identity.Value {
							t.Fatalf("scope: %s", r.URL)
						}
						if len(query) > 1 {
							t.Fatalf("body fields leaked to query: %s", r.URL)
						}
						if r.Body != nil {
							body, _ := io.ReadAll(r.Body)
							if strings.Contains(string(body), "identity_id") {
								t.Fatalf("identity leaked into body: %s", body)
							}
						}
						if tc.path == "files/file/content" {
							return reply(r, 200, `{"url":"https://storage.test/content?signature=ok"}`), nil
						}
						return reply(r, 200, `{"id":"resource","data":[],"has_more":false}`), nil
					})
					if err := tc.call(c); err != nil {
						t.Fatal(err)
					}
					if calls != 1 {
						t.Fatalf("API calls: %d", calls)
					}
				})
			}
		})
	}
}

func TestPairChannelByCodeOrID(t *testing.T) {
	for _, byID := range []bool{false, true} {
		t.Run(map[bool]string{false: "code", true: "ID"}[byID], func(t *testing.T) {
			params := forward.ChannelPairingNewParams{IdentityID: "idn_one", TemplateID: "tmpl_one", IdempotencyKey: forward.String("pair-once")}
			expected := map[string]any{"identity_id": "idn_one", "template_id": "tmpl_one"}
			if byID {
				params.ChannelPairingID = "pair_one"
				expected["channel_pairing_id"] = "pair_one"
			} else {
				params.Code = "K7MP92"
				expected["code"] = "K7MP92"
			}
			c := testClient(func(r *http.Request) (*http.Response, error) {
				if r.URL.RawQuery != "" || r.Header.Get("Idempotency-Key") != "pair-once" {
					t.Fatal("pairing query/header")
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if len(body) != len(expected) {
					t.Fatalf("unexpected body: %+v", body)
				}
				for k, v := range expected {
					if body[k] != v {
						t.Fatalf("body: %+v", body)
					}
				}
				return reply(r, 200, `{"id": "pair_one", "type": "channel_pairing", "channel_id": "channel_one", "scope_type": "room", "scope_external_id": "room_one", "scope_display_name": "Team", "identity_id": "idn_one", "template_id": "tmpl_one", "status": "active", "paired_at": "2026-10-20T01:00:00Z"}`), nil
			})
			res, err := c.ChannelPairings.New(context.Background(), params)
			if err != nil {
				t.Fatal(err)
			}
			if res.ScopeType != "room" || res.ScopeExternalID != "room_one" || res.ScopeDisplayName != "Team" {
				t.Fatalf("pairing=%+v", res)
			}
		})
	}
}

func TestIdentityOwnershipPagingAndOverride(t *testing.T) {
	calls := 0
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Query().Get("identity_id") != "idn_override" {
			t.Fatalf("lost scope: %s", r.URL)
		}
		calls++
		if calls == 1 {
			return reply(r, 200, `{"data":[{"id":"env_one"}],"has_more":true,"next_page":"next"}`), nil
		}
		if calls == 2 && r.URL.Query().Get("page") != "next" {
			t.Fatalf("lost cursor: %s", r.URL)
		}
		return reply(r, 200, `{"data":[{"id":"env_two"}],"has_more":false}`), nil
	})
	pager := c.Environments.ListAutoPaging(context.Background(), forward.EnvironmentListParams{IdentityID: forward.String("idn_typed")}, forward.WithIdentityID("idn_override"))
	count := 0
	for pager.Next() {
		count++
	}
	if err := pager.Err(); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("items=%d", count)
	}
	_, err := c.Environments.New(context.Background(), forward.EnvironmentNewParams{Name: "env", IdentityID: forward.String("idn_typed")}, forward.WithIdentityID("idn_override"))
	if err != nil {
		t.Fatal(err)
	}
}
