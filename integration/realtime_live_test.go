//go:build live

package integration_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
)

func realtimeEnvironment(path string, environment map[string]string) (map[string]string, error) {
	values, err := liveEnvironment(path, environment)
	if err != nil {
		return nil, err
	}
	for _, name := range []string{"QODER_REALTIME_LIVE_ALLOW_WRITE", "QODER_REALTIME_LIVE_ALLOW_EXECUTION"} {
		if values[name] != "true" {
			return nil, fmt.Errorf("%s=true is required", name)
		}
	}
	for _, region := range []string{"CN", "GLOBAL"} {
		sat, key := "QODER_REALTIME_"+region+"_SAT", "QODER_SA_"+region+"_KEY"
		if strings.TrimSpace(values[sat]) == "" && strings.TrimSpace(values[key]) == "" {
			return nil, fmt.Errorf("%s or %s is required", sat, key)
		}
		pat := "QODER_FORWARD_" + region + "_PAT"
		if strings.TrimSpace(values[pat]) == "" {
			return nil, fmt.Errorf("%s is required", pat)
		}
		for _, mode := range []string{"PAT", "SA"} {
			name := "QODER_REALTIME_" + region + "_" + mode + "_TEMPLATE_ID"
			if strings.TrimSpace(values[name]) == "" {
				return nil, fmt.Errorf("%s is required", name)
			}
		}
	}
	return values, nil
}

func TestRealtimeConfiguration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no-env-file")
	valid := map[string]string{"QODER_SA_CN_KEY": "fake-cn-key", "QODER_SA_GLOBAL_KEY": "fake-global-key", "QODER_FORWARD_CN_PAT": "fake-cn-pat", "QODER_FORWARD_GLOBAL_PAT": "fake-global-pat", "QODER_REALTIME_LIVE_ALLOW_WRITE": "true", "QODER_REALTIME_LIVE_ALLOW_EXECUTION": "true"}
	for _, region := range []string{"CN", "GLOBAL"} {
		for _, mode := range []string{"PAT", "SA"} {
			valid["QODER_REALTIME_"+region+"_"+mode+"_TEMPLATE_ID"] = "tmpl_fake"
		}
	}
	if _, err := realtimeEnvironment(path, valid); err != nil {
		t.Fatal(err)
	}
	for key := range valid {
		copy := map[string]string{}
		for k, v := range valid {
			copy[k] = v
		}
		copy[key] = ""
		if _, err := realtimeEnvironment(path, copy); err == nil {
			t.Fatalf("missing %s did not fail before requests", key)
		}
	}
	withSAT := map[string]string{}
	for key, value := range valid {
		withSAT[key] = value
	}
	for _, region := range []string{"CN", "GLOBAL"} {
		withSAT["QODER_SA_"+region+"_KEY"] = ""
		withSAT["QODER_REALTIME_"+region+"_SAT"] = "fake-bound-sat"
	}
	if _, err := realtimeEnvironment(path, withSAT); err != nil {
		t.Fatal("a supplied bound SAT must not require a platform SA Key")
	}
	pcmPath := filepath.Join(t.TempDir(), "input.pcm")
	for _, content := range [][]byte{nil, {0}, {0, 0, 1, 0}} {
		if err := os.WriteFile(pcmPath, content, 0600); err != nil {
			t.Fatal(err)
		}
		pcm, err := realtimePCMInput(pcmPath)
		if len(content) == 4 {
			if err != nil || len(pcm) != 4 {
				t.Fatal("valid raw PCM fixture rejected")
			}
		} else if err == nil {
			t.Fatal("invalid raw PCM fixture accepted")
		}
	}
}

func realtimePCMInput(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}
	pcm, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read Realtime PCM fixture")
	}
	if len(pcm) == 0 || len(pcm)%2 != 0 || len(pcm) > 30*16000*2 {
		return nil, fmt.Errorf("Realtime PCM fixture must be nonempty complete PCM16 samples and at most 30 seconds")
	}
	if len(pcm) >= 12 && string(pcm[:4]) == "RIFF" && string(pcm[8:12]) == "WAVE" {
		return nil, fmt.Errorf("Realtime PCM fixture must be raw PCM16, not a WAV file")
	}
	return pcm, nil
}

