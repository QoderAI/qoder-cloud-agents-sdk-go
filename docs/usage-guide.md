# Usage guide

Resource-specific behavior and examples for the Qoder Cloud Agents SDK. See the [README](../README.md) for installation, authentication and shared client configuration. Run shell commands from the repository root; code fragments use the clients and imports shown in the README unless stated otherwise.

- [API coverage](#api-coverage)
- [SDK-created Session metadata](#sdk-created-session-metadata)
- [Forward SAT management](#forward-sat-management)
- [Usage and cancellation](#usage-and-cancellation)
- [Conversations](#conversations)
- [System prompts and tools](#system-prompts-and-tools)
- [Streaming](#streaming)
- [Long-running operations](#long-running-operations)
- [Forward Decisions (Beta)](#forward-decisions-beta)
- [Forward Drive (Beta)](#forward-drive-beta)
- [File uploads and downloads](#file-uploads-and-downloads)
- [Examples](#examples)
- [Forward Webhooks](#forward-webhooks)
- [Resource search](#resource-search)
- [Official contract alignment](#official-contract-alignment)

## API coverage

| Mode | Package | Resources |
| --- | --- | --- |
| Forward | `forward` | Identity, Template, Session, Schedule, Batch, Channel, Environment, File, Skill, Vault, MemoryStore, Model, Usage, Realtime, Decision |
| Managed | `managed` | Agent, Session, Deployment, Dream, Environment, File, Skill, Vault, MemoryStore, Model |

## SDK-created Session metadata

Direct Forward and Managed Session creation automatically adds these string-valued
metadata fields:

| Field | Value |
| --- | --- |
| `__qca_sdk_created` | `"true"` |
| `__qca_sdk_language` | `"go"` |
| `__qca_sdk_version` | The installed SDK's version |

Reserve the `__qca_sdk_` prefix for SDK provenance. Business metadata is preserved,
and the caller's metadata object is not mutated. Omitted or null metadata becomes
an object containing the markers. Passing an existing marker with the exact same
string value is accepted; a different value or type (including null) fails locally
before any HTTP request. This also applies after request body/JSON options are merged.

Session updates, reads and event operations do not inject or repair these fields.
Sessions created indirectly through Schedules, Batches or Realtime are outside this
behavior. These markers describe the creating SDK; they are client-supplied metadata,
not an authentication or authorization mechanism.

Migration: rename any business fields that already use these three keys before
upgrading. Automatic markers count toward API metadata limits; the SDK never drops
business fields to make room. No additional parameters or opt-in are required.

## Forward SAT management

Forward exposes explicit admin/Identity SAT issuance and revocation in addition
to automatic platform exchange. These three endpoints authenticate with an
**SA Key**, supplied per method call; ordinary business requests still use PAT
or SAT. The Key overrides client/request Authorization only for that call, is
never included in the JSON body, and is not saved on the client.

`ttl_seconds` is optional (1-43200); omission uses the issuer's default. Issuance
and revocation are not automatically retried, even when the client or request
allows retries. Do not blindly repeat an issuance after an uncertain network
failure: it may already have issued a token. Metadata is for audit information,
not secrets. Keep the issuance response's Token ID, Credential ID and expiry for
precise revocation. A returned static SAT does not renew automatically, and
revocation affects its Forward binding; it does not cancel an already-running
Session or revoke unrelated Managed/platform access. Never log the SA Key, SAT,
serialized token response or Authorization headers.

Run `make test-live-sa-tokens` with `QODER_SA_CN_KEY`,
`QODER_SA_GLOBAL_KEY` and `QODER_SA_LIVE_ALLOW_WRITE=true` to verify both regions.
This separate suite issues admin and Identity SATs, checks business access and
cross-Identity refusal, then revokes the tokens, verifies access is refused and
deletes its temporary Identities. It does not execute a model or require
`QODER_SA_LIVE_ALLOW_EXECUTION=true`.

```go
issuer := forward.NewClient(option.WithBaseURL(forwardBaseURL))
// Admin: issuer.ServiceAccountTokens.New(ctx, the same parameters).
token, err := issuer.Identities.ServiceAccountTokens.New(ctx, identityID,
    forward.ServiceAccountTokenNewParams{SAKey: saKey, TTLSeconds: forward.Int(600)})
if err != nil { return err }
defer func() {
    // Use an independent cleanup context if the business context has expired.
    _, revokeErr := issuer.ServiceAccountTokens.Revoke(context.Background(), token.AuthTokenID,
        forward.ServiceAccountTokenRevokeParams{
            SAKey: saKey,
            TargetCredentialID: forward.String(token.CredentialID),
            TargetExpiresAt: param.NewOpt(token.ExpiresAt),
            Reason: forward.String("runtime complete"),
        })
    // Report cleanup failures; do not log credentials.
    if revokeErr != nil { log.Printf("Forward SAT cleanup failed (%T)", revokeErr) }
}()
runtime := forward.NewClient(option.WithBaseURL(forwardBaseURL),
    option.WithAuthMode("sa"), option.WithAuthToken(token.AccessToken))
page, err := runtime.Sessions.List(ctx, forward.SessionListParams{
    IdentityIDs: []string{identityID}, Limit: forward.Int(1),
})
```

Runnable example: `go run ./examples/auth/forward-tokens`. It revokes its token
even if the business call fails. Set `QODER_SA_KEY` and `QODER_SA_REGION=cn|global`;
set `QODER_IDENTITY_ID` for Identity mode, or leave it unset for admin mode.

See the official [admin issuance](https://docs.qoder.com/cloud-agents/api/forward/service-account-tokens/create),
[Identity issuance](https://docs.qoder.com/cloud-agents/api/forward/service-account-tokens/create-identity)
and [revocation](https://docs.qoder.com/cloud-agents/api/forward/service-account-tokens/revoke) contracts.

## Usage and cancellation

Forward Usage requires PAT or Admin SAT.

```go
usage, err := client.Usage.ListIdentities(ctx, forward.UsageListParams{
    StartAt: "2026-09-14T09:00:00",
    EndAt: "2026-09-14T12:00:00",
    IdentityIDs: []string{"idn_one", "idn_two"},
})
if err != nil {
    return err
}
for _, item := range usage.Data {
    fmt.Println(item.IdentityID, item.ActiveSeconds, item.Credits)
}
```

`Usage.ListTemplates` uses the same filters. Both bounds are whole hours in
Asia/Shanghai, including Global: the start is inclusive, the end exclusive, and
the maximum window is 744 hours. Use `ActiveSeconds` without truncating fractions.
Legacy `start_time`, `end_time`, and `duration_seconds` are not exposed.

Rotate Forward secrets with `client.Vaults.Credentials.Update`; only `auth` and
`metadata` are patched, and omitted fields are preserved. Use
`param.NullMap[map[string]any]()` to clear metadata, or a map value of `nil` to remove
a key. Updates never retry automatically; secret values are not returned by GET.

Managed adds `client.Sessions.Cancel(ctx, sessionID, managed.SessionCancelParams{})`,
which returns a lightweight `canceling` acknowledgement (HTTP 202 for active work,
200 for an idle no-op). It also exposes `client.Deployments.Runs.List` and `Get`,
with the deployment ID as a path argument; the existing global `DeploymentRuns`
resource remains available.

## Conversations

Forward creates a Session from an Identity and a Template. The fragments below assume both already exist, with a usable execution environment configured on the Template; [the Forward session example](../examples/forward/session/main.go) shows how to create them.

```go
func startForwardSession(ctx context.Context, client forward.Client, identityID, templateID string) (*forward.Session, error) {
	return client.Sessions.New(ctx, forward.SessionNewParams{
		IdentityID: identityID,
		TemplateID: templateID,
		Title:      forward.String("project assistant"),
	})
}

func sendForwardMessage(ctx context.Context, client forward.Client, sessionID, text, requestKey string) (string, error) {
	result, err := client.Sessions.Events.Send(ctx, sessionID, forward.SessionEventSendParams{
		Events: []forward.SessionEventParam{{
			Type: "user.message",
			Content: forward.EventContentUnionParam{OfBlocks: []forward.ContentBlockParam{{
				Type: "text", Text: forward.String(text),
			}}},
		}},
		IdempotencyKey: forward.String(requestKey),
	})
	if err != nil {
		return "", err
	}
	if len(result.Data) != 1 {
		return "", fmt.Errorf("expected one user event, got %d", len(result.Data))
	}
	return result.Data[0].ID, nil
}
```

Managed creates a Session from an existing Agent and Environment; see [the Managed session example](../examples/managed/session/main.go).

```go
func startManagedSession(ctx context.Context, client managed.Client, agentID, environmentID string) (*managed.ManagedAgentsSession, error) {
	return client.Sessions.New(ctx, managed.SessionNewParams{
		EnvironmentID: environmentID,
		Agent: managed.SessionNewParamsAgentUnion{
			OfString: managed.String(agentID),
		},
	})
}

func sendManagedMessage(ctx context.Context, client managed.Client, sessionID, text, requestKey string) (string, error) {
	result, err := client.Sessions.Events.Send(ctx, sessionID, managed.SessionEventSendParams{
		Events: []managed.ManagedAgentsEventParamsUnion{{
			OfUserMessage: &managed.ManagedAgentsUserMessageEventParams{
				Type: "user.message",
				Content: []managed.ManagedAgentsUserMessageEventParamsContentUnion{{
					OfText: &managed.ManagedAgentsTextBlockParam{Type: "text", Text: text},
				}},
			},
		}},
	}, option.WithHeader("Idempotency-Key", requestKey))
	if err != nil {
		return "", err
	}
	if len(result.Data) != 1 {
		return "", fmt.Errorf("expected one user event, got %d", len(result.Data))
	}
	return result.Data[0].ID, nil
}
```

`requestKey` is a non-empty unique key that the caller generates and stores for one logical message: reuse it when retrying that message, and use a fresh key for a new one. The return value is the user message event ID, which is also the right starting point for reading the events that follow.

Reuse `sessionID` to continue the same conversation. A successful `Events.Send` only means the message was accepted; the agent's answer and execution status arrive through the event list or the stream.

## System prompts and tools

The system prompt, model and tools live on the Forward Template or the Managed Agent. Take `modelID` from `Models.List` on the same account and check `IsEnabled` first.

```go
func createForwardTemplate(ctx context.Context, client forward.Client, environmentID, modelID string) (*forward.Template, error) {
	return client.Templates.New(ctx, forward.TemplateNewParams{
		Name:          "project-assistant",
		EnvironmentID: environmentID,
		Model:         forward.ModelConfigUnionParam{OfString: forward.String(modelID)},
		System:        forward.String("Answer from the project material you can read."),
		Tools:         []forward.ToolParam{{Type: "agent_toolset_20260401"}},
	})
}

func createManagedAgent(ctx context.Context, client managed.Client, modelID string) (*managed.ManagedAgentsAgent, error) {
	return client.Agents.New(ctx, managed.AgentNewParams{
		Name:   "project-assistant",
		Model:  managed.ManagedAgentsModelConfigParams{ID: modelID},
		System: managed.String("Answer from the project material you can read."),
		Tools: []managed.AgentNewParamsToolUnion{{
			OfAgentToolset20260401: &managed.ManagedAgentsAgentToolset20260401Params{
				Type: "agent_toolset_20260401",
			},
		}},
	})
}
```

Built-in tools run in the cloud runtime, and each call produces tool events. Note that the Managed fragment sets `Type` on every event, content block and tool: request construction does not fill those discriminators in for you.

## Streaming

Both modes stream through `Sessions.Events.StreamEvents`. Pass the ID of the user message you just sent as the starting point, so events from earlier turns are not mistaken for the current one.

```go
func readForwardTurn(ctx context.Context, client forward.Client, sessionID, afterID string) error {
	stream := client.Sessions.Events.StreamEvents(ctx, sessionID, forward.SessionEventStreamParams{
		LastEventID: forward.String(afterID),
	})
	defer stream.Close()

	for stream.Next() {
		event := stream.Current()
		switch event.Type {
		case "agent.message":
			blocks, err := event.ContentBlocks()
			if err != nil {
				return err
			}
			for _, block := range blocks {
				if block.Type == "text" {
					fmt.Println(block.Text)
				}
			}
		case "session.status_idle":
			fmt.Printf("turn paused or finished: %v\n", event.StopReason)
			return nil
		case "session.error", "session.status_terminated":
			return fmt.Errorf("session stopped: type=%s event_id=%s", event.Type, event.ID)
		}
	}
	if err := stream.Err(); err != nil {
		return err
	}
	return fmt.Errorf("stream ended before an idle event; last_event_id=%s", stream.LastEventID())
}
```

Managed takes the starting point as a request header and exposes message content through a response union:

```go
func readManagedTurn(ctx context.Context, client managed.Client, sessionID, afterID string) error {
	stream := client.Sessions.Events.StreamEvents(ctx, sessionID, managed.SessionEventStreamParams{
		LastEventID: managed.String(afterID),
	})
	defer stream.Close()

	for stream.Next() {
		event := stream.Current()
		switch event.Type {
		case "agent.message":
			for _, block := range event.AsAgentMessage().Content {
				if block.Type == "text" {
					fmt.Println(block.Text)
				}
			}
		case "session.status_idle":
			fmt.Printf("turn paused or finished: %s\n", event.StopReason.Type)
			return nil
		case "session.error", "session.status_terminated":
			return fmt.Errorf("session stopped: type=%s event_id=%s", event.Type, event.ID)
		}
	}
	if err := stream.Err(); err != nil {
		return err
	}
	return fmt.Errorf("stream ended before an idle event; last_event_id=%s", stream.LastEventID())
}
```

Close the stream, and check `Err()` once the loop ends. An idle event does not always mean success: it can also mean the run is waiting on the caller or hit a budget, so judge the outcome from `stop_reason` together with the final answer. [The scenario helper](../examples/internal/live/turn.go) shows a full set of execution assertions.

To render text as it is produced, set `EventDeltas` on the stream params: `[]string{"agent.message"}` for Forward, `[]managed.ManagedAgentsDeltaType{"agent.message"}` for Managed. You then also receive `event_start` and `event_delta` previews, and the final complete event repeats the whole content. Update one message in place instead of appending both the deltas and the final content, and do not deduplicate deltas by event ID alone — one ID legitimately produces many delta events.

`StreamEvents` is a one-shot connection. For automatic recovery, call `NewResumableStream` with the same arguments and iterate through the same `Next`, `Current`, `Err`, `LastEventID` and `Close` methods:

```go
stream := client.Sessions.Events.NewResumableStream(ctx, sessionID, forward.SessionEventStreamParams{
	EventDeltas: []string{"agent.message"},
	LastEventID: forward.String(afterID),
})
defer stream.Close()
for stream.Next() {
	handle(stream.Current())
}
if err := stream.Err(); err != nil {
	return err
}
```

The resumable stream checkpoints `LastEventID` only after yielding a decoded event, then reconnects transport/read failures and unexpected EOF with cancellable jittered exponential backoff from 500 ms to 10 seconds. Backoff resets only after an event arrives on a connection that has remained healthy for more than five seconds. It preserves stream params and request options on every connection. HTTP retry classification matches ordinary SDK requests: 408, 429, 5xx and transport failures retry (subject to `x-should-retry`), while 409 and other non-retryable 4xx errors stop. It neither clears an invalid cursor nor queries event history. Events are not deduplicated because multiple delta frames can legitimately share one ID. A yielded `session.status_terminated` or `session.deleted` event ends the stream normally; context cancellation and `Close` both stop an active connection or backoff wait.


### Thread event resumption

Forward and Managed also support automatic resumption for a single Thread. Use the
same retry and cancellation behavior as Session streams, with a thread-scoped
`Last-Event-ID` cursor:

```go
// Forward
stream := forwardClient.Sessions.Threads.Events.NewResumableStream(ctx, sessionID, threadID, forward.SessionThreadEventStreamParams{
    LastEventID: forward.String(afterID),
})
defer stream.Close()
for stream.Next() {
    handle(stream.Current())
}
if err := stream.Err(); err != nil {
    return err
}

// Managed uses the parent session in params.
managedStream := managedClient.Sessions.Threads.Events.NewResumableStream(ctx, threadID, managed.SessionThreadEventStreamParams{
    SessionID: sessionID,
    LastEventID: managed.String(afterID),
})
defer managedStream.Close()
```

`session.thread_status_idle` ends only the current turn, so the subscription stays
open. `session.thread_status_terminated`, `session.status_terminated`, or
`session.deleted` ends it normally. A non-retryable error (including an invalid
cursor or missing thread) is returned to the caller; the SDK never clears the
cursor to silently start over. The existing one-shot Thread stream remains
available. See the [Forward Thread stream](https://docs.qoder.com/cloud-agents/api/forward/sessions/stream-thread-events)
and [Managed Thread stream](https://docs.qoder.com/cloud-agents/api/sessions/threads/events/stream)
contracts for cursor validity.

## Long-running operations

For Session, Schedule, Batch, Deployment and Dream, the response to a create or trigger call is a different stage from the execution finishing. Confirm the outcome through the event stream, the run status or the task result.

| Case | Where the result arrives |
| --- | --- |
| Session | the `Sessions.Events` stream or event list, with the idle `stop_reason` and the final answer |
| Forward Schedule | `ScheduleRuns`, plus the linked Session's result |
| Forward Batch | `Batches.Get`, `Batches.Tasks.List`, `Batches.GetOutput` |
| Managed Deployment | `DeploymentRuns`, plus the linked Session's result |
| Managed Dream | the Dream status, plus the contents of the output Memory Store |

Server-side windows and local deadlines are independent: a Batch `completion_window: "24h"` is submitted to the server, while the examples' `-timeout 5m` is only how long the local program waits. A Batch sitting in `queued` may not have started, and waiting longer does not change how the server schedules it. The examples cancel the work they created when they time out, so that run's output cannot be used to judge later execution.

## Forward Decisions (Beta)

The Decision API is in Beta. Its behavior may change in future updates.

Evaluate a nonempty `state` object or an accessible `session_id` (exactly one)
with 1–32 choice or score questions. Choice values are strings or booleans;
score questions select from your numeric levels. Creation waits for the result.
The service selects the model; no `model` parameter is accepted. Normal PAT or
SAT authentication determines the caller's access.

```go
// client is a configured forward.Client.
decision, err := client.Decisions.New(ctx, forward.DecisionNewParams{
    State: map[string]any{"message": "My issue is resolved."},
    Questions: []forward.DecisionQuestionParam{{
        Type: "choice", Name: forward.String("resolved"),
        Instructions: "Has the user confirmed resolution?",
        Choices: []forward.DecisionChoiceParam{{Value: true}, {Value: false}},
    }},
    IdempotencyKey: forward.String("resolution-check-unique-request-id"),
})
if err != nil { return err }
stored, err := client.Decisions.Get(ctx, decision.ID)
if err != nil { return err }
fmt.Println(stored.Answers)
// Or use ListAutoPaging for all pages.
page, err := client.Decisions.List(ctx, forward.DecisionListParams{Limit: forward.Int(20)})
if err != nil { return err }
fmt.Println(page.Data)
```

For Session input, replace `State` with `SessionID: forward.String(sessionID)`.
For score questions, set `Type: "score"` and `Levels` with labels and scores.
Use `answer.JSON.Score.Valid()` to distinguish an omitted score from zero.

Session context is selected by the server. `context.mode` identifies the input
mode; Session results may include event IDs and counts describing the selected
history. Answers have type `choice` or `score`. Oversized input returns HTTP 422.

Reuse an idempotency key only for identical request bytes. A changed request or
a still-running request with that key returns HTTP 409; the SDK does not retry
that conflict automatically. Creation errors can contain `decision_id` in the
standard API error body. Retrieve and list can return running or failed records.
Lists preserve server order and support `after_id` / `before_id` cursors (one at
a time), optional `session_id`, and a limit of 1–100. There is no delete endpoint;
Decision records expire according to server retention.

See the [runnable example](../examples/forward/decisions/main.go) and the
[opt-in live test instructions](../CONTRIBUTING.md#decision-live-tests).

## Forward Drive (Beta)

Drive is a Beta API. It uses `entries` / `next_page_token`, not the ordinary
resource paginator. Pass `identity_id` even with an Identity SAT, preserve the
server order, and reuse a cursor only within the same Identity, directory and
authentication scope. Paths are relative; the SDK preserves slashes when encoding
the list query and leaves path validation to the service.

Upload/download methods return a signed URL, method, headers and expiry. Send the
storage request with a separate HTTP client using only those returned headers;
do not attach QCA credentials or log signed URLs. Replaying an `Idempotency-Key`
does not renew the URL: use a new key or omit it when requesting a fresh URL.

Deleting a nonempty directory requires `recursive=true`; `clear` deletes the
entire selected Identity Drive. Partial failures retain `drive_delete_partial`
and `deleted_count` in the standard API error body. Direct Drive APIs do not
require Template tool configuration.

```go
page, err := client.Drives.ListEntries(ctx, forward.DriveListEntriesParams{
    IdentityID: identityID, Path: forward.String("reports"), Limit: forward.Int(20),
})
if err != nil { return err }
for _, entry := range page.Entries { fmt.Println(entry.Path) }
if page.NextPageToken != "" {
    page, err = client.Drives.ListEntries(ctx, forward.DriveListEntriesParams{
        IdentityID: identityID, Path: forward.String("reports"),
        PageToken: forward.String(page.NextPageToken),
    })
}
// Replace the complete managed-tool list, retaining any other tools you need.
_, err = client.Templates.Update(ctx, templateID, forward.TemplateUpdateParams{
    ManagedToolConfig: forward.ManagedToolConfigParam{EnabledTools: []string{"drive"}},
})
```

Run the [file lifecycle example](../examples/forward/drive/main.go) with
`go run ./examples/forward/drive -env .env.live -region cn` (or `global`).
It supports the existing PAT/SA configuration and creates a fresh Identity.
`DriveEntry.JSON.Size.Valid()` distinguishes zero-byte files from directories.
For maintainers, `make test-live-drive` tests CN/Global × PAT/admin SAT/Identity SAT.

To enable Agent-managed Drive tools, use Template `managed_tool_config.enabled_tools`.
It is separate from `tools[].enabled_tools` and replaces the entire list; an empty
list disables managed tools, omission preserves the setting, and `null` clears it.
Full Agent Drive usage also requires the DeliverArtifacts tool. Internal
`/forward-mcp` routes are invoked by the runtime and are not public SDK methods.

The Drive live suite requires `QODER_DRIVE_LIVE_ALLOW_WRITE=true`,
`QODER_FORWARD_CN_PAT`, `QODER_FORWARD_GLOBAL_PAT`, `QODER_SA_CN_KEY`, and
`QODER_SA_GLOBAL_KEY` in `.env.live` (or `LIVE_ENV_FILE`). It creates isolated
Identities, verifies storage transfers, pagination, deletion, Identity isolation,
and Template configuration, then clears only those new Identities. It does not
run Agent/model tool execution.

## File uploads and downloads

Upload parameters accept an `io.Reader`. Use `convention.UploadFile` to set the filename and media type:

```go
func uploadText(ctx context.Context, client forward.Client, text string) (*forward.FileMetadata, error) {
	return client.Files.Upload(ctx, forward.FileUploadParams{
		File: convention.UploadFile{
			Reader:    strings.NewReader(text),
			Name:      "notes.txt",
			MediaType: "text/plain",
		},
		Purpose: forward.String("session_resource"),
	})
}
```

A local file works the same way through `os.Open`, closed once the call returns. The Managed equivalent is `client.Files.Upload(ctx, managed.FileUploadParams{File: ...})`, which also takes a `convention.UploadFile`.

Uploading a file and attaching it to a Session are two steps; [the Forward](../examples/forward/resources/main.go) and [Managed](../examples/managed/resources/main.go) resource examples show the attachment. Skill file trees go through `Skills.New` and `Skills.Versions.New`, where each reader in `files` is named with its relative path, such as `my-skill/SKILL.md`.

Downloads return an `*http.Response` whose body the caller reads and closes:

```go
func downloadFile(ctx context.Context, client forward.Client, fileID string, destination io.Writer) error {
	response, err := client.Files.Download(ctx, fileID)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	_, err = io.Copy(destination, response.Body)
	return err
}
```

The Managed method is `client.Files.Download(ctx, fileID, managed.FileDownloadParams{})`. The SDK first asks the API for a temporary download URL, then fetches the stored content. The API PAT, request headers and middleware are not attached to that storage request, though a custom HTTP transport still governs it.

## Examples

### Realtime Beta

`client.Realtime` connects to the public Forward gateway. Realtime must already
be enabled on the Template; this SDK does not change its Realtime configuration.
Creating a Conversation requires a stable `Idempotency-Key` (1–256 UTF-8 bytes)
for that operation. It stays in a header during retries, rather than the JSON body.

```go
conversation, err := client.Realtime.Conversations.New(ctx,
	forward.RealtimeConversationNewParams{
		IdentityID: identityID,
		TemplateID: templateID,
		IdempotencyKey: forward.String(operationKey),
	})
if err != nil { return err }
connection, err := client.Realtime.Connect(ctx, conversation.ID,
	forward.RealtimeConnectParams{})
if err != nil { return err }
defer connection.Close()
// Connect returns after voice.ready; that event and earlier voice.state events
// remain in the receive queue. ctx spans the connection lifetime.
if err := connection.SendText(ctx, "Hello"); err != nil { return err }
event, err := connection.Recv(ctx)
if err != nil { return err }
var transcript forward.RealtimeTranscriptPayload
if event.Type == "transcript.final" {
	if err := event.DecodePayload(&transcript); err != nil { return err }
}
// Stop capture/playback before graceful close. Ordinary Close only releases
// transport resources; it does not confirm that history was saved.
closed, err := connection.CloseGracefully(ctx)
if err != nil { return err }
fmt.Println(closed.Outcome)
```

REST and each new WebSocket connection use the same PAT, direct SAT, or shared
`Credential`. A credential provider is called on each new connection; an
established socket cannot update its Authorization header. The current Realtime
service requires an SA Token bound through Forward's Token issuance APIs.
The platform SAT obtained by `ServiceAccountCredential` alone does not create
this binding and is currently rejected by Realtime. Use
`option.WithAuthMode("sa")` and `option.WithAuthToken(boundSAT)` for a directly
supplied bound SAT, or a custom Credential that supplies one. SDK transport
support does not grant additional Token permissions. Credentials stay in headers,
and upgrade redirects are rejected. Native `*http.Client` transports (including
TLS and proxy configuration), headers and middleware apply to the upgrade;
arbitrary custom HTTP doers are unsupported. The client's HTTP `Timeout` and
`option.WithRequestTimeout` do not limit a Realtime connection's lifetime.

`RealtimeConnectParams` defaults to 30s handshake, 60s ready, 30s send, 10s
graceful close, 256 queued events / 4 MiB, and 1 MiB per received message.
Exceeding receive limits terminates the connection with an explicit error.
`Recv(ctx)` preserves order, including buffered events before a terminal error.
Known payload types cover transcripts, audio, state, errors, close and Work
events; `Payload`, `DecodePayload` and `RawJSON()` preserve unknown extensions.
The timestamp stays text, including nanosecond precision. `voice.replaced`
ends the old connection; a nonterminal `error` or `work.failed` remains an event.

`AppendAudio` takes raw **16 kHz mono PCM16 little-endian** bytes, not WAV, and
sends Base64 in JSON text frames. Output is **24 kHz mono PCM16**;
`RealtimeEvent.Audio()` decodes `audio.delta`. The final client JSON message
must fit within 256 KiB; text is limited to 16 KiB after trimming whitespace.
Applications own sampling, playback and its lifecycle: `audio.done` means
generation finished, while `PlaybackStarted`, `PlaybackEnded` and
`PlaybackCancelled` report what the speaker actually did. Copy `work_id` and
`announcement_id` from the corresponding audio event into the receipt.
`Interrupt` does not stop an application's speaker by itself.

The service limits a connection to 60 minutes. The SDK does not automatically
reconnect or replay text/audio; applications may reconnect to the same
Conversation. Closing a connection neither cancels background Work nor deletes
the Conversation. Graceful closure succeeds only after a matching request ID
and successful persistence outcome; timeouts and persistence failures stay errors.

History has its own `conversation / events / page` envelope and `before` cursor,
and iterates towards older records. Historical `id` values are independent of
the WebSocket `event_id` and sequence. Filters are `message`, `work` or
`message,work`, and page size is 1–100. `HasNextPage` / `GetNextPage(ctx)` expose
manual paging; `HistoryAutoPaging` also crosses empty pages and detects cursor cycles.

```go
history := client.Realtime.Conversations.HistoryAutoPaging(ctx, conversation.ID,
	forward.RealtimeConversationHistoryParams{
		Types: forward.String("message,work"), Limit: forward.Int(50),
	})
for history.Next() {
	event := history.Current()
	fmt.Println(event.ID, event.Type, event.OccurredAt)
}
if err := history.Err(); err != nil { return err }
```

The standalone [Realtime example](../examples/forward/realtime/main.go) uses
process environment variables and a pre-enabled Template/Identity. It defaults
to ping/close/history without sending model input. Optional `-text` and `-pcm`
use model quota; choose either text or raw speech PCM (at most 30 seconds,
including trailing silence), sent in 20 ms chunks. The demo does not play audio and
reports received output as cancelled. New Conversations remain stored because
the API exposes no Conversation deletion endpoint.

```sh
# PAT defaults; SA uses the same QODER_AUTH_MODE/SA_KEY/ACCESS_TOKEN convention.
QODER_SA_REGION=cn go run ./examples/forward/realtime \
  -identity "$QODER_IDENTITY_ID" -template "$QODER_TEMPLATE_ID" \
  -idempotency-key "$QODER_IDEMPOTENCY_KEY"
# Supply a Forward-issued bound SAT through QODER_ACCESS_TOKEN; do not set QODER_SA_KEY.
QODER_AUTH_MODE=sa QODER_SA_REGION=global go run ./examples/forward/realtime \
  -conversation "$QODER_CONVERSATION_ID" -text "Hello"
go run ./examples/forward/realtime -conversation "$QODER_CONVERSATION_ID" -history-only
```

Maintainers can run `make test-live-realtime` for an explicit CN/Global × PAT/SA
acceptance suite. Configure the four pre-enabled
`QODER_REALTIME_{CN,GLOBAL}_{PAT,SA}_TEMPLATE_ID` values and matching
`QODER_FORWARD_{CN,GLOBAL}_PAT` / `QODER_REALTIME_{CN,GLOBAL}_SAT` credentials, then
set both `QODER_REALTIME_LIVE_ALLOW_WRITE=true` and
`QODER_REALTIME_LIVE_ALLOW_EXECUTION=true`. The suite validates all inputs
before creating test Identities, verifies text/output PCM and persisted history,
and cleans up its Sessions/Identities. It retains Conversation IDs and leaves
the supplied Templates intact. The SA cases use the supplied Forward-issued
bound SATs. Without them, they use `QODER_SA_{CN,GLOBAL}_KEY` through
`ServiceAccountCredential` to retain platform Token compatibility checks; the
current service rejects such unbound platform SATs. No Token issuance is
performed implicitly. Optionally set `QODER_REALTIME_LIVE_PCM_FILE` to
a raw 16 kHz mono PCM16 little-endian speech fixture of at most 30 seconds,
including trailing silence. It adds an audio round requiring a user transcript
and an assistant transcript/audio response. Without it, acceptance covers text
input; the standalone example also supports `-pcm`.

`examples/forward/` and `examples/managed/` hold 18 business programs, one `main.go` per scenario, each runnable directly. They read `.env.live` from the working directory — run them from the repository root, or point `-env` elsewhere. Copy [.env.live.example](../.env.live.example) and configure PAT or SA authentication. Examples default to CN; `-region cn|global` overrides the configured Qoder hostname.

```sh
# Run from the repository root. Listing models is read-only.
go run ./examples/forward/models
go run ./examples/managed/models

# Create resources, run one real conversation, then clean up.
go run ./examples/forward/session
go run ./examples/managed/session

# Write a memory, then check that a brand-new session's answer uses it.
go run ./examples/forward/memory
go run ./examples/managed/memory
```

| Scenario | Forward | Managed |
| --- | --- | --- |
| Models and a basic session | `models`, `session` | `models`, `session` |
| Multi-turn chat and text deltas | `conversation`, `streaming-deltas` | `conversation`, `streaming-deltas` |
| Personalization and business tools | `identity-config`: per-user overrides on a shared template | `custom-tools`: running a Go function and returning its result |
| Files, environment variables, Skills | `resources` | `resources` |
| Memory | `memory`: bound to Identity and Template | `memory`: attached as a Session resource |
| Scheduling and batches | `schedule`, `batch` | `deployment`, `dream` |

The programs print their steps, messages, assistant replies and cleanup results; `-output json` switches to structured logs, and `-timeout` bounds each scenario, excluding cleanup. Running a scenario creates real resources and consumes model usage.

### Authentication examples

The two read-only examples under `examples/auth/` show client construction directly:
SA Key exchange with automatic caching/renewal, and direct SAT access with renewal
managed by the caller. Each lists Forward templates and Managed agents; no model
execution or resources are created. The SA needs both API scopes/permissions and
its Managed workspace binding in the selected region.

```sh
# These standalone demos read process environment variables (not .env.live).
# Configure QODER_SA_REGION=cn|global and only the credential needed by the demo.
go run ./examples/auth/sa-key     # QODER_SA_KEY; unset QODER_ACCESS_TOKEN
go run ./examples/auth/sa-token   # QODER_ACCESS_TOKEN; unset QODER_SA_KEY
```

The existing business examples keep PAT as the default and read `.env.live`.
Set `QODER_AUTH_MODE=sa` to select SA explicitly, then configure exactly one of
`QODER_SA_KEY` and `QODER_ACCESS_TOKEN`. `QODER_SA_REGION=cn|global` selects both the
business and exchange hosts; `-region` overrides it (`international` remains an
alias for `global`). Leave the unused credential empty. Ambient PATs are ignored
in SA mode; missing SA credentials and conflicting Key/SAT inputs fail before a
request is sent. This selection is example configuration, not implicit SDK behavior.

```sh
# With the SA credential configured in .env.live, reuse any business example.
QODER_AUTH_MODE=sa QODER_SA_REGION=cn go run ./examples/forward/models
QODER_AUTH_MODE=sa QODER_SA_REGION=global go run ./examples/managed/session
```

For a centralized SAT service, implement `convention.Credential.Token(ctx)` and
use `option.WithAuthMode("sa"), option.WithCredential(myTokenProvider)`. The provider
is responsible for obtaining and renewing the SAT; business calls stay the same.

## Forward Webhooks

Manage webhook endpoints and verify signed event deliveries. See [Webhook usage and receiver examples](../docs/webhooks.md).

## Resource search

Search uses `POST /<resource>/search` with typed filters in the JSON body.
The SDK automatically includes `x-qoder-beta: search-2026-08-31`, preserving
other effective beta flags. Calling Search with no filters sends `{}`.

- Forward: Environments, Skills, Vaults, and Files.
- Managed: Agents, Sessions, Deployments, Environments, Skills, Vaults, Files,
  and Memory Stores.

Metadata filters are exact string matches combined with AND. Each resource also
supports its documented filters, such as name, creation time, Session status,
or Skill source/display title. See the official
[Forward search](https://docs.qoder.com/cloud-agents/api/forward/environments/search)
and [Managed Session search](https://docs.qoder.com/cloud-agents/api/sessions/search)
contracts. This searches QCA resources; it does not perform a web search.

Forward's `identity_id` selector is sent only in the URL query. Set it when a PAT
should search an Identity's resources; omit it with all SATs, whose scope comes
from the token. Search cursors are opaque: automatic pagination puts `next_page`
back into the JSON body's `page` while preserving filters, ownership, and request
options. Existing List methods retain their GET/query pagination behavior.

```go
// client is a forward.Client.
items := client.Environments.SearchAutoPaging(ctx, forward.EnvironmentSearchParams{
    Metadata: map[string]string{"team": "core"},
    Name: forward.String("dev"),
    IdentityID: forward.String("idn_example"), // PAT only; omit for SAT.
    Limit: forward.Int(20),
})
for items.Next() {
    fmt.Println(items.Current().ID)
}
if err := items.Err(); err != nil {
    return err
}

// managedClient is a managed.Client. Dates are typed time.Time values.
sessions, err := managedClient.Sessions.Search(ctx, managed.SessionSearchParams{
    Metadata: map[string]string{"ticket": "A-123"},
    Statuses: []string{"idle", "running"},
})
```

Each resource exposes `Search` for a page and `SearchAutoPaging` for iteration.
Use the page's `GetNextPage()` for manual pagination.

## Official contract alignment

The next release aligns these public types with the official API documentation
retrieved on 2026-10-27. It introduces source changes relative to 0.12.0:
Outcome `rubric` is an object; Managed Agent `effort` is a string; Managed Agent
and Dream model objects do not accept `speed`. Managed Dream `model` is optional
and defaults to `auto`. Discover model IDs and supported efforts/context windows
with List models. These Managed model changes do not remove Forward model features.

```go
rubric := forward.OutcomeRubricParam{Type: "text", Content: forward.String("report exists")}
model := managed.ManagedAgentsModelConfigParams{ID: "ultimate", Effort: managed.String("high")}
// Dream model can be omitted; its default is auto.
dream, err := client.Dreams.New(ctx, forward.DreamNewParams{
    Inputs: []forward.DreamInputParam{{Type: "identity_template",
        IdentityID: forward.String("idn_..."), TemplateID: forward.String("tmpl_...")}},
})
```

Mount list callers must supply query params, even if empty. The result now uses
the standard page type. Use `ListAutoPaging` to traverse every page:

```go
pager := client.Identities.MemoryStores.ListAutoPaging(ctx, "idn_...", "tmpl_...",
    forward.IdentityMemoryStoreListParams{Limit: forward.Int(20)})
for pager.Next() {
    fmt.Println(pager.Current().MemoryStoreID)
}
if err := pager.Err(); err != nil { /* handle error */ }
```

Identity tool maps now use `IdentityToolsUnionParam{OfOverrides: toolMap}`;
`OfArray: []forward.ToolParam{}` explicitly replaces all tools with an empty list.
Use `param.NullStruct[forward.IdentityToolsUnionParam]()` to restore inherited
tools, and `param.NullMap[map[string]forward.ManagedCapabilityOverrideParam]()`
to clear all managed-capability overrides. A selector value of
`param.NullStruct[forward.ManagedCapabilityOverrideParam]()` clears that selector.
For roster string entries use `MultiagentEntryParam{OfString: forward.String("tmpl_...")}`;
Advisor entries use `Type: "advisor", Model: forward.String("ultimate")`.

New Forward Dream calls require a PAT or an admin SAT; an Identity SAT is not
supported. Inputs contain exactly one primary Memory Store or Identity/Template
source, optionally plus Session sources. Pairing `update` changes an active
binding; supply at least one nonempty Identity or Template ID. Use pairing
`create` to complete a pending pairing. Pending records preserve empty
`identity_id`, `template_id`, and `paired_at` strings.

Identity `auto_memory_store=false` disables automatic default-store creation.
Independent mounts now accept `access="read_write"`; omit access for a new
read-only mount, or to preserve an existing mount's permission. Schedule
`resources` accepts file entries and `[]` clears them; Run retrieval supports
`include_initial_events=true` to return the saved input snapshot.

Managed cloud Environment `config.storage.buckets` supports up to eight OSS
mounts. When the resulting configuration has buckets, explicitly include beta
`oss-bucket-2026-09-23` (`Betas` in Go, `betas` in Python/TS). Send credentials
only when creating or replacing them; omit `credential` on an existing named
bucket to keep it. Responses expose a read view without role ARN or AccessKey
secret. `storage: {buckets: []}` clears mounts.

Managed `OAuth.Start` (Go) / `oauth.start` (Python/TS) returns `authorization_url`,
`state`, and `callback_origin`; the user completes authorization in a browser.
Client secrets go in the POST body, never in a URL. This method does not open a
browser or complete consent automatically.

See the official [Outcome contract](https://docs.qoder.com/cloud-agents/api/forward/sessions/send-events),
[Identity overrides](https://docs.qoder.com/cloud-agents/api/forward/identities/upsert-config),
[Dream inputs](https://docs.qoder.com/cloud-agents/api/forward/dreams/create),
[OSS storage](https://docs.qoder.com/cloud-agents/api/environments/storage-buckets),
and [OAuth flow](https://docs.qoder.com/cloud-agents/api/vaults/start-oauth).
