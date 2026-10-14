package forward_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestDriveAPIContracts(t *testing.T) {
	for _, status := range []int{200, 400} {
		t.Run(statusName(status), func(t *testing.T) {
			ctx := context.Background()
			client0 := contractClient(t, "listDriveEntries", status)
			a, err := client0.Drives.ListEntries(ctx, contractParams[forward.DriveListEntriesParams](t, "listDriveEntries"))
			checkError(t, status, err)
			if status == 200 {
				checkDecoded(t, a)
			}
			client1 := contractClient(t, "createDriveUploadURL", status)
			b, err := client1.Drives.CreateUploadURL(ctx, contractParams[forward.DriveCreateUploadURLParams](t, "createDriveUploadURL"))
			checkError(t, status, err)
			if status == 200 {
				checkDecoded(t, b)
			}
			client2 := contractClient(t, "createDriveDownloadURL", status)
			c, err := client2.Drives.CreateDownloadURL(ctx, contractParams[forward.DriveCreateDownloadURLParams](t, "createDriveDownloadURL"))
			checkError(t, status, err)
			if status == 200 {
				checkDecoded(t, c)
			}
			client3 := contractClient(t, "deleteDriveEntry", status)
			d, err := client3.Drives.DeleteEntry(ctx, contractParams[forward.DriveDeleteEntryParams](t, "deleteDriveEntry"))
			checkError(t, status, err)
			if status == 200 {
				checkDecoded(t, d)
			}
			client4 := contractClient(t, "clearDrive", status)
			e, err := client4.Drives.Clear(ctx, contractParams[forward.DriveClearParams](t, "clearDrive"))
			checkError(t, status, err)
			if status == 200 {
				checkDecoded(t, e)
			}
		})
	}
}

func TestDriveEncodingAndPaginationPreserveScopeAndOrder(t *testing.T) {
	cursor := "opaque/+%=&"
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Query().Get("identity_id") != "idn_test" || r.URL.Query().Get("path") != "目录/a b+~*" || r.URL.Query().Get("extra") != "a/b" {
			t.Fatalf("query: %s", r.URL.RawQuery)
		}
		if !strings.Contains(r.URL.RawQuery, "path=%E7%9B%AE%E5%BD%95/a+b%2B~%2A") || !strings.Contains(r.URL.RawQuery, "extra=a%2Fb") {
			t.Fatalf("wire encoding: %s", r.URL.RawQuery)
		}
		if calls == 1 {
			return reply(r, 200, `{"entries":[{"path":"a.txt","name":"a.txt","type":"file","size":0},{"path":"a","name":"a","type":"directory"}],"next_page_token":"opaque/+%=&"}`), nil
		}
		if r.URL.Query().Get("page_token") != cursor {
			t.Fatal("cursor modified")
		}
		return reply(r, 200, `{"entries":[]}`), nil
	}, option.WithAuthMode("sa"), option.WithAuthToken("test-sat"))
	p := forward.DriveListEntriesParams{IdentityID: "idn_test", Path: forward.String("replaced"), Limit: forward.Int(2)}
	opts := []option.RequestOption{option.WithQuery("path", "目录/a b+~*"), option.WithQuery("extra", "a/b")}
	first, err := client.Drives.ListEntries(context.Background(), p, opts...)
	if err != nil {
		t.Fatal(err)
	}
	if first.Entries[0].Path != "a.txt" || first.Entries[1].Path != "a" {
		t.Fatal("service order changed")
	}
	if !first.Entries[0].JSON.Size.Valid() || first.Entries[1].JSON.Size.Valid() {
		t.Fatal("zero-byte file and directory indistinguishable")
	}
	p.PageToken = forward.String(first.NextPageToken)
	last, err := client.Drives.ListEntries(context.Background(), p, opts...)
	if err != nil {
		t.Fatal(err)
	}
	if last.NextPageToken != "" || len(last.Entries) != 0 || calls != 2 {
		t.Fatal("unexpected final page")
	}
}

