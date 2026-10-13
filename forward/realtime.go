package forward

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/coder/websocket"
)

const RealtimeMaxClientMessageBytes = 256 * 1024

type RealtimeService struct {
	Options       []option.RequestOption
	Conversations RealtimeConversationService
}

func NewRealtimeService(opts ...option.RequestOption) RealtimeService {
	return RealtimeService{Options: slices.Clone(opts), Conversations: NewRealtimeConversationService(opts...)}
}

// RealtimeConnectParams controls per-stage timeouts and bounded buffering.
// Zero values select defaults; negative values are rejected. The caller's ctx
// spans the connection lifetime, independently of ordinary HTTP timeouts.
type RealtimeConnectParams struct {
	HandshakeTimeout time.Duration // Default 30s.
	ReadyTimeout     time.Duration // Default 60s.
	SendTimeout      time.Duration // Default 30s.
	CloseTimeout     time.Duration // Default 10s.
	QueueCapacity    int           // Default 256 events.
	QueueBytes       int64         // Default 4 MiB.
	ReadLimit        int64         // Default 1 MiB per server message.
}

func (p RealtimeConnectParams) defaults() (RealtimeConnectParams, error) {
	if p.HandshakeTimeout < 0 || p.ReadyTimeout < 0 || p.SendTimeout < 0 || p.CloseTimeout < 0 || p.QueueCapacity < 0 || p.QueueBytes < 0 || p.ReadLimit < 0 {
		return p, fmt.Errorf("Realtime timeouts and limits must not be negative")
	}
	if p.HandshakeTimeout == 0 {
		p.HandshakeTimeout = 30 * time.Second
	}
	if p.ReadyTimeout == 0 {
		p.ReadyTimeout = 60 * time.Second
	}
	if p.SendTimeout == 0 {
		p.SendTimeout = 30 * time.Second
	}
	if p.CloseTimeout == 0 {
		p.CloseTimeout = 10 * time.Second
	}
	if p.QueueCapacity == 0 {
		p.QueueCapacity = 256
	}
	if p.QueueBytes == 0 {
		p.QueueBytes = 4 * 1024 * 1024
	}
	if p.ReadLimit == 0 {
		p.ReadLimit = 1024 * 1024
	}
	return p, nil
}

// Connect upgrades once and waits for a validated voice.ready. It never follows
// redirects or automatically reconnects/replays input. Events received while
// waiting for ready, including ready itself, remain available through Recv.
func (r *RealtimeService) Connect(ctx context.Context, conversationID string, params RealtimeConnectParams, opts ...option.RequestOption) (*RealtimeConnection, error) {
	if strings.TrimSpace(conversationID) == "" {
		return nil, fmt.Errorf("missing required conversation_id parameter")
	}
	params, err := params.defaults()
	if err != nil {
		return nil, err
	}
	handshakeCtx, cancel := context.WithTimeout(ctx, params.HandshakeTimeout)
	defer cancel()
	cfg, err := convention.NewRequestConfig(handshakeCtx, http.MethodGet, "realtime", nil, nil, slices.Concat(r.Options, opts)...)
	if err != nil {
		return nil, err
	}
	client, request, err := cfg.WebSocketHandshake()
	if err != nil {
		return nil, err
	}
	query := request.URL.Query()
	query.Set("conversation_id", conversationID)
	request.URL.RawQuery = query.Encode()
	ws, response, err := websocket.Dial(handshakeCtx, request.URL.String(), &websocket.DialOptions{HTTPClient: client, HTTPHeader: request.Header})
	cancel()
	if cfg.ResponseInto != nil {
		*cfg.ResponseInto = response
	}
	if err != nil {
		if response != nil && response.StatusCode != http.StatusSwitchingProtocols {
			api := &apierror.Error{StatusCode: response.StatusCode, Response: response, Request: response.Request, RequestID: response.Header.Get("X-Request-Id")}
			if response.Body != nil {
				body, _ := io.ReadAll(response.Body)
				_ = response.Body.Close()
				_ = api.UnmarshalJSON(body)
			}
			return nil, api
		}
		return nil, err
	}
	lifetime, cancelLifetime := context.WithCancel(ctx)
	c := &RealtimeConnection{ConversationID: conversationID, ws: ws, params: params, ctx: lifetime, cancel: cancelLifetime, events: make(chan realtimeQueuedEvent, params.QueueCapacity), done: make(chan struct{}), readyDone: make(chan struct{}), sendGate: make(chan struct{}, 1), closeGate: make(chan struct{}, 1)}
	if response != nil {
		c.RequestID = response.Header.Get("X-Request-Id")
	}
	ws.SetReadLimit(params.ReadLimit)
	go c.receiveLoop()
	readyCtx, cancelReady := context.WithTimeout(ctx, params.ReadyTimeout)
	defer cancelReady()
	select {
	case <-c.readyDone:
		c.mu.Lock()
		isReady := c.ready
		c.mu.Unlock()
		if !isReady {
			return nil, c.Err()
		}
		return c, nil
	case <-c.done:
		return nil, c.Err()
	case <-readyCtx.Done():
		err := fmt.Errorf("qca realtime ready wait: %w", readyCtx.Err())
		c.finish(err)
		return nil, err
	}
}

