//go:build live

package forward_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestForwardSharedCasesThroughGoSDKLive(t *testing.T) {
	s := newLiveSuite(t)
	ctx := s.context(t)
	handlers := map[string]func(*testing.T){
		"valid_token": func(t *testing.T) {
			_, err := s.client.Templates.List(ctx, forward.TemplateListParams{Limit: forward.Int(1)})
			liveCheck(t, err)
		},
		"create_identity_minimal": func(t *testing.T) {
			s := newLiveSuite(t, "WRITE")
			identity, err := s.client.Identities.New(ctx, forward.IdentityNewParams{ExternalID: liveName("minimal")})
			liveCheck(t, err)
			s.cleanup(t, "minimal identity", func(ctx context.Context) error { _, err := s.client.Identities.Delete(ctx, identity.ID); return err })
			if !strings.HasPrefix(identity.ID, "idn_") || !identity.Enabled {
				t.Fatalf("unexpected identity: %s", identity.RawJSON())
			}
		},
		"create_identity_metadata": func(t *testing.T) {
			s := newLiveSuite(t, "WRITE")
			identity := s.identity(t)
			if identity.Metadata["suite"] != "sdk-live" {
				t.Fatal("metadata did not round trip")
			}
		},
		"list_templates_limit": func(t *testing.T) {
			page, err := s.client.Templates.List(ctx, forward.TemplateListParams{Limit: forward.Int(1)})
			liveCheck(t, err)
			if len(page.Data) > 1 {
				t.Fatal("limit was ignored")
			}
		},
		"list_templates_after_cursor": func(t *testing.T) {
			page, err := s.client.Templates.List(ctx, forward.TemplateListParams{Limit: forward.Int(1)})
			liveCheck(t, err)
			if page.LastID == "" {
				t.Skip("no template cursor")
			}
			_, err = s.client.Templates.List(ctx, forward.TemplateListParams{Limit: forward.Int(1), AfterID: forward.String(page.LastID)})
			liveCheck(t, err)
		},
		"get_template_not_found": func(t *testing.T) {
			_, err := s.client.Templates.Get(ctx, "tmpl_does_not_exist")
			var apiErr *convention.Error
			if !errors.As(err, &apiErr) || apiErr.StatusCode != 404 {
				t.Fatalf("expected 404, got %v", err)
			}
		},
	}
	for _, invalid := range []bool{false, true} {
		name := "stream_session_events"
		if invalid {
			name += "_invalid_cursor"
		}
		handlers[name] = func(t *testing.T) {
			id := os.Getenv("QODER_FORWARD_LIVE_SESSION_ID")
			if id == "" {
				t.Skip("no existing live session configured")
			}
			ctx, cancel := context.WithTimeout(ctx, s.timeout)
			defer cancel()
			params := forward.SessionEventStreamParams{}
			if invalid {
				params.LastEventID = forward.String("evt_not_in_this_session")
			}
			stream := s.client.Sessions.Events.StreamEvents(ctx, id, params)
			defer stream.Close()
			gotEvent := stream.Next()
			if invalid {
				var apiErr *convention.Error
				if !errors.As(stream.Err(), &apiErr) || apiErr.StatusCode != 404 {
					t.Fatalf("expected invalid-cursor 404, got %v", stream.Err())
				}
			} else if !gotEvent && stream.Err() != nil && !errors.Is(stream.Err(), context.DeadlineExceeded) {
				t.Fatal(stream.Err())
			}
		}
	}
	data, err := os.ReadFile("testdata/sdk-forward-cases.json")
	liveCheck(t, err)
	var manifest struct {
		Cases []struct{ ID, Scenario string }
	}
	liveCheck(t, json.Unmarshal(data, &manifest))
	if len(manifest.Cases) != len(handlers) {
		t.Fatal("shared live scenario inventory drift")
	}
	for _, item := range manifest.Cases {
		handler, ok := handlers[item.Scenario]
		if !ok {
			t.Fatalf("missing scenario %s", item.Scenario)
		}
		t.Run(item.ID, handler)
	}
}

func TestForwardCollectionsLive(t *testing.T) {
	s := newLiveSuite(t)
	ctx := s.context(t)
	_, err := s.client.Models.List(ctx)
	liveCheck(t, err)
	_, err = s.client.Identities.List(ctx, forward.IdentityListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	_, err = s.client.Sessions.List(ctx, forward.SessionListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	_, err = s.client.Schedules.List(ctx, forward.ScheduleListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	_, err = s.client.Batches.List(ctx, forward.BatchListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	_, err = s.client.Channels.List(ctx, forward.ChannelListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	_, err = s.client.Environments.List(ctx, forward.EnvironmentListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	_, err = s.client.Files.List(ctx, forward.FileListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	_, err = s.client.Skills.List(ctx, forward.SkillListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	_, err = s.client.Vaults.List(ctx, forward.VaultListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
	_, err = s.client.MemoryStores.List(ctx, forward.MemoryStoreListParams{Limit: forward.Int(1)})
	liveCheck(t, err)
}
