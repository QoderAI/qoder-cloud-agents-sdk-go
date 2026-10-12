//go:build live

package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func saCheck(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(testsupport.SafeError(err))
	}
}

func saCleanup(t *testing.T, label string, action func(context.Context) error) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := action(ctx); err != nil {
			t.Errorf("cleanup %s: %s", label, testsupport.SafeError(err))
		} else {
			t.Logf("cleanup %s: completed", label)
		}
	})
}

func saFinishSession(ctx context.Context, get func() (string, error), cancel, finish func() error) error {
	status, err := get()
	if err != nil {
		return err
	}
	if status != "idle" && status != "terminated" {
		if err = cancel(); err != nil {
			return err
		}
		for status != "idle" && status != "terminated" {
			if err = testsupport.PollPause(ctx); err != nil {
				return err
			}
			status, err = get()
			if err != nil {
				return err
			}
		}
	}
	return finish()
}

func TestServiceAccountExecutionLive(t *testing.T) {
	saRequireConfiguration(t, true)
	model := os.Getenv("QODER_SA_MODEL")
	if model == "" {
		model = "auto"
	}
	for _, region := range []string{"cn", "global"} {
		t.Run(region, func(t *testing.T) {
			variable, host := "QODER_SA_GLOBAL_KEY", "https://api.qoder.com"
			if region == "cn" {
				variable, host = "QODER_SA_CN_KEY", "https://api.qoder.com.cn"
			}
			key := os.Getenv(variable)
			if key == "" {
				t.Fatalf("%s is required", variable)
			}
			credential, err := convention.NewServiceAccountCredential(convention.ServiceAccountCredentialConfig{SAKey: key, Region: region, Scopes: []string{"qca.access", "forward.access"}})
			saCheck(t, err)
			for _, mode := range []string{"forward", "managed"} {
				t.Run(mode, func(t *testing.T) {
					opts := []option.RequestOption{option.WithCredential(credential), option.WithMaxRetries(0), option.WithRequestTimeout(30 * time.Second), testsupport.RequestLog(t)}
					ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
					defer cancel()
					if mode == "forward" {
						saForwardExecution(t, ctx, forward.NewClient(append(opts, option.WithBaseURL(host+"/api/v1/forward/"))...), model)
					} else {
						saManagedExecution(t, ctx, managed.NewClient(append(opts, option.WithBaseURL(host+"/api/v1/cloud/"))...), model)
					}
				})
			}
		})
	}
}

func saForwardExecution(t *testing.T, ctx context.Context, client forward.Client, model string) {
	environment, err := client.Environments.New(ctx, forward.EnvironmentNewParams{Name: "sdk-go-sa-env-" + testsupport.Marker(t), Config: map[string]any{"type": "cloud"}})
	saCheck(t, err)
	saCleanup(t, "environment "+environment.ID, func(ctx context.Context) error {
		_, err := client.Environments.Archive(ctx, environment.ID)
		return err
	})
	gotEnvironment, err := client.Environments.Get(ctx, environment.ID)
	saCheck(t, err)
	if gotEnvironment.ID != environment.ID {
		t.Fatal("Environment read did not match created ID")
	}
	identity, err := client.Identities.New(ctx, forward.IdentityNewParams{ExternalID: "sdk-go-sa-identity-" + testsupport.Marker(t)})
	saCheck(t, err)
	saCleanup(t, "identity "+identity.ID, func(ctx context.Context) error { _, err := client.Identities.Delete(ctx, identity.ID); return err })
	template, err := client.Templates.New(ctx, forward.TemplateNewParams{Name: "sdk-go-sa-template-" + testsupport.Marker(t), EnvironmentID: environment.ID, Model: forward.ModelConfigUnionParam{OfString: forward.String(model)}, System: forward.String("Echo the supplied text directly without using tools.")})
	saCheck(t, err)
	saCleanup(t, "template "+template.ID, func(ctx context.Context) error {
		_, err := client.Templates.Archive(ctx, template.ID, forward.TemplateArchiveParams{})
		return err
	})
	gotTemplate, err := client.Templates.Get(ctx, template.ID)
	saCheck(t, err)
	if gotTemplate.ID != template.ID {
		t.Fatal("Template read did not match created ID")
	}
	session, err := client.Sessions.New(ctx, forward.SessionNewParams{IdentityID: identity.ID, TemplateID: template.ID})
	saCheck(t, err)
	t.Logf("created session=%s model=%s", session.ID, model)
	saCleanup(t, "session "+session.ID, func(ctx context.Context) error {
		return saFinishSession(ctx, func() (string, error) {
			s, err := client.Sessions.Get(ctx, session.ID)
			if err != nil {
				return "", err
			}
			return string(s.Status), nil
		}, func() error {
			_, err := client.Sessions.Cancel(ctx, session.ID, forward.SessionCancelParams{})
			return err
		}, func() error {
			_, err := client.Sessions.Archive(ctx, session.ID, forward.SessionArchiveParams{})
			return err
		})
	})
	gotSession, err := client.Sessions.Get(ctx, session.ID)
	saCheck(t, err)
	if gotSession.ID != session.ID {
		t.Fatal("Session read did not match created ID")
	}
	expected := testsupport.Marker(t)
	sent, err := client.Sessions.Events.Send(ctx, session.ID, forward.SessionEventSendParams{Events: []forward.SessionEventParam{{Type: "user.message", Content: forward.EventContentUnionParam{OfBlocks: []forward.ContentBlockParam{{Type: "text", Text: forward.String("Reply with exactly " + expected)}}}}}, IdempotencyKey: forward.String("sdk-go-sa-event-" + testsupport.Marker(t))})
	saCheck(t, err)
	if sent == nil || len(sent.Data) != 1 || sent.Data[0].ID == "" {
		t.Fatal("Send must return one user event ID")
	}
	reply := testsupport.TurnResult{LastID: sent.Data[0].ID}
	stream := client.Sessions.Events.StreamEvents(ctx, session.ID, forward.SessionEventStreamParams{LastEventID: forward.String(sent.Data[0].ID)}, option.WithRequestTimeout(180*time.Second))
	defer stream.Close()
	for stream.Next() {
		saCheck(t, reply.Observe(stream.Current().RawJSON()))
		if reply.Complete {
			break
		}
	}
	saCheck(t, stream.Err())
	if err := reply.Verify([]string{expected}, false); err != nil {
		t.Fatalf("session=%s: %v", session.ID, err)
	}
	t.Logf("verified session=%s last_event_id=%s via SSE", session.ID, reply.LastID)
}

