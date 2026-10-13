package forward_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/coder/websocket"
)

func realtimeFixtures(t *testing.T) []map[string]any {
	t.Helper()
	var fixture struct {
		Events []map[string]any `json:"events"`
	}
	readJSON(t, "testdata/realtime-events.json", &fixture)
	return fixture.Events
}
func writeRealtime(t *testing.T, ws *websocket.Conn, event map[string]any) {
	t.Helper()
	b, err := json.Marshal(event)
	if err != nil {
		t.Error(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := ws.Write(ctx, websocket.MessageText, b); err != nil {
		t.Error(err)
	}
}
func readRealtime(t *testing.T, ws *websocket.Conn) map[string]any {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	kind, b, err := ws.Read(ctx)
	if err != nil {
		t.Error(err)
		return nil
	}
	if kind != websocket.MessageText {
		t.Error("client sent binary frame")
	}
	var event map[string]any
	if err = json.Unmarshal(b, &event); err != nil {
		t.Error(err)
	}
	return event
}
func realtimeServer(t *testing.T, handle func(*websocket.Conn, *http.Request)) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer ws.CloseNow()
		handle(ws, r)
	}))
	t.Cleanup(server.Close)
	return server
}
func realtimeClient(url string, opts ...option.RequestOption) forward.Client {
	return forward.NewClient(append([]option.RequestOption{option.WithPAT("test-pat"), option.WithBaseURL(url + "/prefix/api/v1/forward/")}, opts...)...)
}
func realtimeConnect(t *testing.T, client forward.Client, params forward.RealtimeConnectParams) *forward.RealtimeConnection {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	conn, err := client.Realtime.Connect(ctx, "conv_test", params)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

func TestRealtimeProtocolAndConvenienceMethods(t *testing.T) {
	fixture := realtimeFixtures(t)
	sent := make(chan []map[string]any, 1)
	server := realtimeServer(t, func(ws *websocket.Conn, r *http.Request) {
		if r.URL.Path != "/prefix/api/v1/forward/realtime" || r.URL.Query().Get("conversation_id") != "conv_test" || r.Header.Get("Authorization") != "Bearer test-pat" || r.Header.Get("X-Request-Header") != "request" {
			t.Error("handshake did not preserve URL/auth/options")
		}
		for _, event := range fixture[:2] {
			writeRealtime(t, ws, event)
		}
		var got []map[string]any
		for i := 0; i < 9; i++ {
			got = append(got, readRealtime(t, ws))
		}
		sent <- got
		for _, event := range fixture[2:16] {
			writeRealtime(t, ws, event)
		}
		closeEvent := readRealtime(t, ws)
		requestID := closeEvent["payload"].(map[string]any)["request_id"]
		wrong := fixture[16]
		wrong["payload"] = map[string]any{"request_id": "wrong", "outcome": "saved"}
		writeRealtime(t, ws, wrong)
		good := fixture[16]
		good["payload"] = map[string]any{"request_id": requestID, "outcome": "saved_interrupted"}
		writeRealtime(t, ws, good)
		_, _, _ = ws.Read(context.Background())
	})
	client := realtimeClient(server.URL, option.WithHeader("X-Request-Header", "client"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := client.Realtime.Connect(ctx, "conv_test", forward.RealtimeConnectParams{}, option.WithHeader("X-Request-Header", "request"))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if conn.Ready.InputAudio.SampleRate != 16000 || conn.Ready.OutputAudio.SampleRate != 24000 {
		t.Fatal("ready configuration missing")
	}
	for i, kind := range []string{"voice.state", "voice.ready"} {
		event, err := conn.Recv(ctx)
		if err != nil || event.Type != kind || event.Sequence != int64(i+1) {
			t.Fatalf("startup events lost: %#v %v", event, err)
		}
	}
	receipt := forward.RealtimePlaybackReceipt{WorkID: "work_one", AnnouncementID: "ann_one"}
	for _, send := range []func() error{
		func() error { return conn.SendText(ctx, "你好") }, func() error { return conn.AppendAudio(ctx, []byte{0, 0, 1, 0}) }, func() error { return conn.Interrupt(ctx) }, func() error { return conn.PlaybackStarted(ctx, receipt) }, func() error { return conn.PlaybackEnded(ctx, receipt) }, func() error { return conn.PlaybackCancelled(ctx, receipt) }, func() error { return conn.Ping(ctx) }, func() error {
			return conn.SendEvent(ctx, forward.RealtimeClientEvent{Type: "future.client", Payload: map[string]any{"extra": true}})
		}, func() error { return conn.PlaybackEnded(ctx, forward.RealtimePlaybackReceipt{}) },
	} {
		if err := send(); err != nil {
			t.Fatal(err)
		}
	}
	got := <-sent
	for i, kind := range []string{"text.message", "audio.append", "interrupt", "playback.started", "playback.ended", "playback.cancelled", "ping", "future.client", "playback.ended"} {
		if got[i]["version"] != "voice.realtime.v1" || got[i]["type"] != kind || got[i]["payload"] == nil {
			t.Errorf("bad client event %#v", got[i])
		}
	}
	if got[0]["payload"].(map[string]any)["text"] != "你好" || got[1]["payload"].(map[string]any)["audio"] != "AAABAA==" || !reflect.DeepEqual(got[3]["payload"], map[string]any{"work_id": "work_one", "announcement_id": "ann_one"}) || len(got[8]["payload"].(map[string]any)) != 0 {
		t.Fatalf("client wire payloads %#v", got)
	}
	for _, expected := range fixture[2:16] {
		event, err := conn.Recv(ctx)
		if err != nil || event.Type != expected["type"] || !strings.Contains(event.RawJSON(), `"future_envelope":"preserved"`) || event.Timestamp != "2026-10-13T01:02:03.123456789Z" {
			t.Fatalf("bad server event %#v %v", event, err)
		}
		if event.Type == "audio.delta" {
			pcm, err := event.Audio()
			if err != nil || !reflect.DeepEqual(pcm, []byte{0, 0, 1, 0}) {
				t.Fatal("audio decode failed")
			}
		}
		if event.Type == "work.failed" {
			var p forward.RealtimeWorkFailedPayload
			if event.DecodePayload(&p) != nil || p.Error.Code != "task_error" {
				t.Fatal("typed payload decoding failed")
			}
		}
	}
	result, err := conn.CloseGracefully(ctx)
	if err != nil || result.Outcome != "saved_interrupted" {
		t.Fatalf("close failed %#v %v", result, err)
	}
	again, err := conn.CloseGracefully(ctx)
	if err != nil || *again != *result {
		t.Fatal("repeated graceful close did not retain confirmation")
	}
	if conn.SendText(ctx, "after") == nil {
		t.Fatal("send after close succeeded")
	}
}

func TestRealtimeInvalidInputBounds(t *testing.T) {
	fixture := realtimeFixtures(t)
	server := realtimeServer(t, func(ws *websocket.Conn, _ *http.Request) {
		writeRealtime(t, ws, fixture[1])
		_, _, _ = ws.Read(context.Background())
	})
	conn := realtimeConnect(t, realtimeClient(server.URL), forward.RealtimeConnectParams{})
	ctx := context.Background()
	for _, text := range []string{"", " \t", strings.Repeat("你", 5462), string([]byte{0xff})} {
		if conn.SendText(ctx, text) == nil {
			t.Fatal("invalid text accepted")
		}
	}
	for _, pcm := range [][]byte{nil, {1}, make([]byte, 256*1024)} {
		if conn.AppendAudio(ctx, pcm) == nil {
			t.Fatal("invalid audio accepted")
		}
	}
	if conn.SendEvent(ctx, forward.RealtimeClientEvent{Type: "future", Payload: map[string]any{"text": strings.Repeat("x", 256*1024)}}) == nil {
		t.Fatal("oversized JSON accepted")
	}
	if conn.SendEvent(ctx, forward.RealtimeClientEvent{Version: "bad", Type: "ping"}) == nil {
		t.Fatal("bad client version accepted")
	}
}

func TestRealtimeReadyCancellationAndHTTPTimeoutIndependence(t *testing.T) {
	fixture := realtimeFixtures(t)
	server := realtimeServer(t, func(ws *websocket.Conn, _ *http.Request) {
		time.Sleep(50 * time.Millisecond)
		writeRealtime(t, ws, fixture[1])
		_, _, _ = ws.Read(context.Background())
	})
	client := realtimeClient(server.URL, option.WithHTTPClient(&http.Client{Timeout: time.Millisecond}), option.WithRequestTimeout(time.Millisecond))
	conn := realtimeConnect(t, client, forward.RealtimeConnectParams{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _ = conn.Recv(ctx) // queued ready remains observable even with a cancelled wait.
	_, err := conn.Recv(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("recv cancellation %v", err)
	}
	blocked := realtimeServer(t, func(ws *websocket.Conn, _ *http.Request) { _, _, _ = ws.Read(context.Background()) })
	client = realtimeClient(blocked.URL)
	_, err = client.Realtime.Connect(context.Background(), "conv_test", forward.RealtimeConnectParams{ReadyTimeout: 10 * time.Millisecond})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("ready timeout %v", err)
	}
	ctx, cancel = context.WithCancel(context.Background())
	time.AfterFunc(10*time.Millisecond, cancel)
	_, err = client.Realtime.Connect(ctx, "conv_test", forward.RealtimeConnectParams{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("connect cancellation %v", err)
	}
}

func TestRealtimeProtocolFailures(t *testing.T) {
	fixture := realtimeFixtures(t)
	for _, kind := range []string{"invalid_json", "binary", "version", "missing_envelope", "null_sequence", "timestamp_comma", "timestamp_offset", "payload_type", "ready_audio", "ready_capability_missing", "ready_capability_type", "initialization_error"} {
		t.Run(kind, func(t *testing.T) {
			server := realtimeServer(t, func(ws *websocket.Conn, _ *http.Request) {
				event := realtimeFixtures(t)[1]
				switch kind {
				case "invalid_json":
					_ = ws.Write(context.Background(), websocket.MessageText, []byte("{broken"))
				case "binary":
					_ = ws.Write(context.Background(), websocket.MessageBinary, []byte{0, 0})
				case "version":
					event["version"] = "v2"
					writeRealtime(t, ws, event)
				case "missing_envelope":
					delete(event, "sequence")
					writeRealtime(t, ws, event)
				case "null_sequence":
					event["sequence"] = nil
					writeRealtime(t, ws, event)
				case "timestamp_comma":
					event["timestamp"] = "2026-10-13T01:02:03,123Z"
					writeRealtime(t, ws, event)
				case "timestamp_offset":
					event["timestamp"] = "2026-10-13T01:02:03+24:00"
					writeRealtime(t, ws, event)
				case "payload_type":
					event["payload"] = []any{}
					writeRealtime(t, ws, event)
				case "ready_audio":
					event["payload"].(map[string]any)["input_audio"] = map[string]any{"format": "wav", "sample_rate": 16000, "channels": 1}
					writeRealtime(t, ws, event)
				case "ready_capability_missing":
					event["payload"].(map[string]any)["capabilities"] = map[string]any{}
					writeRealtime(t, ws, event)
				case "ready_capability_type":
					event["payload"].(map[string]any)["capabilities"] = map[string]any{"graceful_close": "true"}
					writeRealtime(t, ws, event)
				case "initialization_error":
					event["type"] = "error"
					event["payload"] = map[string]any{"code": "provider_initialization_failed", "message": "init failed", "retry_after_ms": 5}
					writeRealtime(t, ws, event)
				}
				_, _, _ = ws.Read(context.Background())
			})
			client := realtimeClient(server.URL)
			_, err := client.Realtime.Connect(context.Background(), "conv_test", forward.RealtimeConnectParams{ReadyTimeout: time.Second})
			if err == nil {
				t.Fatal("bad server protocol accepted")
			}
			if kind == "initialization_error" {
				var serverErr *forward.RealtimeServerError
				if !errors.As(err, &serverErr) || serverErr.Code != "provider_initialization_failed" {
					t.Fatalf("lost server error %v", err)
				}
			} else {
				var protocol *forward.RealtimeProtocolError
				if !errors.As(err, &protocol) {
					t.Fatalf("bad error type %T %v", err, err)
				}
			}
		})
	}
	release := make(chan struct{})
	server := realtimeServer(t, func(ws *websocket.Conn, _ *http.Request) {
		writeRealtime(t, ws, fixture[1])
		<-release
		writeRealtime(t, ws, fixture[17])
		_, _, _ = ws.Read(context.Background())
	})
	client := realtimeClient(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := client.Realtime.Connect(ctx, "conv_test", forward.RealtimeConnectParams{})
	if err != nil {
		t.Fatal(err)
	}
	close(release)
	defer conn.Close()
	_, _ = conn.Recv(ctx)
	event, err := conn.Recv(ctx)
	if err != nil || event.Type != "voice.replaced" {
		t.Fatalf("replacement event missing %v", err)
	}
	_, err = conn.Recv(ctx)
	if err == nil {
		t.Fatal("replaced connection remained active")
	}
}

func TestRealtimeGracefulCloseFailure(t *testing.T) {
	for _, kind := range []string{"wrong_request_id", "persistence_failed", "bad_outcome"} {
		t.Run(kind, func(t *testing.T) {
			fixture := realtimeFixtures(t)
			server := realtimeServer(t, func(ws *websocket.Conn, _ *http.Request) {
				writeRealtime(t, ws, fixture[1])
				event := readRealtime(t, ws)
				id := event["payload"].(map[string]any)["request_id"]
				response := fixture[16]
				response["payload"] = map[string]any{"request_id": id, "outcome": "saved"}
				switch kind {
				case "wrong_request_id":
					response["payload"].(map[string]any)["request_id"] = "wrong"
				case "persistence_failed":
					response["type"] = "error"
					response["payload"] = map[string]any{"code": "event_persistence_failed", "message": "storage failed"}
				case "bad_outcome":
					response["payload"].(map[string]any)["outcome"] = "lost"
				}
				writeRealtime(t, ws, response)
				_, _, _ = ws.Read(context.Background())
			})
			conn := realtimeConnect(t, realtimeClient(server.URL), forward.RealtimeConnectParams{CloseTimeout: 30 * time.Millisecond})
			_, err := conn.CloseGracefully(context.Background())
			if err == nil {
				t.Fatal("unconfirmed persistence reported success")
			}
			if conn.SendText(context.Background(), "closed") == nil {
				t.Fatal("failed graceful close leaked active connection")
			}
		})
	}
}

func TestRealtimeBackpressureAndReceiveLimit(t *testing.T) {
	for _, kind := range []string{"events", "bytes", "message"} {
		t.Run(kind, func(t *testing.T) {
			fixture := realtimeFixtures(t)
			release := make(chan struct{})
			server := realtimeServer(t, func(ws *websocket.Conn, _ *http.Request) {
				writeRealtime(t, ws, fixture[1])
				<-release
				for i := 0; i < 4; i++ {
					event := fixture[15]
					size := 400
					if kind == "message" {
						size = 800
					}
					event["payload"] = map[string]any{"text": strings.Repeat("x", size)}
					b, _ := json.Marshal(event)
					if ws.Write(context.Background(), websocket.MessageText, b) != nil {
						return
					}
				}
				_, _, _ = ws.Read(context.Background())
			})
			params := forward.RealtimeConnectParams{}
			switch kind {
			case "events":
				params.QueueCapacity = 1
			case "bytes":
				params.QueueBytes = 900
			case "message":
				params.ReadLimit = 650
			}
			conn := realtimeConnect(t, realtimeClient(server.URL), params)
			_, _ = conn.Recv(context.Background())
			close(release)
			deadline := time.Now().Add(time.Second)
			for conn.Err() == nil && time.Now().Before(deadline) {
				time.Sleep(time.Millisecond)
			}
			if conn.Err() == nil {
				t.Fatal("slow consumer/message limit did not end connection")
			}
		})
	}
}

type realtimeCredential struct{ calls atomic.Int32 }

func (c *realtimeCredential) Token(context.Context) (string, error) {
	n := c.calls.Add(1)
	if n == 1 {
		return "sat-one", nil
	}
	return "sat-two", nil
}
func TestRealtimeAuthenticationAndRedirectSafety(t *testing.T) {
	fixture := realtimeFixtures(t)
	credential := &realtimeCredential{}
	var auths []string
	var mu sync.Mutex
	server := realtimeServer(t, func(ws *websocket.Conn, r *http.Request) {
		mu.Lock()
		auths = append(auths, r.Header.Get("Authorization"))
		mu.Unlock()
		if strings.Contains(r.URL.String(), "sat-") {
			t.Error("token leaked in URL")
		}
		writeRealtime(t, ws, fixture[1])
		_, _, _ = ws.Read(context.Background())
	})
	client := forward.NewClient(option.WithCredential(credential), option.WithBaseURL(server.URL+"/api/v1/forward"))
	for i := 0; i < 2; i++ {
		conn := realtimeConnect(t, client, forward.RealtimeConnectParams{})
		_ = conn.Close()
	}
	mu.Lock()
	if !reflect.DeepEqual(auths, []string{"Bearer sat-one", "Bearer sat-two"}) {
		t.Errorf("credential reused stale token: %v", auths)
	}
	mu.Unlock()
	conn, err := client.Realtime.Connect(context.Background(), "conv_test", forward.RealtimeConnectParams{}, option.WithHeader("Authorization", "Bearer explicit"))
	if err != nil {
		t.Fatal(err)
	}
	_ = conn.Close()
	if credential.calls.Load() != 2 {
		t.Fatal("explicit header did not override credential")
	}
	targetCalls := atomic.Int32{}
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { targetCalls.Add(1); w.WriteHeader(401) }))
	defer target.Close()
	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer redirect.Close()
	client = realtimeClient(redirect.URL)
	_, err = client.Realtime.Connect(context.Background(), "conv_test", forward.RealtimeConnectParams{})
	var api *forward.Error
	if !errors.As(err, &api) || api.StatusCode != 307 || targetCalls.Load() != 0 {
		t.Fatalf("redirect followed: %v", err)
	}
	for _, body := range []string{`{"code":"TOKEN_INVALID","message":"expired"}`, "proxy rejected"} {
		bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-Id", "handshake-failure")
			w.WriteHeader(401)
			_, _ = io.WriteString(w, body)
		}))
		client = realtimeClient(bad.URL)
		_, err = client.Realtime.Connect(context.Background(), "conv_test", forward.RealtimeConnectParams{})
		bad.Close()
		if !errors.As(err, &api) || api.StatusCode != 401 || api.RequestID != "handshake-failure" || api.Message == "" {
			t.Fatalf("handshake error lost %v", err)
		}
	}
}

type realtimeBlockedCredential struct{}

func (realtimeBlockedCredential) Token(ctx context.Context) (string, error) {
	<-ctx.Done()
	return "", ctx.Err()
}

func TestRealtimeHandshakeTimeoutIncludesCredential(t *testing.T) {
	client := forward.NewClient(option.WithCredential(realtimeBlockedCredential{}), option.WithBaseURL("https://example.invalid/api/v1/forward"))
	_, err := client.Realtime.Connect(context.Background(), "conv_test", forward.RealtimeConnectParams{HandshakeTimeout: 10 * time.Millisecond})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("credential wait escaped handshake timeout: %v", err)
	}
}

