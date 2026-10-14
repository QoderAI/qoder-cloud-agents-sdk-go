//go:build live

package integration_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
)

// TestDriveLive only clears newly created Identities; no model execution is needed.
func TestDriveLive(t *testing.T) {
	environment := map[string]string{}
	for _, entry := range os.Environ() {
		key, value, _ := strings.Cut(entry, "=")
		environment[key] = value
	}
	values, err := liveEnvironment(os.Getenv("LIVE_ENV_FILE"), environment)
	saCheck(t, err)
	if values["QODER_DRIVE_LIVE_ALLOW_WRITE"] != "true" {
		t.Fatal("QODER_DRIVE_LIVE_ALLOW_WRITE=true is required")
	}
	for _, region := range []string{"cn", "global"} {
		for _, auth := range []string{"pat", "admin_sat", "identity_sat"} {
			t.Run(region+"/"+auth, func(t *testing.T) {
				host := "https://api.qoder.com"
				if region == "cn" {
					host = "https://api.qoder.com.cn"
				}
				opts := []option.RequestOption{option.WithBaseURL(host + "/api/v1/forward/"), option.WithMaxRetries(0), option.WithRequestTimeout(30 * time.Second), testsupport.RequestLog(t)}
				ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
				defer cancel()
				issuer := forward.NewClient(slices.Concat(opts, []option.RequestOption{option.WithAuthMode("sa")})...)
				key := values["QODER_SA_"+strings.ToUpper(region)+"_KEY"]
				issue := func(identityID string) forward.Client {
					if key == "" {
						t.Fatal("region SA key is required")
					}
					params := forward.ServiceAccountTokenNewParams{SAKey: key, TTLSeconds: forward.Int(600)}
					var token *forward.ServiceAccountToken
					var err error
					if identityID == "" {
						token, err = issuer.ServiceAccountTokens.New(ctx, params)
					} else {
						token, err = issuer.Identities.ServiceAccountTokens.New(ctx, identityID, params)
					}
					saCheck(t, err)
					saCleanup(t, "SAT "+token.AuthTokenID, func(ctx context.Context) error {
						_, err := issuer.ServiceAccountTokens.Revoke(ctx, token.AuthTokenID, forward.ServiceAccountTokenRevokeParams{SAKey: key, TargetCredentialID: forward.String(token.CredentialID), TargetExpiresAt: param.NewOpt(token.ExpiresAt)})
						return err
					})
					return forward.NewClient(slices.Concat(opts, []option.RequestOption{option.WithAuthMode("sa"), option.WithAuthToken(token.AccessToken)})...)
				}
				var admin forward.Client
				if auth == "pat" {
					pat := values["QODER_FORWARD_"+strings.ToUpper(region)+"_PAT"]
					if pat == "" {
						t.Fatal("region PAT is required")
					}
					admin = forward.NewClient(slices.Concat(opts, []option.RequestOption{option.WithAuthMode("pat"), option.WithPAT(pat)})...)
				} else {
					admin = issue("")
				}
				identity, err := admin.Identities.New(ctx, forward.IdentityNewParams{ExternalID: "sdk-drive-" + testsupport.Marker(t)})
				saCheck(t, err)
				id := identity.ID
				saCleanup(t, "identity "+id, func(ctx context.Context) error { _, err := admin.Identities.Delete(ctx, id); return err })
				saCleanup(t, "isolated Drive "+id, func(ctx context.Context) error {
					_, err := admin.Drives.Clear(ctx, forward.DriveClearParams{IdentityID: id})
					return err
				})
				runtime := admin
				if auth == "identity_sat" {
					runtime = issue(id)
				}
				page, err := runtime.Drives.ListEntries(ctx, forward.DriveListEntriesParams{IdentityID: id})
				saCheck(t, err)
				if len(page.Entries) != 0 {
					t.Fatal("new Drive is not empty")
				}
				paths := []string{"sdk 目录/a+b.txt", "sdk 目录/zero.txt", "sdk 目录/nested/child.txt"}
				content := "SDK Drive verification"
				storage := &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
				defer storage.CloseIdleConnections()
				transfer := func(method, url string, headers map[string]string, body io.Reader) []byte {
					req, err := http.NewRequestWithContext(ctx, method, url, body)
					if err != nil {
						t.Fatal("invalid storage request")
					}
					for k, v := range headers {
						if strings.EqualFold(k, "Authorization") {
							t.Fatal("unexpected storage Authorization")
						}
						req.Header.Set(k, v)
					}
					response, err := storage.Do(req)
					if err != nil {
						t.Fatal("storage transport failed")
					}
					defer response.Body.Close()
					if response.StatusCode < 200 || response.StatusCode >= 300 {
						t.Fatalf("storage HTTP %d", response.StatusCode)
					}
					data, err := io.ReadAll(io.LimitReader(response.Body, 1024))
					saCheck(t, err)
					return data
				}
				for _, path := range paths {
					signed, err := runtime.Drives.CreateUploadURL(ctx, forward.DriveCreateUploadURLParams{IdentityID: id, Path: path, ContentType: forward.String("text/plain")})
					saCheck(t, err)
					if signed.Method != "PUT" {
						t.Fatal("upload method")
					}
					body := content
					if strings.Contains(path, "zero") {
						body = ""
					}
					transfer(signed.Method, signed.URL, signed.Headers, strings.NewReader(body))
				}
				var entries []forward.DriveEntry
				cursor := ""
				seen := map[string]bool{}
				for {
					params := forward.DriveListEntriesParams{IdentityID: id, Path: forward.String("sdk 目录"), Limit: forward.Int(1)}
					if cursor != "" {
						params.PageToken = forward.String(cursor)
					}
					page, err := runtime.Drives.ListEntries(ctx, params)
					saCheck(t, err)
					entries = append(entries, page.Entries...)
					cursor = page.NextPageToken
					if cursor == "" {
						break
					}
					if seen[cursor] || len(seen) >= 10 {
						t.Fatal("pagination did not finish")
					}
					seen[cursor] = true
				}
				if len(entries) != 3 {
					t.Fatalf("directory has %d entries", len(entries))
				}
				found := map[string]bool{}
				for _, entry := range entries {
					found[entry.Path] = true
					if entry.Path == paths[1] && (!entry.JSON.Size.Valid() || entry.Size != 0) {
						t.Fatal("zero-byte file size")
					}
					if entry.Type == "directory" && entry.JSON.Size.Valid() {
						t.Fatal("directory has file size")
					}
				}
				for _, path := range []string{paths[0], paths[1], "sdk 目录/nested"} {
					if !found[path] {
						t.Fatal("missing directory entry")
					}
				}
				signed, err := runtime.Drives.CreateDownloadURL(ctx, forward.DriveCreateDownloadURLParams{IdentityID: id, Path: paths[0]})
				saCheck(t, err)
				if signed.Method != "GET" || string(transfer(signed.Method, signed.URL, signed.Headers, nil)) != content {
					t.Fatal("download differs")
				}
				_, err = runtime.Drives.DeleteEntry(ctx, forward.DriveDeleteEntryParams{IdentityID: id, Path: "sdk 目录"})
				saTokenDenied(t, err, 409)
				for _, count := range []int64{1, 0} {
					result, err := runtime.Drives.DeleteEntry(ctx, forward.DriveDeleteEntryParams{IdentityID: id, Path: paths[0]})
					saCheck(t, err)
					if result.DeletedCount != count {
						t.Fatal("delete count")
					}
				}
				deleted, err := runtime.Drives.DeleteEntry(ctx, forward.DriveDeleteEntryParams{IdentityID: id, Path: "sdk 目录/nested", Recursive: forward.Bool(true)})
				saCheck(t, err)
				if deleted.DeletedCount != 1 {
					t.Fatal("recursive count")
				}
				for _, count := range []int64{1, 0} {
					result, err := runtime.Drives.Clear(ctx, forward.DriveClearParams{IdentityID: id})
					saCheck(t, err)
					if result.DeletedCount != count {
						t.Fatal("clear count")
					}
				}
				if auth == "identity_sat" {
					other, err := admin.Identities.New(ctx, forward.IdentityNewParams{ExternalID: "sdk-drive-other-" + testsupport.Marker(t)})
					saCheck(t, err)
					saCleanup(t, "other identity "+other.ID, func(ctx context.Context) error { _, err := admin.Identities.Delete(ctx, other.ID); return err })
					_, err = runtime.Drives.ListEntries(ctx, forward.DriveListEntriesParams{IdentityID: other.ID})
					saTokenDenied(t, err, 403)
				} else {
					env, err := admin.Environments.New(ctx, forward.EnvironmentNewParams{Name: "sdk-drive-" + testsupport.Marker(t), Config: map[string]any{"type": "cloud"}})
					saCheck(t, err)
					saCleanup(t, "environment "+env.ID, func(ctx context.Context) error { _, err := admin.Environments.Archive(ctx, env.ID); return err })
					template, err := admin.Templates.New(ctx, forward.TemplateNewParams{Name: "sdk-drive-" + testsupport.Marker(t), Model: forward.ModelConfigUnionParam{OfString: forward.String("auto")}, EnvironmentID: env.ID, ManagedToolConfig: forward.ManagedToolConfigParam{EnabledTools: []string{"drive"}}})
					saCheck(t, err)
					saCleanup(t, "template "+template.ID, func(ctx context.Context) error {
						_, err := admin.Templates.Archive(ctx, template.ID, forward.TemplateArchiveParams{})
						return err
					})
					current, err := admin.Templates.Get(ctx, template.ID)
					saCheck(t, err)
					if !slices.Equal(current.ManagedToolConfig.EnabledTools, []string{"drive"}) {
						t.Fatal("template config readback")
					}
					current, err = admin.Templates.Update(ctx, template.ID, forward.TemplateUpdateParams{ManagedToolConfig: forward.ManagedToolConfigParam{EnabledTools: []string{}}})
					saCheck(t, err)
					if len(current.ManagedToolConfig.EnabledTools) != 0 {
						t.Fatal("template disable")
					}
				}
			})
		}
	}
}