func saManagedExecution(t *testing.T, ctx context.Context, client managed.Client, model string) {
	environment, err := client.Environments.New(ctx, managed.EnvironmentNewParams{Name: "sdk-go-sa-env-" + testsupport.Marker(t), Config: managed.EnvironmentNewParamsConfigUnion{OfCloud: &managed.CloudConfigParams{Type: "cloud"}}})
	saCheck(t, err)
	saCleanup(t, "environment "+environment.ID, func(ctx context.Context) error {
		_, err := client.Environments.Archive(ctx, environment.ID, managed.EnvironmentArchiveParams{})
		return err
	})
	gotEnvironment, err := client.Environments.Get(ctx, environment.ID, managed.EnvironmentGetParams{})
	saCheck(t, err)
	if gotEnvironment.ID != environment.ID {
		t.Fatal("Environment read did not match created ID")
	}
	agent, err := client.Agents.New(ctx, managed.AgentNewParams{Name: "sdk-go-sa-agent-" + testsupport.Marker(t), Model: managed.ManagedAgentsModelConfigParams{ID: managed.ManagedAgentsModel(model)}, System: managed.String("Echo the supplied text directly without using tools.")})
	saCheck(t, err)
	saCleanup(t, "agent "+agent.ID, func(ctx context.Context) error {
		_, err := client.Agents.Archive(ctx, agent.ID, managed.AgentArchiveParams{})
		return err
	})
	gotAgent, err := client.Agents.Get(ctx, agent.ID, managed.AgentGetParams{})
	saCheck(t, err)
	if gotAgent.ID != agent.ID {
		t.Fatal("Agent read did not match created ID")
	}
	session, err := client.Sessions.New(ctx, managed.SessionNewParams{Agent: managed.SessionNewParamsAgentUnion{OfString: managed.String(agent.ID)}, EnvironmentID: environment.ID})
	saCheck(t, err)
	t.Logf("created session=%s model=%s", session.ID, model)
	saCleanup(t, "session "+session.ID, func(ctx context.Context) error {
		return saFinishSession(ctx, func() (string, error) {
			s, err := client.Sessions.Get(ctx, session.ID, managed.SessionGetParams{})
			if err != nil {
				return "", err
			}
			return string(s.Status), nil
		}, func() error {
			_, err := client.Sessions.Cancel(ctx, session.ID, managed.SessionCancelParams{})
			return err
		}, func() error {
			_, err := client.Sessions.Delete(ctx, session.ID, managed.SessionDeleteParams{})
			return err
		})
	})
	gotSession, err := client.Sessions.Get(ctx, session.ID, managed.SessionGetParams{})
	saCheck(t, err)
	if gotSession.ID != session.ID {
		t.Fatal("Session read did not match created ID")
	}
	expected := testsupport.Marker(t)
	sent, err := client.Sessions.Events.Send(ctx, session.ID, managed.SessionEventSendParams{Events: []managed.ManagedAgentsEventParamsUnion{{OfUserMessage: &managed.ManagedAgentsUserMessageEventParams{Type: "user.message", Content: []managed.ManagedAgentsUserMessageEventParamsContentUnion{{OfText: &managed.ManagedAgentsTextBlockParam{Type: "text", Text: "Reply with exactly " + expected}}}}}}}, option.WithHeader("Idempotency-Key", "sdk-go-sa-event-"+testsupport.Marker(t)))
	saCheck(t, err)
	if sent == nil || len(sent.Data) != 1 || sent.Data[0].ID == "" {
		t.Fatal("Send must return one user event ID")
	}
	reply := testsupport.TurnResult{LastID: sent.Data[0].ID}
	stream := client.Sessions.Events.StreamEvents(ctx, session.ID, managed.SessionEventStreamParams{}, option.WithHeader("Last-Event-ID", sent.Data[0].ID), option.WithRequestTimeout(180*time.Second))
	defer stream.Close()
	for stream.Next() {
		saCheck(t, reply.Observe(stream.Current().RawJSON()))
		if reply.Complete {
			break
		}
	}
	saCheck(t, stream.Err())
	if err := reply.Verify([]string{expected}, false); err != nil {
		t.Fatalf("session=%s: %v", session.ID, err)
	}
	t.Logf("verified session=%s last_event_id=%s via SSE", session.ID, reply.LastID)
}