func TestRealtimeSharedServiceAccountCredentialTLSAndMiddleware(t *testing.T) {
	fixture := realtimeFixtures(t)
	exchanges := atomic.Int32{}
	credential, err := convention.NewServiceAccountCredential(convention.ServiceAccountCredentialConfig{SAKey: "fake-sa-key", Region: "cn", Scopes: []string{"forward.access"}, HTTPClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://openapi.qoder.com.cn/api/v1/serviceToken/exchange" || r.Header.Get("Authorization") != "Bearer fake-sa-key" {
			t.Error("unexpected SA exchange")
		}
		exchanges.Add(1)
		return reply(r, 200, `{"access_token":"cached-sat","token_type":"Bearer","expires_in":3600}`), nil
	})}})
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer cached-sat" || r.Header.Get("X-Middleware") != "present" {
			t.Error("HTTP/WS authentication or middleware differed")
		}
		if strings.HasSuffix(r.URL.Path, "/models") {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `{"data":[]}`)
			return
		}
		ws, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer ws.CloseNow()
		writeRealtime(t, ws, fixture[1])
		_, _, _ = ws.Read(context.Background())
	}))
	defer server.Close()
	client := forward.NewClient(option.WithCredential(credential), option.WithHTTPClient(server.Client()), option.WithBaseURL(server.URL+"/api/v1/forward"), option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		r = r.Clone(r.Context())
		r.Header.Set("X-Middleware", "present")
		return next(r)
	}))
	if _, err := client.Models.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	conn := realtimeConnect(t, client, forward.RealtimeConnectParams{})
	_ = conn.Close()
	if exchanges.Load() != 1 {
		t.Fatalf("REST/WS did not reuse shared SAT cache: exchanges=%d", exchanges.Load())
	}
}