type realtimeQueuedEvent struct {
	event RealtimeEvent
	bytes int64
}

// RealtimeConnection has one receiver and serializes sends. Ready is immutable
// after Connect returns. Applications own audio playback and its receipts.
type RealtimeConnection struct {
	ConversationID string
	RequestID      string
	Ready          RealtimeReadyPayload
	ws             *websocket.Conn
	params         RealtimeConnectParams
	ctx            context.Context
	cancel         context.CancelFunc
	events         chan realtimeQueuedEvent
	done           chan struct{}
	readyDone      chan struct{}
	sendGate       chan struct{}
	closeGate      chan struct{}
	mu             sync.Mutex
	ready          bool
	queuedBytes    int64
	err            error
	closing        bool
	closeRequestID string
	closeResult    *RealtimeClosedPayload
	closeWaiter    chan RealtimeClosedPayload
	finishOnce     sync.Once
}

type RealtimeTransportError struct {
	ConversationID string
	CloseCode      int
	Reason         string
	Err            error
}

func (e *RealtimeTransportError) Error() string {
	return fmt.Sprintf("qca realtime transport closed (conversation_id=%s, code=%d): %s", e.ConversationID, e.CloseCode, e.Reason)
}
func (e *RealtimeTransportError) Unwrap() error { return e.Err }

func (c *RealtimeConnection) receiveLoop() {
	for {
		kind, b, err := c.ws.Read(c.ctx)
		if err != nil {
			var ce websocket.CloseError
			if errors.As(err, &ce) {
				c.finish(&RealtimeTransportError{ConversationID: c.ConversationID, CloseCode: int(ce.Code), Reason: ce.Reason, Err: err})
			} else {
				c.finish(err)
			}
			return
		}
		if kind != websocket.MessageText {
			c.finish(&RealtimeProtocolError{Message: "server sent a binary frame"})
			return
		}
		event, err := decodeRealtimeEvent(b, c.ConversationID)
		if err != nil {
			c.finish(err)
			return
		}
		var ready *RealtimeReadyPayload
		if event.Type == "voice.ready" {
			var p RealtimeReadyPayload
			var capabilityFields struct {
				Capabilities map[string]json.RawMessage `json:"capabilities"`
			}
			capErr := event.DecodePayload(&capabilityFields)
			graceful := capabilityFields.Capabilities["graceful_close"]
			if capErr != nil || graceful == nil || string(graceful) == "null" {
				c.finish(&RealtimeProtocolError{Message: "voice.ready requires capabilities.graceful_close"})
				return
			}
			if err := event.DecodePayload(&p); err != nil || p.Provider == "" || p.Config.Audio.Output.Voice == "" || p.InputAudio.Format != "pcm16" || p.InputAudio.SampleRate != 16000 || p.InputAudio.Channels != 1 || p.OutputAudio.Format != "pcm16" || p.OutputAudio.SampleRate != 24000 || p.OutputAudio.Channels != 1 {
				c.finish(&RealtimeProtocolError{Message: "invalid voice.ready configuration"})
				return
			}
			ready = &p
		}
		c.mu.Lock()
		if event.Type == "connection.closed" && c.closeWaiter != nil {
			var p RealtimeClosedPayload
			if event.DecodePayload(&p) == nil && p.RequestID == c.closeRequestID {
				select {
				case c.closeWaiter <- p:
				default:
				}
			}
		}
		if c.queuedBytes+int64(len(b)) > c.params.QueueBytes {
			c.mu.Unlock()
			c.finish(&RealtimeProtocolError{Message: "receive queue byte limit exceeded"})
			return
		}
		c.queuedBytes += int64(len(b))
		c.mu.Unlock()
		select {
		case c.events <- realtimeQueuedEvent{event: event, bytes: int64(len(b))}:
		default:
			c.finish(&RealtimeProtocolError{Message: "receive queue event limit exceeded"})
			return
		}
		if ready != nil {
			c.mu.Lock()
			if !c.ready {
				c.Ready = *ready
				c.ready = true
				close(c.readyDone)
			}
			c.mu.Unlock()
		}
		if event.Type == "voice.replaced" {
			c.finish(&RealtimeProtocolError{Message: "connection replaced"})
			return
		}
		if event.Type == "error" {
			var p RealtimeErrorPayload
			if err := event.DecodePayload(&p); err != nil {
				c.finish(&RealtimeProtocolError{Message: "invalid error payload"})
				return
			}
			c.mu.Lock()
			isReady := c.ready
			c.mu.Unlock()
			if !isReady || p.Code == "provider_initialization_failed" || p.Code == "context_restore_failed" || p.Code == "provider_unavailable" || p.Code == "service_restarting" || p.Code == "event_persistence_failed" {
				c.finish(&RealtimeServerError{RealtimeErrorPayload: p, ConversationID: c.ConversationID})
				return
			}
		}
	}
}