func TestDriveClearHasNoBodyAndPartialDeletionIsNotRetried(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != "DELETE" || r.URL.Path != "/api/v1/forward/drives" || r.URL.Query().Get("identity_id") != "idn_test" {
			t.Fatalf("request: %s %s", r.Method, r.URL)
		}
		if r.Body != nil {
			b, _ := io.ReadAll(r.Body)
			if len(b) != 0 {
				t.Fatalf("unexpected body %s", b)
			}
		}
		return reply(r, 503, `{"request_id":"req_partial","error":{"type":"api_error","code":"drive_delete_partial","deleted_count":7,"message":"partial"}}`), nil
	}, option.WithMaxRetries(2))
	_, err := client.Drives.Clear(context.Background(), forward.DriveClearParams{IdentityID: "idn_test"})
	var api *apierror.Error
	if !errors.As(err, &api) || api.Code != "drive_delete_partial" || api.RequestID != "req_partial" || !strings.Contains(api.RawJSON(), `"deleted_count":7`) || calls != 1 {
		t.Fatalf("partial error: %v calls=%d", err, calls)
	}
}

func TestTemplateManagedToolConfigPresence(t *testing.T) {
	for _, tc := range []struct {
		name     string
		config   forward.ManagedToolConfigParam
		expected string
	}{
		{"omitted", forward.ManagedToolConfigParam{}, "{}"},
		{"enabled", forward.ManagedToolConfigParam{EnabledTools: []string{"drive", "schedule"}}, `{"managed_tool_config":{"enabled_tools":["drive","schedule"]}}`},
		{"empty", forward.ManagedToolConfigParam{EnabledTools: []string{}}, `{"managed_tool_config":{"enabled_tools":[]}}`},
		{"null", param.NullStruct[forward.ManagedToolConfigParam](), `{"managed_tool_config":null}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := testClient(func(r *http.Request) (*http.Response, error) {
				var actual, expected any
				if err := json.NewDecoder(r.Body).Decode(&actual); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(tc.expected), &expected); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(actual, expected) {
					t.Fatalf("body=%v expected=%v", actual, expected)
				}
				return reply(r, 200, `{"managed_tool_config":{"enabled_tools":["drive","schedule"]}}`), nil
			})
			result, err := client.Templates.Update(context.Background(), "tmpl_test", forward.TemplateUpdateParams{ManagedToolConfig: tc.config})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result.ManagedToolConfig.EnabledTools, []string{"drive", "schedule"}) {
				t.Fatal("typed response missing")
			}
		})
	}
}

func TestDrivePreservesRejectedPathsAndLimits(t *testing.T) {
	for _, tc := range []struct {
		path  string
		limit int64
		ok    bool
	}{
		{"", 1, true}, {"", 100, true}, {"/absolute", 20, false}, {"../other", 20, false}, {`a\b`, 20, false}, {"a%2Fb", 20, false}, {" a ", 20, false}, {"", 0, false}, {"", 101, false},
	} {
		client := testClient(func(r *http.Request) (*http.Response, error) {
			q := r.URL.Query()
			if q.Get("path") != tc.path || q.Get("limit") != fmt.Sprint(tc.limit) || len(q["identity_id"]) != 1 {
				t.Fatalf("changed query %s", r.URL.RawQuery)
			}
			if r.Body != nil {
				body, _ := io.ReadAll(r.Body)
				if len(body) > 0 {
					t.Fatal("GET body")
				}
			}
			if tc.ok {
				return reply(r, 200, `{"entries":[]}`), nil
			}
			return reply(r, 400, `{"error":{"code":"invalid_argument","message":"invalid Drive query"}}`), nil
		})
		_, err := client.Drives.ListEntries(context.Background(), forward.DriveListEntriesParams{IdentityID: "idn_test", Path: forward.String(tc.path), Limit: forward.Int(tc.limit)})
		if tc.ok {
			if err != nil {
				t.Fatal(err)
			}
		} else {
			var api *apierror.Error
			if !errors.As(err, &api) || api.StatusCode != 400 || api.Code != "invalid_argument" {
				t.Fatalf("wrong error: %v", err)
			}
		}
	}
}