func TestRealtimeConcurrentWritesAndGracefulCapability(t *testing.T) {
	fixture := realtimeFixtures(t)
	received := make(chan int, 1)
	server := realtimeServer(t, func(ws *websocket.Conn, _ *http.Request) {
		ready := fixture[1]
		ready["payload"].(map[string]any)["capabilities"] = map[string]any{"graceful_close": false}
		writeRealtime(t, ws, ready)
		seen := map[string]bool{}
		for i := 0; i < 32; i++ {
			event := readRealtime(t, ws)
			if event == nil {
				return
			}
			text := event["payload"].(map[string]any)["text"].(string)
			seen[text] = true
		}
		received <- len(seen)
		_, _, _ = ws.Read(context.Background())
	})
	conn := realtimeConnect(t, realtimeClient(server.URL), forward.RealtimeConnectParams{})
	var writers sync.WaitGroup
	for i := 0; i < 32; i++ {
		writers.Add(1)
		go func(i int) {
			defer writers.Done()
			if err := conn.SendText(context.Background(), statusName(i)); err != nil {
				t.Error(err)
			}
		}(i)
	}
	writers.Wait()
	if count := <-received; count != 32 {
		t.Fatalf("writes were lost/corrupted: %d", count)
	}
	if _, err := conn.CloseGracefully(context.Background()); err == nil || !strings.Contains(err.Error(), "does not support") {
		t.Fatalf("unsupported graceful close was accepted: %v", err)
	}
}

