//go:build live

package forward_test

import (
	"context"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"testing"
)

func TestForwardScheduleE2ELive(t *testing.T) {
	testutil.RequireE2E(t, "FORWARD")
	s := newLiveSuite(t, "WRITE", "EXECUTION")
	ctx, cancel := context.WithTimeout(context.Background(), testutil.ExecutionTimeout(t))
	defer cancel()
	env := s.environment(t)
	identity := s.identity(t)
	template := s.template(t, env.ID)
	marker := testutil.Marker(t)
	schedule, err := s.client.Schedules.New(ctx, forward.ScheduleNewParams{IdentityID: identity.ID, TemplateID: template.ID, EnvironmentID: env.ID, Name: liveName("schedule-e2e"), InitialEvents: []map[string]any{{"type": "user.message", "content": "Reply with exactly " + marker}}, TriggerPolicy: map[string]any{"type": "manual"}, Execution: map[string]any{"max_attempts": 1, "max_concurrent_runs": 1}})
	liveCheck(t, err)
	s.cleanup(t, "schedule "+schedule.ID, func(ctx context.Context) error {
		_, err := s.client.Schedules.Archive(ctx, schedule.ID, forward.ScheduleArchiveParams{})
		return err
	})
	run, err := s.client.Schedules.Run(ctx, schedule.ID, forward.ScheduleRunParams{IdempotencyKey: forward.String(liveName("schedule-run"))})
	liveCheck(t, err)
	runID := run.ID
	s.cleanup(t, "schedule run "+runID, func(ctx context.Context) error { return s.finishScheduleRun(ctx, runID, identity.ID) })
	t.Logf("schedule=%s run=%s", schedule.ID, runID)
	sessionID := ""
	for {
		run, err = s.client.ScheduleRuns.Get(ctx, runID, forward.ScheduleRunGetParams{IdentityID: forward.String(identity.ID)})
		liveCheck(t, err)
		if run.SessionID != "" && sessionID == "" {
			sessionID = run.SessionID
		}
		if run.Status == "completed" {
			break
		}
		if run.Status == "failed" || run.Status == "skipped" {
			t.Fatalf("schedule run=%s status=%s", runID, run.Status)
		}
		if err = testutil.PollPause(ctx); err != nil {
			t.Fatalf("schedule run=%s session=%s status=%s: %v", runID, sessionID, run.Status, err)
		}
	}
	if sessionID == "" {
		t.Fatal("completed schedule has no session")
	}
	s.waitTurn(t, sessionID, "", []string{marker}, false, false)
}

func (s *liveSuite) finishScheduleRun(ctx context.Context, runID, identityID string) error {
	for {
		run, err := s.client.ScheduleRuns.Get(ctx, runID, forward.ScheduleRunGetParams{IdentityID: forward.String(identityID)})
		if err != nil {
			return fmt.Errorf("run=%s cleanup: %w", runID, err)
		}
		if run.SessionID != "" {
			return s.finishSession(ctx, run.SessionID)
		}
		if run.Status == "failed" || run.Status == "skipped" {
			return nil
		}
		if run.Status == "completed" {
			return fmt.Errorf("completed run %s has no session", runID)
		}
		if err = testutil.PollPause(ctx); err != nil {
			return fmt.Errorf("run=%s remains %s without a session: %w", runID, run.Status, err)
		}
	}
}