func (c *RealtimeConnection) finish(err error) {
	c.finishOnce.Do(func() {
		c.mu.Lock()
		c.err = err
		c.mu.Unlock()
		close(c.done)
		c.cancel()
		_ = c.ws.CloseNow()
	})
}

// Err returns the terminal receive/transport error, if any.
func (c *RealtimeConnection) Err() error { c.mu.Lock(); defer c.mu.Unlock(); return c.err }

// Recv preserves receive order and drains buffered events before reporting a
// terminal error. Cancelling this wait does not close the connection.
func (c *RealtimeConnection) Recv(ctx context.Context) (RealtimeEvent, error) {
	dequeue := func(q realtimeQueuedEvent) (RealtimeEvent, error) {
		c.mu.Lock()
		c.queuedBytes -= q.bytes
		c.mu.Unlock()
		return q.event, nil
	}
	select {
	case q := <-c.events:
		return dequeue(q)
	default:
	}
	select {
	case q := <-c.events:
		return dequeue(q)
	case <-ctx.Done():
		return RealtimeEvent{}, ctx.Err()
	case <-c.done:
		select {
		case q := <-c.events:
			return dequeue(q)
		default:
		}
		if err := c.Err(); err != nil {
			return RealtimeEvent{}, err
		}
		return RealtimeEvent{}, io.EOF
	}
}

func (c *RealtimeConnection) SendEvent(ctx context.Context, event RealtimeClientEvent) error {
	return c.sendEvent(ctx, event, false)
}
func (c *RealtimeConnection) sendEvent(ctx context.Context, event RealtimeClientEvent, allowClosing bool) error {
	if event.Version == "" {
		event.Version = RealtimeProtocolVersion
	}
	if event.Version != RealtimeProtocolVersion || event.Type == "" {
		return &RealtimeProtocolError{Message: "invalid client event envelope"}
	}
	if event.Payload == nil {
		event.Payload = map[string]any{}
	}
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if len(b) > RealtimeMaxClientMessageBytes {
		return &RealtimeProtocolError{Message: "client message exceeds 256 KiB"}
	}
	sendCtx, cancel := context.WithTimeout(ctx, c.params.SendTimeout)
	defer cancel()
	select {
	case c.sendGate <- struct{}{}:
		defer func() { <-c.sendGate }()
	case <-sendCtx.Done():
		return sendCtx.Err()
	case <-c.done:
		if err := c.Err(); err != nil {
			return err
		}
		return io.ErrClosedPipe
	}
	select {
	case <-c.done:
		if err := c.Err(); err != nil {
			return err
		}
		return io.ErrClosedPipe
	default:
	}
	c.mu.Lock()
	closing := c.closing
	c.mu.Unlock()
	if closing && !allowClosing {
		return &RealtimeProtocolError{Message: "connection is closing"}
	}
	if err = c.ws.Write(sendCtx, websocket.MessageText, b); err != nil {
		c.finish(err)
	}
	return err
}
func (c *RealtimeConnection) SendText(ctx context.Context, text string) error {
	if !utf8.ValidString(text) {
		return &RealtimeProtocolError{Message: "text must be valid UTF-8"}
	}
	text = strings.TrimSpace(text)
	if text == "" || len([]byte(text)) > 16*1024 {
		return &RealtimeProtocolError{Message: "text must be nonempty valid UTF-8 and at most 16 KiB"}
	}
	return c.SendEvent(ctx, RealtimeClientEvent{Type: "text.message", Payload: map[string]any{"text": text}})
}