// TestRealtimeLive is explicit and never runs in the offline suite. Both regions
// and both auth modes are preflighted before any test Identity is created.
func TestRealtimeLive(t *testing.T) {
	path := os.Getenv("LIVE_ENV_FILE")
	if path == "" {
		path = filepath.Join("..", ".env.live")
	}
	environment := map[string]string{}
	for _, entry := range os.Environ() {
		name, value, _ := strings.Cut(entry, "=")
		environment[name] = value
	}
	values, err := realtimeEnvironment(path, environment)
	if err != nil {
		t.Fatal(err)
	}
	pcmPath := values["QODER_REALTIME_LIVE_PCM_FILE"]
	if pcmPath == "" {
		pcmPath = values["QODER_REALTIME_PCM_FILE"]
	}
	pcm, err := realtimePCMInput(pcmPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, region := range []string{"cn", "global"} {
		for _, mode := range []string{"pat", "sa"} {
			t.Run(region+"/"+mode, func(t *testing.T) {
				host := "https://api.qoder.com"
				if region == "cn" {
					host = "https://api.qoder.com.cn"
				}
				opts := []option.RequestOption{option.WithBaseURL(host + "/api/v1/forward"), option.WithMaxRetries(0), option.WithRequestTimeout(30 * time.Second), testsupport.RequestLog(t)}
				if mode == "pat" {
					opts = append(opts, option.WithPAT(values["QODER_FORWARD_"+strings.ToUpper(region)+"_PAT"]))
				} else if sat := values["QODER_REALTIME_"+strings.ToUpper(region)+"_SAT"]; strings.TrimSpace(sat) != "" {
					opts = append(opts, option.WithAuthMode("sa"), option.WithAuthToken(sat))
					t.Log("auth_mode=sa source=forward_bound_sat")
				} else {
					credential, err := convention.NewServiceAccountCredential(convention.ServiceAccountCredentialConfig{SAKey: values["QODER_SA_"+strings.ToUpper(region)+"_KEY"], Region: region, Scopes: []string{"forward.access"}})
					saCheck(t, err)
					opts = append(opts, option.WithAuthMode("sa"), option.WithCredential(credential))
					t.Log("auth_mode=sa source=platform_credential")
				}
				client := forward.NewClient(opts...)
				ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
				defer cancel()
				templateID := values["QODER_REALTIME_"+strings.ToUpper(region)+"_"+strings.ToUpper(mode)+"_TEMPLATE_ID"]
				identity, err := client.Identities.New(ctx, forward.IdentityNewParams{ExternalID: "sdk-go-realtime-" + testsupport.Marker(t), IdempotencyKey: forward.String("sdk-go-realtime-identity-" + testsupport.Marker(t))})
				saCheck(t, err)
				t.Logf("identity_id=%s template_id=%s", identity.ID, templateID)
				saCleanup(t, "realtime identity "+identity.ID, func(cleanCtx context.Context) error {
					// The Conversation's primary Forward Session belongs to this unique
					// test Identity. Finish/archive it before deleting the Identity.
					pager := client.Sessions.ListAutoPaging(cleanCtx, forward.SessionListParams{IdentityIDs: []string{identity.ID}, Limit: forward.Int(100)})
					var cleanupErrors []error
					for pager.Next() {
						session := pager.Current()
						if err := saFinishSession(cleanCtx, func() (string, error) {
							s, err := client.Sessions.Get(cleanCtx, session.ID)
							if err != nil {
								return "", err
							}
							return string(s.Status), nil
						}, func() error {
							_, err := client.Sessions.Cancel(cleanCtx, session.ID, forward.SessionCancelParams{})
							return err
						}, func() error {
							_, err := client.Sessions.Archive(cleanCtx, session.ID, forward.SessionArchiveParams{})
							return err
						}); err != nil {
							cleanupErrors = append(cleanupErrors, err)
						}
					}
					if pager.Err() != nil {
						cleanupErrors = append(cleanupErrors, pager.Err())
					}
					_, err := client.Identities.Delete(cleanCtx, identity.ID)
					if err != nil {
						cleanupErrors = append(cleanupErrors, err)
					}
					return errors.Join(cleanupErrors...)
				})
				conversation, err := client.Realtime.Conversations.New(ctx, forward.RealtimeConversationNewParams{IdentityID: identity.ID, TemplateID: templateID, IdempotencyKey: forward.String("sdk-go-realtime-conversation-" + testsupport.Marker(t)), Title: forward.String("SDK Realtime acceptance"), Metadata: map[string]any{"source": "sdk-go-live"}})
				saCheck(t, err)
				if conversation.ID == "" || conversation.Type != "voice.conversation" || conversation.Status != "ready" {
					t.Fatal("invalid Conversation response")
				}
				t.Logf("conversation_id=%s retained=true", conversation.ID)
				conn, err := client.Realtime.Connect(ctx, conversation.ID, forward.RealtimeConnectParams{})
				saCheck(t, err)
				defer conn.Close()
				t.Logf("websocket conversation_id=%s request_id=%s", conn.ConversationID, conn.RequestID)
				saCheck(t, conn.Ping(ctx))
				seenReady, seenPong := false, false
				for !seenPong {
					event, err := conn.Recv(ctx)
					saCheck(t, err)
					if event.Type == "voice.ready" {
						seenReady = true
					}
					if event.Type == "pong" {
						seenPong = true
					}
					realtimeLiveError(t, event)
				}
				if !seenReady {
					t.Fatal("voice.ready was consumed during Connect")
				}
				saCheck(t, conn.SendText(ctx, "请用一句话简单打招呼。"))
				audioBytes := realtimeLiveTurn(t, ctx, conn, false)
				if len(pcm) > 0 {
					sent := make(chan error, 1)
					go func() {
						for start := 0; start < len(pcm); start += 640 {
							if err := conn.AppendAudio(ctx, pcm[start:min(start+640, len(pcm))]); err != nil {
								sent <- err
								return
							}
							timer := time.NewTimer(20 * time.Millisecond)
							select {
							case <-ctx.Done():
								timer.Stop()
								sent <- ctx.Err()
								return
							case <-timer.C:
							}
						}
						sent <- nil
					}()
					audioBytes += realtimeLiveTurn(t, ctx, conn, true)
					select {
					case err := <-sent:
						saCheck(t, err)
					case <-ctx.Done():
						saCheck(t, ctx.Err())
					}
					t.Logf("conversation_id=%s input_pcm_bytes=%d user_transcript_final=true", conversation.ID, len(pcm))
				}
				result, err := conn.CloseGracefully(ctx)
				saCheck(t, err)
				// response.done can already have committed the message before close,
				// legitimately yielding no_active_response/no_content. The SDK validates
				// the acknowledgement; persisted user+assistant history remains mandatory.
				t.Logf("conversation_id=%s transcript_final=true audio_done=true audio_bytes=%d close_outcome=%s", conversation.ID, audioBytes, result.Outcome)
				for _, types := range []string{"message", "work", "message,work"} {
					pager := client.Realtime.Conversations.HistoryAutoPaging(ctx, conversation.ID, forward.RealtimeConversationHistoryParams{Limit: forward.Int(1), Types: forward.String(types)})
					ids := map[string]bool{}
					count := 0
					seenUser, seenAssistant := false, false
					for pager.Next() {
						event := pager.Current()
						if event.ID == "" || ids[event.ID] || event.OccurredAt == "" {
							t.Fatal("invalid/duplicate persisted history ID or timestamp")
						}
						ids[event.ID] = true
						count++
						seenUser = seenUser || strings.HasPrefix(event.Type, "voice.user_message.")
						seenAssistant = seenAssistant || strings.HasPrefix(event.Type, "voice.assistant_message.")
						if types == "message" && !strings.HasPrefix(event.Type, "voice.user_message.") && !strings.HasPrefix(event.Type, "voice.assistant_message.") {
							t.Fatal("message filter returned a Work record")
						}
						if types == "work" && !strings.HasPrefix(event.Type, "voice.work.") {
							t.Fatal("work filter returned a message record")
						}
					}
					saCheck(t, pager.Err())
					if types != "work" && (count < 2 || !seenUser || !seenAssistant) {
						t.Fatal("history did not contain persisted user/assistant messages with limit=1 pagination")
					}
					t.Logf("conversation_id=%s history_filter=%s count=%d", conversation.ID, types, count)
				}
			})
		}
	}
}

func realtimeLiveTurn(t *testing.T, ctx context.Context, conn *forward.RealtimeConnection, requireUserTranscript bool) int {
	t.Helper()
	seenUser := !requireUserTranscript
	seenFinal, seenAudioDone, audioBytes := false, false, 0
	receipts := map[string]bool{}
	for !seenUser || !seenFinal || !seenAudioDone {
		event, err := conn.Recv(ctx)
		saCheck(t, err)
		realtimeLiveError(t, event)
		switch event.Type {
		case "transcript.final":
			var p forward.RealtimeTranscriptPayload
			saCheck(t, event.DecodePayload(&p))
			if strings.TrimSpace(p.Text) != "" {
				if p.Role == "assistant" {
					seenFinal = true
				}
				if p.Role == "user" {
					seenUser = true
				}
			}
		case "audio.delta":
			pcm, err := event.Audio()
			saCheck(t, err)
			audioBytes += len(pcm)
			key := event.WorkID + "/" + event.AnnouncementID
			if !receipts[key] {
				receipts[key] = true
				saCheck(t, conn.PlaybackCancelled(ctx, forward.RealtimePlaybackReceipt{WorkID: event.WorkID, AnnouncementID: event.AnnouncementID}))
			}
		case "audio.done":
			seenAudioDone = true
		}
	}
	if audioBytes == 0 {
		t.Fatal("response produced no decodable PCM audio")
	}
	return audioBytes
}

func realtimeLiveError(t *testing.T, event forward.RealtimeEvent) {
	t.Helper()
	if event.Type == "error" {
		var p forward.RealtimeErrorPayload
		if event.DecodePayload(&p) != nil {
			t.Fatal("malformed server error event")
		}
		t.Fatalf("Realtime server error code=%s conversation_id=%s event_id=%s", p.Code, event.ConversationID, event.EventID)
	}
	if event.Type == "voice.replaced" {
		t.Fatalf("Realtime connection replaced conversation_id=%s", event.ConversationID)
	}
}