func TestRealtimeLifetimeContextAndCloseCode(t *testing.T) {
	for _, kind := range []string{"caller_cancel", "peer_close"} {
		t.Run(kind, func(t *testing.T) {
			fixture := realtimeFixtures(t)
			release := make(chan struct{})
			server := realtimeServer(t, func(ws *websocket.Conn, _ *http.Request) {
				writeRealtime(t, ws, fixture[1])
				<-release
				if kind == "peer_close" {
					_ = ws.Close(websocket.StatusServiceRestart, "restart")
				} else {
					_, _, _ = ws.Read(context.Background())
				}
			})
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			client := realtimeClient(server.URL)
			conn, err := client.Realtime.Connect(ctx, "conv_test", forward.RealtimeConnectParams{})
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			_, _ = conn.Recv(context.Background())
			close(release)
			if kind == "caller_cancel" {
				cancel()
			}
			waitCtx, waitCancel := context.WithTimeout(context.Background(), time.Second)
			defer waitCancel()
			_, err = conn.Recv(waitCtx)
			if kind == "caller_cancel" {
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("lifetime context ignored: %v", err)
				}
			} else {
				var transport *forward.RealtimeTransportError
				if !errors.As(err, &transport) || transport.CloseCode != 1012 || transport.Reason != "restart" || transport.ConversationID != "conv_test" {
					t.Fatalf("close diagnostics lost: %v", err)
				}
			}
		})
	}
}
