//go:build live

package managed_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func (s *managedScenarioSuite) waitTurn(t *testing.T, id, after string, expected []string, tool, streaming bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), testutil.ExecutionTimeout(t))
	defer cancel()
	result := testutil.TurnResult{LastID: after}
	observe := func(raw string) bool {
		t.Helper()
		if err := result.Observe(raw); err != nil {
			t.Fatalf("session=%s: %v", id, err)
		}
		return result.Complete
	}
	if streaming {
		stream := s.client.Sessions.Events.StreamEvents(ctx, id, managed.SessionEventStreamParams{}, option.WithHeader("Last-Event-ID", after), option.WithRequestTimeout(testutil.ExecutionTimeout(t)))
		defer stream.Close()
		for stream.Next() {
			if observe(stream.Current().RawJSON()) {
				break
			}
		}
		if err := stream.Err(); err != nil {
			t.Fatal(testutil.SafeError(err))
		}
	} else {
		for !result.Complete {
			page := s.client.Sessions.Events.ListAutoPaging(ctx, id, managed.SessionEventListParams{AfterID: managed.String(result.LastID), Order: "asc", Limit: managed.Int(100)})
			count := 0
			for page.Next() {
				count++
				if count > 2000 {
					t.Fatal("execution exceeded 2000 events")
				}
				if observe(page.Current().RawJSON()) {
					break
				}
			}
			if err := page.Err(); err != nil {
				t.Fatal(testutil.SafeError(err))
			}
			if result.Complete {
				break
			}
			if err := testutil.PollPause(ctx); err != nil {
				t.Fatal(testutil.SafeError(err))
			}
		}
	}
	if err := result.Verify(expected, tool); err != nil {
		t.Fatalf("session=%s: %v", id, err)
	}
	t.Logf("verified session=%s last_event_id=%s tool_used=%t", id, result.LastID, result.ToolUsed)
}
func (s *managedScenarioSuite) sendTurn(t *testing.T, id, prompt string) string {
	t.Helper()
	ctx, cancel := s.context()
	defer cancel()
	params := liveJSON[managed.SessionEventSendParams](t, map[string]any{"events": []any{map[string]any{"type": "user.message", "content": []any{map[string]any{"type": "text", "text": prompt}}}}})
	res, err := s.client.Sessions.Events.Send(ctx, id, params, option.WithHeader("Idempotency-Key", managedUnique("send")))
	if err != nil {
		t.Fatal(testutil.SafeError(err))
	}
	if res == nil || len(res.Data) != 1 || res.Data[0].ID == "" {
		t.Fatal("send did not return one user event ID")
	}
	return res.Data[0].ID
}
func (s *managedScenarioSuite) cleanupSession(t *testing.T, id string) {
	t.Helper()
	t.Logf("created session=%s", id)
	s.cleanup(t, "session "+id, func(ctx context.Context) error { return s.finishSession(ctx, id) })
}
func (s *managedScenarioSuite) finishSession(ctx context.Context, id string) error {
	session, err := s.client.Sessions.Get(ctx, id, managed.SessionGetParams{})
	if err != nil {
		return err
	}
	if session.Status != "idle" && session.Status != "terminated" {
		interrupt := managed.SessionEventSendParams{Events: []managed.ManagedAgentsEventParamsUnion{{OfUserInterrupt: &managed.ManagedAgentsUserInterruptEventParams{Type: "user.interrupt"}}}}
		if _, err = s.client.Sessions.Events.Send(ctx, id, interrupt); err != nil {
			return err
		}
		for {
			session, err = s.client.Sessions.Get(ctx, id, managed.SessionGetParams{})
			if err != nil {
				return err
			}
			if session.Status == "idle" || session.Status == "terminated" {
				break
			}
			if err = testutil.PollPause(ctx); err != nil {
				return err
			}
		}
	}
	_, err = s.client.Sessions.Delete(ctx, id, managed.SessionDeleteParams{})
	return err
}
func TestManagedExecutionE2ELive(t *testing.T) {
	testutil.RequireE2E(t, "MANAGED")
	s := newManagedScenarioSuite(t)
	s.requireExecution(t)
	ctx, cancel := context.WithTimeout(context.Background(), testutil.ExecutionTimeout(t))
	defer cancel()
	env := s.createEnvironment(t)
	fileToken, envToken, skillToken, memoryToken := testutil.Marker(t), testutil.Marker(t), testutil.Marker(t), testutil.Marker(t)
	file := liveResult(s.client.Files.Upload(ctx, managed.FileUploadParams{File: convention.UploadFile{Name: "sdk-e2e.txt", Reader: strings.NewReader(fileToken)}})).require(t)
	s.cleanup(t, "file "+file.ID, func(ctx context.Context) error {
		_, err := s.client.Files.Delete(ctx, file.ID, managed.FileDeleteParams{})
		return err
	})
	skillName := managedUnique("proof")
	skill := liveResult(s.client.Skills.New(ctx, managed.SkillNewParams{Files: []io.Reader{convention.UploadFile{Name: skillName + "/SKILL.md", Reader: strings.NewReader(fmt.Sprintf("---\nname: %s\ndescription: Provides SDK_E2E_SKILL_TOKEN for SDK verification.\n---\nWhen asked for SDK_E2E_SKILL_TOKEN return exactly: %s\n", skillName, skillToken))}}})).require(t)
	s.cleanup(t, "skill "+skill.ID, func(ctx context.Context) error {
		_, err := s.client.Skills.Delete(ctx, skill.ID, managed.SkillDeleteParams{})
		return err
	})
	store := s.createMemoryStore(t)
	liveResult(s.client.MemoryStores.Memories.New(ctx, store.ID, managed.MemoryStoreMemoryNewParams{Path: "sdk-e2e/proof.md", Content: managed.String("SDK_E2E_MEMORY_TOKEN=" + memoryToken)})).require(t)
	agentParams := liveJSON[managed.AgentNewParams](t, map[string]any{"name": managedUnique("proof"), "model": map[string]any{"id": os.Getenv("QODER_MANAGED_MODEL")}, "system": "Complete the requested SDK verification. Use the available tools to read files, environment variables, skills and memory. Do not guess missing values.", "tools": []any{map[string]any{"type": "agent_toolset_20260401"}}, "skills": []any{map[string]any{"type": "custom", "skill_id": skill.ID, "version": skill.LatestVersionID}}})
	agent := liveResult(s.client.Agents.New(ctx, agentParams)).require(t)
	s.cleanup(t, "agent "+agent.ID, func(ctx context.Context) error {
		_, err := s.client.Agents.Archive(ctx, agent.ID, managed.AgentArchiveParams{})
		return err
	})
	sessionParams := liveJSON[managed.SessionNewParams](t, map[string]any{"agent": agent.ID, "environment_id": env.ID, "environment_variables": map[string]string{"SDK_E2E_VALUE": envToken}, "resources": []any{map[string]any{"type": "file", "file_id": file.ID, "mount_path": "/data/workspace/sdk-e2e.txt"}, map[string]any{"type": "memory_store", "memory_store_id": store.ID, "access": "read_only", "instructions": "Read sdk-e2e/proof.md when asked for SDK_E2E_MEMORY_TOKEN."}}})
	session := liveResult(s.client.Sessions.New(ctx, sessionParams)).require(t)
	s.cleanupSession(t, session.ID)
	t.Logf("model=%s agent=%s file=%s skill=%s memory_store=%s", os.Getenv("QODER_MANAGED_MODEL"), agent.ID, file.ID, skill.ID, store.ID)
	echo := testutil.Marker(t)
	for _, scenario := range []struct {
		name, prompt string
		expected     []string
		tool, stream bool
	}{
		{"completion_and_sse", "Reply with exactly " + echo, []string{echo}, false, true},
		{"file_and_environment_config", "Use tools to read /data/workspace/sdk-e2e.txt and the SDK_E2E_VALUE environment variable. Reply with both exact values.", []string{fileToken, envToken}, true, false},
		{"skill_and_memory", "Use skill " + skillName + " to obtain SDK_E2E_SKILL_TOKEN. Read sdk-e2e/proof.md from the mounted memory store to obtain SDK_E2E_MEMORY_TOKEN. Reply with both exact tokens.", []string{skillToken, memoryToken}, true, false},
	} {
		if !t.Run(scenario.name, func(t *testing.T) {
			after := s.sendTurn(t, session.ID, scenario.prompt)
			s.waitTurn(t, session.ID, after, scenario.expected, scenario.tool, scenario.stream)
		}) {
			return
		}
	}
}
