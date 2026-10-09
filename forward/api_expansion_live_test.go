//go:build live

package forward_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestForwardUsageHourlyLive(t *testing.T) {
	s := newLiveSuite(t)
	ctx := s.context(t)
	end := time.Now().In(time.FixedZone("Asia/Shanghai", 8*60*60)).Truncate(time.Hour)
	params := forward.UsageListParams{StartAt: end.Add(-24 * time.Hour).Format("2006-01-02T15:00:00"), EndAt: end.Format("2006-01-02T15:00:00"), Limit: forward.Int(2)}
	identities, err := s.client.Usage.ListIdentities(ctx, params)
	liveCheck(t, err)
	checkUsageWindow(t, identities.RawJSON(), "identity_usage.list", params, len(identities.Data))
	for _, row := range identities.Data {
		if row.Type != "identity_usage" || row.IdentityID == "" || !row.JSON.ActiveSeconds.Valid() || row.ActiveSeconds < 0 || row.Credits < 0 || row.SessionCount < 0 {
			t.Fatal("invalid identity usage item")
		}
		params.IdentityIDs = append(params.IdentityIDs, row.IdentityID)
	}
	if identities.HasMore {
		next, err := identities.GetNextPage()
		liveCheck(t, err)
		if next == nil {
			t.Fatal("identity usage has_more has no usable cursor")
		}
		checkUsageWindow(t, next.RawJSON(), "identity_usage.list", params, len(next.Data))
	}
	if len(params.IdentityIDs) > 0 {
		filtered, err := s.client.Usage.ListIdentities(ctx, params)
		liveCheck(t, err)
		for _, row := range filtered.Data {
			found := false
			for _, id := range params.IdentityIDs {
				found = found || row.IdentityID == id
			}
			if !found {
				t.Fatal("identity usage ignored multi-ID filter")
			}
		}
	}
	params.IdentityIDs = nil
	templates, err := s.client.Usage.ListTemplates(ctx, params)
	liveCheck(t, err)
	checkUsageWindow(t, templates.RawJSON(), "template_usage.list", params, len(templates.Data))
	for _, row := range templates.Data {
		if row.Type != "template_usage" || row.TemplateID == "" || !row.JSON.ActiveSeconds.Valid() || row.ActiveSeconds < 0 || row.Credits < 0 || row.SessionCount < 0 || row.ActiveIdentities < 0 {
			t.Fatal("invalid template usage item")
		}
		params.TemplateIDs = append(params.TemplateIDs, row.TemplateID)
	}
	if templates.HasMore {
		next, err := templates.GetNextPage()
		liveCheck(t, err)
		if next == nil {
			t.Fatal("template usage has_more has no usable cursor")
		}
		checkUsageWindow(t, next.RawJSON(), "template_usage.list", params, len(next.Data))
	}
	if len(params.TemplateIDs) > 0 {
		filtered, err := s.client.Usage.ListTemplates(ctx, params)
		liveCheck(t, err)
		for _, row := range filtered.Data {
			found := false
			for _, id := range params.TemplateIDs {
				found = found || row.TemplateID == id
			}
			if !found {
				t.Fatal("template usage ignored multi-ID filter")
			}
		}
	}
	t.Logf("usage identity_rows=%d template_rows=%d; empty pages cover the collection only", len(identities.Data), len(templates.Data))
}

func checkUsageWindow(t *testing.T, raw, kind string, params forward.UsageListParams, count int) {
	t.Helper()
	var window struct {
		Type    string `json:"type"`
		StartAt string `json:"start_at"`
		EndAt   string `json:"end_at"`
	}
	liveCheck(t, json.Unmarshal([]byte(raw), &window))
	if window.Type != kind || window.StartAt != params.StartAt || window.EndAt != params.EndAt || count > 2 {
		t.Fatal("usage response changed the type, window, or requested limit")
	}
}
