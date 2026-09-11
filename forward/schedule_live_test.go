//go:build live

package forward_test

import (
	"context"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestScheduleLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	env := s.environment(t)
	identity := s.identity(t)
	template := s.template(t, env.ID)
	schedule, err := s.client.Schedules.New(ctx, forward.ScheduleNewParams{IdentityID: identity.ID, TemplateID: template.ID, EnvironmentID: env.ID, Name: liveName("schedule"), InitialEvents: []map[string]any{{"type": "user.message", "content": "Reply with SDK-LIVE."}}, TriggerPolicy: map[string]any{"type": "manual"}})
	liveCheck(t, err)
	s.cleanup(t, "schedule", func(ctx context.Context) error {
		got, err := s.client.Schedules.Get(ctx, schedule.ID)
		if err != nil || !got.ArchivedAt.IsZero() {
			return err
		}
		_, err = s.client.Schedules.Archive(ctx, schedule.ID, forward.ScheduleArchiveParams{})
		return err
	})
	got, err := s.client.Schedules.Get(ctx, schedule.ID)
	liveCheck(t, err)
	if got.ID != schedule.ID {
		t.Fatal("schedule did not round trip")
	}
	updated, err := s.client.Schedules.Update(ctx, schedule.ID, forward.ScheduleUpdateParams{Description: forward.String("SDK updated schedule")})
	liveCheck(t, err)
	if updated.Description != "SDK updated schedule" {
		t.Fatal("schedule description not updated")
	}
	paused, err := s.client.Schedules.Pause(ctx, schedule.ID, forward.SchedulePauseParams{})
	liveCheck(t, err)
	if paused.Status != "paused" {
		t.Fatal("schedule not paused")
	}
	active, err := s.client.Schedules.Unpause(ctx, schedule.ID, forward.ScheduleUnpauseParams{})
	liveCheck(t, err)
	if active.Status != "active" {
		t.Fatal("schedule not active")
	}
	_, err = s.client.ScheduleRuns.List(ctx, forward.ScheduleRunListParams{IdentityID: identity.ID, ScheduleID: forward.String(schedule.ID)})
	liveCheck(t, err)
	archived, err := s.client.Schedules.ArchiveMany(ctx, forward.ScheduleArchiveManyParams{
		Scope:       "by_schedule_ids",
		ScheduleIDs: []string{schedule.ID, schedule.ID},
	})
	liveCheck(t, err)
	if archived.ArchivedCount != 1 {
		t.Fatal("archive did not deduplicate IDs")
	}
	got, err = s.client.Schedules.Get(ctx, schedule.ID)
	liveCheck(t, err)
	if got.ArchivedAt.IsZero() {
		t.Fatal("batch archive did not persist")
	}
}
