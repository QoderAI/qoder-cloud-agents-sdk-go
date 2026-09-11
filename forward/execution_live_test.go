//go:build live

package forward_test

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
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func (s *liveSuite) waitTurn(t *testing.T, sessionID, after string, expected []string, tool, streaming bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), testutil.ExecutionTimeout(t))
	defer cancel()
	result := testutil.TurnResult{LastID: after}
	observe := func(raw string) bool {
		t.Helper()
		if err := result.Observe(raw); err != nil {
			t.Fatalf("session=%s: %v", sessionID, err)
		}
		return result.Complete
	}
	if streaming {
		stream := s.client.Sessions.Events.StreamEvents(ctx, sessionID, forward.SessionEventStreamParams{LastEventID: forward.String(after), IncludeToolCalls: forward.Bool(true)}, option.WithRequestTimeout(testutil.ExecutionTimeout(t)))
		defer stream.Close()
		for stream.Next() {
			if observe(stream.Current().RawJSON()) {
				break
			}
		}
		liveCheck(t, stream.Err())
	} else {
		for !result.Complete {
			page := s.client.Sessions.Events.ListAutoPaging(ctx, sessionID, forward.SessionEventListParams{AfterID: forward.String(result.LastID), Order: forward.String("asc"), Limit: forward.Int(100), IncludeToolCalls: forward.Bool(true)})
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
			liveCheck(t, page.Err())
			if result.Complete {
				break
			}
			liveCheck(t, testutil.PollPause(ctx))
		}
	}
	if err := result.Verify(expected, tool); err != nil {
		t.Fatalf("session=%s: %v", sessionID, err)
	}
	t.Logf("verified session=%s last_event_id=%s tool_used=%t", sessionID, result.LastID, result.ToolUsed)
}
func (s *liveSuite) sendTurn(t *testing.T, id, prompt string) string {
	t.Helper()
	res, err := s.client.Sessions.Events.Send(s.context(t), id, forward.SessionEventSendParams{Events: []forward.SessionEventParam{{Type: "user.message", Content: forward.EventContentUnionParam{OfBlocks: []forward.ContentBlockParam{{Type: "text", Text: forward.String(prompt)}}}}}, IdempotencyKey: forward.String(liveName("send"))})
	liveCheck(t, err)
	if res == nil || len(res.Data) != 1 || res.Data[0].ID == "" {
		t.Fatal("send did not return one user event ID")
	}
	return res.Data[0].ID
}
func (s *liveSuite) cleanupSession(t *testing.T, id string) {
	t.Helper()
	t.Logf("created session=%s", id)
	s.cleanup(t, "session "+id, func(ctx context.Context) error { return s.finishSession(ctx, id) })
}
func (s *liveSuite) finishSession(ctx context.Context, id string) error {
	session, err := s.client.Sessions.Get(ctx, id)
	if err != nil {
		return err
	}
	if session.Status != "idle" && session.Status != "terminated" {
		if _, err = s.client.Sessions.Cancel(ctx, id, forward.SessionCancelParams{}); err != nil {
			return err
		}
		for {
			session, err = s.client.Sessions.Get(ctx, id)
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
	_, err = s.client.Sessions.Archive(ctx, id, forward.SessionArchiveParams{})
	return err
}
func TestForwardExecutionE2ELive(t *testing.T) {
	testutil.RequireE2E(t, "FORWARD")
	s := newLiveSuite(t, "WRITE", "EXECUTION")
	ctx := s.context(t)
	env := s.environment(t)
	identity := s.identity(t)
	fileToken, envToken, skillToken, memoryToken := testutil.Marker(t), testutil.Marker(t), testutil.Marker(t), testutil.Marker(t)
	file := s.file(t, "sdk-e2e.txt", "session_resource", fileToken)
	skillName := liveName("proof")
	skill, err := s.client.Skills.New(ctx, forward.SkillNewParams{Files: []io.Reader{convention.UploadFile{Name: skillName + "/SKILL.md", Reader: strings.NewReader(fmt.Sprintf("---\nname: %s\ndescription: Provides the SDK_E2E_SKILL_TOKEN for SDK verification.\n---\nWhen asked for SDK_E2E_SKILL_TOKEN return exactly: %s\n", skillName, skillToken))}}})
	liveCheck(t, err)
	s.cleanup(t, "skill "+skill.ID, func(ctx context.Context) error { return s.client.Skills.Delete(ctx, skill.ID) })
	store, err := s.client.MemoryStores.New(ctx, forward.MemoryStoreNewParams{Name: liveName("proof"), IdempotencyKey: liveName("memory-key")})
	liveCheck(t, err)
	s.cleanup(t, "memory store "+store.ID, func(ctx context.Context) error { _, err := s.client.MemoryStores.Delete(ctx, store.ID); return err })
	_, err = s.client.MemoryStores.Memories.New(ctx, store.ID, forward.MemoryStoreMemoryNewParams{Path: "sdk-e2e/proof.md", Content: "SDK_E2E_MEMORY_TOKEN=" + memoryToken})
	liveCheck(t, err)
	template, err := s.client.Templates.New(ctx, forward.TemplateNewParams{Name: liveName("proof"), EnvironmentID: env.ID, Model: forward.ModelConfigUnionParam{OfString: forward.String(os.Getenv("QODER_FORWARD_MODEL"))}, System: forward.String("Complete the requested SDK verification. Use the available tools to read files, environment variables, skills and memory. Do not guess missing values."), Tools: []forward.ToolParam{{Type: "agent_toolset_20260401"}}, Skills: []forward.SkillBindingParam{{Type: "custom", SkillID: skill.ID, Version: forward.String(skill.LatestVersion)}}, EnvironmentVariables: forward.EnvironmentVariablesUnionParam{OfMap: map[string]any{"SDK_E2E_VALUE": "template-default"}}})
	liveCheck(t, err)
	s.cleanup(t, "template "+template.ID, func(ctx context.Context) error {
		_, err := s.client.Templates.Archive(ctx, template.ID, forward.TemplateArchiveParams{})
		return err
	})
	_, err = s.client.Identities.Configs.Upsert(ctx, identity.ID, template.ID, forward.IdentityConfigUpsertParams{IdentityConfig: forward.IdentityConfigSpecParam{EnvironmentVariables: map[string]forward.EnvironmentVariableOverrideParam{"SDK_E2E_VALUE": {Op: "set", Value: forward.String(envToken)}}}})
	liveCheck(t, err)
	_, err = s.client.Identities.MemoryStores.Mount(ctx, identity.ID, template.ID, forward.IdentityMemoryStoreMountParams{MemoryStoreID: store.ID})
	liveCheck(t, err)
	s.cleanup(t, "memory mount "+store.ID, func(ctx context.Context) error {
		_, err := s.client.Identities.MemoryStores.Detach(ctx, identity.ID, template.ID, store.ID)
		return err
	})
	session, err := s.client.Sessions.New(ctx, forward.SessionNewParams{IdentityID: identity.ID, TemplateID: template.ID, Resources: []forward.SessionResourceSpecParam{{Type: "file", FileID: file.ID, MountPath: forward.String("/data/workspace/sdk-e2e.txt")}}})
	liveCheck(t, err)
	s.cleanupSession(t, session.ID)
	t.Logf("model=%s identity=%s template=%s file=%s skill=%s memory_store=%s", os.Getenv("QODER_FORWARD_MODEL"), identity.ID, template.ID, file.ID, skill.ID, store.ID)
	echo := testutil.Marker(t)
	for _, scenario := range []struct {
		name, prompt string
		expected     []string
		tool, stream bool
	}{
		{"completion_and_sse", "Reply with exactly " + echo, []string{echo}, false, true},
		{"file_and_identity_config", "Use tools to read /data/workspace/sdk-e2e.txt and the SDK_E2E_VALUE environment variable. Reply with both exact values.", []string{fileToken, envToken}, true, false},
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