// AppendAudio sends nonempty 16 kHz mono PCM16 little-endian bytes. It does not
// resample, decode WAV or send binary WebSocket frames.
func (c *RealtimeConnection) AppendAudio(ctx context.Context, pcm []byte) error {
	if len(pcm) == 0 || len(pcm)%2 != 0 {
		return &RealtimeProtocolError{Message: "audio must contain whole nonempty PCM16 samples"}
	}
	if base64.StdEncoding.EncodedLen(len(pcm)) > RealtimeMaxClientMessageBytes {
		return &RealtimeProtocolError{Message: "client message exceeds 256 KiB"}
	}
	return c.SendEvent(ctx, RealtimeClientEvent{Type: "audio.append", Payload: map[string]any{"audio": base64.StdEncoding.EncodeToString(pcm)}})
}
func (c *RealtimeConnection) Interrupt(ctx context.Context) error {
	return c.SendEvent(ctx, RealtimeClientEvent{Type: "interrupt"})
}
func (c *RealtimeConnection) Ping(ctx context.Context) error {
	return c.SendEvent(ctx, RealtimeClientEvent{Type: "ping"})
}
func (c *RealtimeConnection) playback(ctx context.Context, kind string, receipt RealtimePlaybackReceipt) error {
	payload := map[string]any{}
	if receipt.WorkID != "" {
		payload["work_id"] = receipt.WorkID
	}
	if receipt.AnnouncementID != "" {
		payload["announcement_id"] = receipt.AnnouncementID
	}
	return c.SendEvent(ctx, RealtimeClientEvent{Type: kind, Payload: payload})
}
func (c *RealtimeConnection) PlaybackStarted(ctx context.Context, r RealtimePlaybackReceipt) error {
	return c.playback(ctx, "playback.started", r)
}
func (c *RealtimeConnection) PlaybackEnded(ctx context.Context, r RealtimePlaybackReceipt) error {
	return c.playback(ctx, "playback.ended", r)
}
func (c *RealtimeConnection) PlaybackCancelled(ctx context.Context, r RealtimePlaybackReceipt) error {
	return c.playback(ctx, "playback.cancelled", r)
}

// CloseGracefully sends a unique close request, waits for matching persistence
// acknowledgement, then closes the transport. Ordinary Close does not confirm
// persistence. Repeated successful calls return the same acknowledgement.
func (c *RealtimeConnection) CloseGracefully(ctx context.Context) (*RealtimeClosedPayload, error) {
	closeCtx, cancel := context.WithTimeout(ctx, c.params.CloseTimeout)
	defer cancel()
	select {
	case c.closeGate <- struct{}{}:
		defer func() { <-c.closeGate }()
	case <-closeCtx.Done():
		return nil, closeCtx.Err()
	}
	c.mu.Lock()
	if c.closeResult != nil {
		res := *c.closeResult
		c.mu.Unlock()
		return &res, nil
	}
	c.mu.Unlock()
	defer c.Close()
	if !c.Ready.Capabilities.GracefulClose {
		return nil, &RealtimeProtocolError{Message: "server does not support graceful close"}
	}
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return nil, err
	}
	requestID := hex.EncodeToString(id[:])
	c.mu.Lock()
	c.closing = true
	c.closeRequestID = requestID
	c.closeWaiter = make(chan RealtimeClosedPayload, 1)
	waiter := c.closeWaiter
	c.mu.Unlock()
	if err := c.sendEvent(closeCtx, RealtimeClientEvent{Type: "connection.close", Payload: map[string]any{"request_id": requestID}}, true); err != nil {
		return nil, err
	}
	select {
	case result := <-waiter:
		if result.Outcome != "no_active_response" && result.Outcome != "no_content" && result.Outcome != "saved" && result.Outcome != "saved_interrupted" {
			return nil, &RealtimeProtocolError{Message: "invalid connection.closed outcome"}
		}
		c.mu.Lock()
		c.closeResult = &result
		c.mu.Unlock()
		return &result, nil
	case <-closeCtx.Done():
		return nil, fmt.Errorf("qca realtime close acknowledgement: %w", closeCtx.Err())
	case <-c.done:
		select {
		case result := <-waiter:
			if result.Outcome == "no_active_response" || result.Outcome == "no_content" || result.Outcome == "saved" || result.Outcome == "saved_interrupted" {
				c.mu.Lock()
				c.closeResult = &result
				c.mu.Unlock()
				return &result, nil
			}
		default:
		}
		if err := c.Err(); err != nil {
			return nil, err
		}
		return nil, io.ErrClosedPipe
	}
}

// Close releases transport resources without confirming history persistence.
func (c *RealtimeConnection) Close() error { c.finish(nil); return nil }
