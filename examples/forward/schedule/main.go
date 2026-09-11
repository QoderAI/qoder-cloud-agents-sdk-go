package main

import (
	"context"
	"fmt"
	"os"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/forwardutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	config, err := live.Load("forward")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	s := &forwardutil.Example{Client: forward.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "schedule",
		Description: "创建手动 Schedule，触发一次运行并检查助手回复。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

// run triggers one manual run, waits for completion and verifies its Session output.
func run(ctx context.Context, r *live.Run, s *forwardutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	identity, err := s.Identity(ctx, r)
	if err != nil {
		return err
	}
	template, err := s.Template(ctx, r, forward.TemplateNewParams{EnvironmentID: env})
	if err != nil {
		return err
	}
	marker := live.Marker()
	r.Log("info", "本次运行将发送以下初始消息：")
	r.Message("user", "Reply with exactly "+marker)
	r.Step("创建只手动触发的 Schedule")
	schedule, err := s.Client.Schedules.New(ctx, forward.ScheduleNewParams{IdentityID: identity, TemplateID: template, EnvironmentID: env, Name: live.Name("schedule"), InitialEvents: []map[string]any{{"type": "user.message", "content": "Reply with exactly " + marker}}, TriggerPolicy: map[string]any{"type": "manual"}, Execution: map[string]any{"max_attempts": 1, "max_concurrent_runs": 1}})
	if err != nil {
		return err
	}
	r.Track("schedule", schedule.ID, func(ctx context.Context) error {
		_, err := s.Client.Schedules.Archive(ctx, schedule.ID, forward.ScheduleArchiveParams{})
		return err
	})
	r.Step("手动触发一次 Schedule 运行")
	run, err := s.Client.Schedules.Run(ctx, schedule.ID, forward.ScheduleRunParams{IdempotencyKey: forward.String(live.Name("run"))})
	if err != nil {
		return err
	}
	runID := run.ID
	r.Track("schedule_run", runID, func(ctx context.Context) error {
		for {
			current, err := s.Client.ScheduleRuns.Get(ctx, runID, forward.ScheduleRunGetParams{IdentityID: forward.String(identity)})
			if err != nil {
				return err
			}
			if current.SessionID != "" {
				r.Log("info", "清理运行关联的会话："+current.SessionID)
				return s.FinishSession(ctx, current.SessionID)
			}
			if current.Status == "failed" || current.Status == "skipped" {
				return nil
			}
			if err = live.Pause(ctx); err != nil {
				return fmt.Errorf("run %s has no session for cleanup: %w", runID, err)
			}
		}
	})
	stopWaiting := r.Wait("等待 Schedule Run 完成")
	defer stopWaiting()
	for {
		run, err = s.Client.ScheduleRuns.Get(ctx, runID, forward.ScheduleRunGetParams{IdentityID: forward.String(identity)})
		if err != nil {
			return err
		}
		if run.Status == "completed" {
			break
		}
		if run.Status == "failed" || run.Status == "skipped" {
			return fmt.Errorf("run=%s status=%s", runID, run.Status)
		}
		if err = live.Pause(ctx); err != nil {
			return fmt.Errorf("run=%s status=%s: %w", runID, run.Status, err)
		}
	}
	stopWaiting()
	if run.SessionID == "" {
		return fmt.Errorf("completed run %s has no session", runID)
	}
	return s.WaitTurn(ctx, r, run.SessionID, "", []string{marker}, false, false)
}
