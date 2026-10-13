// Realtime is Beta and requires a pre-enabled Template. This demo uses process
// environment variables; it does not create/delete Identities or Templates.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Realtime example failed: "+live.Safe(err))
		os.Exit(1)
	}
}

func run() error {
	conversationID := flag.String("conversation", os.Getenv("QODER_CONVERSATION_ID"), "existing Conversation; omit to create one")
	identityID := flag.String("identity", os.Getenv("QODER_IDENTITY_ID"), "Identity for a new Conversation")
	templateID := flag.String("template", os.Getenv("QODER_TEMPLATE_ID"), "pre-enabled Realtime Template")
	idempotencyKey := flag.String("idempotency-key", os.Getenv("QODER_IDEMPOTENCY_KEY"), "stable creation operation key; required when creating")
	text := flag.String("text", "", "optional text input (uses model quota)")
	pcmFile := flag.String("pcm", "", "optional raw 16 kHz mono PCM16 little-endian input; no WAV header")
	historyOnly := flag.Bool("history-only", false, "only paginate an existing Conversation's persisted history")
	timeout := flag.Duration("timeout", 2*time.Minute, "overall example timeout")
	flag.Parse()
	if *text != "" && *pcmFile != "" {
		return fmt.Errorf("choose -text or -pcm for this single-turn example")
	}
	var inputPCM []byte
	if *pcmFile != "" {
		var err error
		inputPCM, err = os.ReadFile(*pcmFile)
		if err != nil {
			return err
		}
		if len(inputPCM) == 0 || len(inputPCM)%2 != 0 || len(inputPCM) > 30*16000*2 {
			return fmt.Errorf("PCM input must contain complete nonempty 16-bit samples and at most 30 seconds")
		}
		if len(inputPCM) >= 12 && string(inputPCM[:4]) == "RIFF" && string(inputPCM[8:12]) == "WAVE" {
			return fmt.Errorf("PCM input must be raw PCM16, not a WAV file")
		}
	}
	region := os.Getenv("QODER_SA_REGION")
	if region == "" {
		region = "cn"
	}
	if region != "cn" && region != "global" {
		return fmt.Errorf("QODER_SA_REGION must be cn or global")
	}
	host := "https://api.qoder.com"
	if region == "cn" {
		host = "https://api.qoder.com.cn"
	}
	base := os.Getenv("QODER_FORWARD_BASE_URL")
	if base == "" {
		base = host + "/api/v1/forward"
	}
	opts := []option.RequestOption{option.WithBaseURL(base)}
	mode := os.Getenv("QODER_AUTH_MODE")
	if mode == "" {
		mode = "pat"
	}
	if mode != "pat" && mode != "sa" {
		return fmt.Errorf("QODER_AUTH_MODE must be pat or sa")
	}
	opts = append(opts, option.WithAuthMode(mode))
	if mode == "sa" && os.Getenv("QODER_SA_KEY") != "" {
		if os.Getenv("QODER_ACCESS_TOKEN") != "" {
			return fmt.Errorf("choose SA Key exchange or direct SAT, not both")
		}
		credential, err := convention.NewServiceAccountCredential(convention.ServiceAccountCredentialConfig{SAKey: os.Getenv("QODER_SA_KEY"), Region: region, Scopes: []string{"forward.access"}})
		if err != nil {
			return err
		}
		opts = append(opts, option.WithCredential(credential))
	}
	client := forward.NewClient(opts...)
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	if *historyOnly && *conversationID == "" {
		return fmt.Errorf("history-only requires -conversation")
	}
	if *conversationID == "" {
		if *identityID == "" || *templateID == "" || *idempotencyKey == "" {
			return fmt.Errorf("creation requires -identity, -template and a stable -idempotency-key")
		}
		conversation, err := client.Realtime.Conversations.New(ctx, forward.RealtimeConversationNewParams{IdentityID: *identityID, TemplateID: *templateID, IdempotencyKey: forward.String(*idempotencyKey), Title: forward.String("SDK Realtime example")})
		if err != nil {
			return err
		}
		*conversationID = conversation.ID
		fmt.Println("Conversation:", *conversationID, "(retained by the service; closing does not delete it)")
	}
	if !*historyOnly {
		conn, err := client.Realtime.Connect(ctx, *conversationID, forward.RealtimeConnectParams{})
		if err != nil {
			return err
		}
		defer conn.Close()
		fmt.Printf("Ready: input=%d Hz output=%d Hz request_id=%s\n", conn.Ready.InputAudio.SampleRate, conn.Ready.OutputAudio.SampleRate, conn.RequestID)
		if *text != "" {
			if err = conn.SendText(ctx, *text); err != nil {
				return err
			}
		}
		audioErr := make(chan error, 1)
		sentAudio := len(inputPCM) == 0
		if len(inputPCM) > 0 {
			go func() {
				// 640 bytes = 20 ms at 16 kHz, mono PCM16. Pace input in real time.
				for start := 0; start < len(inputPCM); start += 640 {
					end := min(start+640, len(inputPCM))
					if err := conn.AppendAudio(ctx, inputPCM[start:end]); err != nil {
						audioErr <- err
						return
					}
					timer := time.NewTimer(20 * time.Millisecond)
					select {
					case <-ctx.Done():
						timer.Stop()
						audioErr <- ctx.Err()
						return
					case <-timer.C:
					}
				}
				audioErr <- nil
			}()
		}
		if err = conn.Ping(ctx); err != nil {
			return err
		}
		seenUser, seenAssistant, seenAudioDone := len(inputPCM) == 0, false, false
		receipts := map[string]bool{}
		for {
			select {
			case sendErr := <-audioErr:
				if sendErr != nil {
					return sendErr
				}
				sentAudio = true
			default:
			}
			event, recvErr := conn.Recv(ctx)
			if recvErr != nil {
				return recvErr
			}
			fmt.Printf("Event: %s work_id=%s announcement_id=%s\n", event.Type, event.WorkID, event.AnnouncementID)
			if event.Type == "transcript.final" {
				var p forward.RealtimeTranscriptPayload
				if err = event.DecodePayload(&p); err != nil {
					return err
				}
				if strings.TrimSpace(p.Text) != "" {
					seenUser = seenUser || p.Role == "user"
					seenAssistant = seenAssistant || p.Role == "assistant"
				}
			}
			if event.Type == "audio.delta" {
				pcm, decodeErr := event.Audio()
				if decodeErr != nil {
					return decodeErr
				}
				fmt.Println("Output PCM bytes:", len(pcm))
				// This demo has no speaker; never claim playback started/ended.
				key := event.WorkID + "/" + event.AnnouncementID
				if !receipts[key] {
					receipts[key] = true
					if err = conn.PlaybackCancelled(ctx, forward.RealtimePlaybackReceipt{WorkID: event.WorkID, AnnouncementID: event.AnnouncementID}); err != nil {
						return err
					}
				}
			}
			if event.Type == "error" {
				var p forward.RealtimeErrorPayload
				_ = json.Unmarshal(event.Payload, &p)
				return fmt.Errorf("Realtime returned error code %s", p.Code)
			}
			if event.Type == "audio.done" {
				seenAudioDone = true
			}
			if (event.Type == "pong" && *text == "" && len(inputPCM) == 0) || (seenUser && seenAssistant && seenAudioDone) {
				break
			}
		}
		if !sentAudio {
			select {
			case err = <-audioErr:
				if err != nil {
					return err
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		closeCtx, cancelClose := context.WithTimeout(context.Background(), 15*time.Second)
		result, closeErr := conn.CloseGracefully(closeCtx)
		cancelClose()
		if closeErr != nil {
			return closeErr
		}
		fmt.Println("Persisted close outcome:", result.Outcome)
	}
	pager := client.Realtime.Conversations.HistoryAutoPaging(ctx, *conversationID, forward.RealtimeConversationHistoryParams{Types: forward.String("message,work"), Limit: forward.Int(50)})
	for pager.Next() {
		event := pager.Current()
		fmt.Printf("History: %s %s\n", event.ID, event.Type)
	}
	return pager.Err()
}
