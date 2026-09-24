# Qoder Cloud Agents Go SDK

[Changelog](CHANGELOG.md) · [Releases](https://github.com/QoderAI/qoder-cloud-agents-sdk-go/releases)

The Qoder Cloud Agents Go SDK provides access to the Qoder Cloud Agents API from Go. It offers typed request parameters, context cancellation, streamed events and per-request configuration.

The SDK ships two clients that share authentication, transport, retries, error handling, parameter encoding and pagination:

| Mode | Package | Resources |
| --- | --- | --- |
| Forward | `forward` | Identity, Template, Session, Schedule, Batch, Channel, Environment, File, Skill, Vault, MemoryStore, Model |
| Managed | `managed` | Agent, Session, Deployment, Dream, Environment, File, Skill, Vault, MemoryStore, Model |

- [Installation](#installation)
- [Requirements](#requirements)
- [Usage](#usage)
- [Conversations](#conversations)
- [System prompts and tools](#system-prompts-and-tools)
- [Streaming](#streaming)
- [Request fields](#request-fields)
- [Response objects](#response-objects)
- [Error handling](#error-handling)
- [Retries](#retries)
- [Timeouts](#timeouts)
- [Long-running operations](#long-running-operations)
- [File uploads and downloads](#file-uploads-and-downloads)
- [Pagination](#pagination)
- [Request options](#request-options)
- [HTTP client customization](#http-client-customization)
- [Accessing raw response data](#accessing-raw-response-data)
- [Making undocumented requests](#making-undocumented-requests)
- [Examples](#examples)
- [Development](#development)
- [Versioning](#versioning)
- [License](#license)

## Installation

```sh
go get github.com/QoderAI/qoder-cloud-agents-sdk-go@v0.1.0
```

`v0.1.0` is the current release, so `@latest` resolves to it. It is pre-1.0, so pin a version explicitly and keep the version your build resolved in `go.mod` and `go.sum`.

## Requirements

- Go **1.23** or later.
- A personal access token (PAT) for the target environment, with access to the resources you call.

## Usage

The complete program below reads a PAT from the environment, connects to the CN environment and lists the models enabled for the account.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	token := os.Getenv("QODER_PAT")
	if token == "" {
		log.Fatal("set QODER_PAT")
	}
	client := forward.NewClient(
		option.WithPAT(token),
		option.WithBaseURL("https://api.qoder.com.cn/api/v1/forward"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	models, err := client.Models.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range models.Data {
		if model.IsEnabled {
			fmt.Println(model.ID)
		}
	}
}
```

Managed works the same way through `managed.NewClient` with the CN base URL `https://api.qoder.com.cn/api/v1/cloud`; its model listing is `client.Models.List(ctx, managed.ModelListParams{})`.

Explicit options always win over the environment. When a base URL or token is not passed, each client falls back to:

| Setting | Forward | Managed |
| --- | --- | --- |
| PAT variable | `QODER_PAT` | `QODER_PAT` |
| Base URL variable | `QODER_FORWARD_BASE_URL` | `QODER_BASE_URL` |
| Default base URL | `https://api.qoder.com/api/v1/forward/` | `https://api.qoder.com/api/v1/cloud/` |

The defaults point at the international site; pass the CN base URL explicitly to reach CN. A base URL must include the full API root path, and the trailing `/` is optional. If the two modes use different accounts, pass each token to its own client.

The clients do not read `.env` files. The programs under `examples/` have their own config loader, which prefers real environment variables and otherwise reads `QODER_FORWARD_PAT`, `QODER_MANAGED_PAT` and friends from `.env.live`, falling back to `QODER_PAT`.

For rotating credentials, `option.WithCredential(provider)` accepts any implementation of `convention.Credential`. The provider is evaluated per request, and an `Authorization` header that is already present — including one derived from `QODER_PAT` — takes precedence, so do not configure a static token alongside a dynamic credential.

The remaining examples are function fragments meant to be dropped into an application; `package` and `import` blocks are omitted. Import the standard library packages each fragment uses, plus the SDK packages:

| Package | Import path |
| --- | --- |
| `forward` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/forward` |
| `managed` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/managed` |
| `convention` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention` |
| `option` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option` |
| `param` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param` |

## Conversations

Forward creates a Session from an Identity and a Template. The fragments below assume both already exist, with a usable execution environment configured on the Template; [the Forward session example](examples/forward/session/main.go) shows how to create them.

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

Managed creates a Session from an existing Agent and Environment; see [the Managed session example](examples/managed/session/main.go).

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

Close the stream, and check `Err()` once the loop ends. An idle event does not always mean success: it can also mean the run is waiting on the caller or hit a budget, so judge the outcome from `stop_reason` together with the final answer. [The scenario helper](examples/internal/live/turn.go) shows a full set of execution assertions.

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

## Request fields

Required fields take plain values; optional scalars take `param.Opt[T]`. Use `forward.String`, `managed.String`, `Bool`, `Int` and `Float` to set an optional value explicitly.

| Go value | Meaning on the wire |
| --- | --- |
| unset `param.Opt[T]` | field is omitted |
| `forward.String("")` | sends an empty string |
| `forward.Bool(false)` | sends `false` |
| `forward.Int(0)` | sends the number `0` |
| `param.Null[string]()` | sends JSON `null` |
| nil map or slice (with `omitzero`) | field is omitted |
| non-nil empty map or slice | sends `{}` or `[]` |

`param.IsOmitted(value)` and `param.IsNull(value)` report the state of a parameter. For an explicit null on a struct, map or slice, use `param.NullStruct[T]()`, `param.NullMap[T]()` or `param.NullSlice[T]()`. Whether a given field may be cleared is still decided by that API's validation.

The fragment below only serializes; it sends nothing:

```go
func encodeOptionalFields() ([]byte, error) {
	params := forward.IdentityUpdateParams{
		Enabled: forward.Bool(false),
	}
	// Name is unset and will be omitted; Enabled is sent explicitly as false.
	return json.Marshal(params)
}
```

### Request unions

A union parameter picks one representation through an `Of...` field, and exactly one branch should be set. A Forward model parameter, for example, is either a model ID or a configuration object:

```go
func configuredModel(modelID string) forward.ModelConfigUnionParam {
	return forward.ModelConfigUnionParam{
		OfConfig: &forward.ModelConfigParam{
			ID:     modelID,
			Effort: forward.String("high"),
		},
	}
}
```

Options such as `Effort` must match what the selected model supports. To send only the ID, use `forward.ModelConfigUnionParam{OfString: forward.String(modelID)}`.

### Extra fields

`SetExtraFields` adds fields that the server already accepts but the SDK does not model yet. An extra field overrides the struct field of the same name, so restrict it to names and values your application controls.

```go
func removeSessionMetadataKey() forward.SessionUpdateParams {
	params := forward.SessionUpdateParams{}
	params.SetExtraFields(map[string]any{
		"metadata": map[string]any{"obsolete": nil},
	})
	return params
}
```

### Deserializing params

`param.SetJSON` sends a request body exactly as given. It controls serialization of the whole object, so fields already set on the struct are not merged in. The server still validates the result.

```go
func paramsFromJSON(raw []byte) (managed.AgentNewParams, error) {
	var params managed.AgentNewParams
	if !json.Valid(raw) {
		return params, fmt.Errorf("invalid JSON")
	}
	param.SetJSON(raw, &params)
	return params, nil
}
```

Do not rely on unmarshalling arbitrary JSON into a param union and getting every `Of...` branch back. Use the response types to read data, and `SetJSON` to preserve a request verbatim.

## Response objects

Response fields are read directly. The `JSON` metadata distinguishes an absent field from `null`, from a type mismatch, from a real value:

```go
func inspectIdentity(identity forward.Identity) {
	fmt.Println(identity.ID)
	if identity.JSON.Name.Valid() {
		fmt.Println("display name:", identity.Name)
	}
	switch identity.JSON.Name.Raw() {
	case "":
		fmt.Println("the response had no name")
	case "null":
		fmt.Println("name was null")
	}
}
```

`RawJSON()` returns the JSON an object was decoded from, and unknown fields remain available through `result.JSON.ExtraFields["field_name"].Raw()`. When `Valid()` is false, read `Raw()` to tell missing from null from a mismatched type.

### Response unions

Managed response unions expose `AsAny()` and `As...()` accessors; check the discriminator before calling one. An `agent.message` event, for instance, is read through `event.AsAgentMessage()` as in the Managed stream above.

Some open-ended Forward content stays as `json.RawMessage` — `SessionEvent.Content` is one. For an array of text content blocks call `ContentBlocks()`; otherwise decode per event type, or keep `RawJSON()`. Both modes preserve unknown fields, but the accessor depends on the response type.

## Error handling

When the API returns an error status, the SDK returns a `*convention.Error`; `forward.Error` and `managed.Error` are aliases of that type. Use `errors.As` to read the status, code and request ID:

```go
func lookupSession(ctx context.Context, client forward.Client, sessionID string) error {
	_, err := client.Sessions.Get(ctx, sessionID)
	if err == nil {
		return nil
	}

	var apiErr *convention.Error
	switch {
	case errors.Is(err, context.Canceled):
		fmt.Println("call was cancelled")
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("timed out waiting for the request")
	case errors.As(err, &apiErr):
		fmt.Printf("status=%d code=%s type=%s request_id=%s\n",
			apiErr.StatusCode, apiErr.Code, apiErr.Type(), apiErr.RequestID)
	default:
		fmt.Printf("network, parameter or decoding error: %v\n", err)
	}
	return err
}
```

| Field or method | Contents |
| --- | --- |
| `StatusCode` | HTTP status code |
| `Code` | server error code, when the error response carries one |
| `Message` | server error description |
| `Type()` | error type |
| `RequestID` | request identifier, empty when the server sends none |
| `Request` / `Response` | the underlying HTTP request and response |
| `RawJSON()` | the error body as received, including non-JSON gateway output |

Network, context, parameter-construction and decoding failures are not guaranteed to be `*convention.Error`, so keep the other branches. When investigating, record the status, code and request ID — but note that the raw request may carry the PAT and the body may carry business data, so neither belongs in a general-purpose log.

For streams, connection and read failures surface through `stream.Err()`, while `session.error`, terminal statuses and unexpected stop reasons are business events the caller has to handle.

## Retries

Some errors are retried automatically, up to **2 times** by default — three HTTP attempts in total. A retry requires a replayable request body and a caller context that is still live.

- `GET` and `HEAD`, plus any request carrying an `Idempotency-Key`: retried on connection errors, 408, 429 and 5xx.
- Any other request without an idempotency key: retried on 429 only.
- **409 is never retried automatically.**
- Within those rules, a server `x-should-retry` header can force a retry on or off.

The wait comes from a valid `Retry-After-Ms` or `Retry-After` header when present, and from exponential backoff otherwise. Retries do not extend the caller's context deadline.

```go
func clientWithoutRetries(token string) forward.Client {
	return forward.NewClient(
		option.WithPAT(token),
		option.WithBaseURL("https://api.qoder.com.cn/api/v1/forward"),
		option.WithMaxRetries(0),
	)
}

func getWithRetries(ctx context.Context, client forward.Client, sessionID string) (*forward.Session, error) {
	return client.Sessions.Get(ctx, sessionID, option.WithMaxRetries(3))
}
```

An idempotency key stands for one logical operation, and the caller generates and stores it. The SDK does not mint keys for write operations, and a retry must not switch to a new key.

## Timeouts

By default, Forward and Managed clients clone `http.DefaultTransport` and set a 10-minute response-header timeout, matching Anthropic's Go SDK. This bounds the wait for response headers after the request has been fully written; it does not limit reading the response body or the total duration of an SSE stream. A custom `WithHTTPClient` takes precedence; if `http.DefaultTransport` has been replaced by a wrapper, that wrapper is preserved and controls its own timeouts.

There is no default whole-call deadline. Bound the whole call and each HTTP attempt separately:

```go
func getWithTimeout(parent context.Context, client forward.Client, sessionID string) (*forward.Session, error) {
	ctx, cancel := context.WithTimeout(parent, time.Minute)
	defer cancel()
	return client.Sessions.Get(ctx, sessionID,
		option.WithRequestTimeout(15*time.Second),
	)
}
```

Here the context's minute covers every attempt and the backoff between them, while each attempt gets at most 15 seconds; whichever deadline arrives first applies. `WithRequestTimeout` also advertises the per-attempt deadline to the server as `X-Qoder-Timeout`, in whole seconds. Whether a timed-out attempt is retried still depends on the method, idempotency key and body rules above.

For streams and downloads the timeout also covers reading the response body. Do not reuse a short request timeout for an event stream that is meant to stay open: override `WithRequestTimeout` on the streaming call, and give the whole subscription an appropriate context.

Cancelling the local context stops the HTTP call or the subscription. If cloud execution has already started, cancelling or interrupting it is a separate API call, after which the server state is worth confirming.

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

Uploading a file and attaching it to a Session are two steps; [the Forward](examples/forward/resources/main.go) and [Managed](examples/managed/resources/main.go) resource examples show the attachment. Skill file trees go through `Skills.New` and `Skills.Versions.New`, where each reader in `files` is named with its relative path, such as `my-skill/SKILL.md`.

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

## Pagination

`ListAutoPaging` fetches subsequent pages on demand; iterate with `Next()`, `Current()` and `Err()`:

```go
func listTemplates(ctx context.Context, client forward.Client) error {
	pager := client.Templates.ListAutoPaging(ctx, forward.TemplateListParams{
		Limit: forward.Int(20), Status: forward.String("active"),
	})
	for pager.Next() {
		fmt.Println(pager.Current().ID)
	}
	return pager.Err()
}
```

Managed lists iterate identically, for example `client.Agents.ListAutoPaging(ctx, managed.AgentListParams{Limit: managed.Int(20)})`.

`List` returns a single page instead, which `GetNextPage()` advances manually:

```go
func listAgentPages(ctx context.Context, client managed.Client) error {
	page, err := client.Agents.List(ctx, managed.AgentListParams{Limit: managed.Int(20)})
	for err == nil && page != nil {
		for _, agent := range page.Data {
			fmt.Println(agent.ID)
		}
		page, err = page.GetNextPage()
	}
	return err
}
```

Different APIs page by `after_id` / `before_id` or by `next_page`. The pagers carry the parameters the next request needs and preserve the filters, so use the pagination methods of the response you have rather than mixing cursor styles.

## Request options

The `option` package configures a client constructor or a single call. A method-level option overrides the same setting from the client or service, while options that accumulate — middleware, added headers — combine by their own rules. Finish configuring a client before sharing it, so that its Options are not mutated concurrently.

| Option | Purpose |
| --- | --- |
| `WithBaseURL` | full API root URL |
| `WithPAT` / `WithCredential` | static token or dynamic credential |
| `WithMaxRetries` / `WithRequestTimeout` | retry budget and per-attempt timeout |
| `WithHeader` / `WithHeaderAdd` / `WithHeaderDel` | set, add or remove a request header |
| `WithQuery` / `WithQueryAdd` / `WithQueryDel` | custom query parameters |
| `WithJSONSet` / `WithJSONDel` | edit fields in the JSON request body |
| `WithResponseInto` | capture the raw HTTP response and its headers |
| `WithResponseBodyInto` | replace the default decoding target |
| `WithRequestBody` | supply a pre-serialized body |
| `WithHTTPClient` / `WithMiddleware` | transport and per-attempt interception |

The full list lives in [requestoption.go](convention/option/requestoption.go).

## HTTP client customization

`WithHTTPClient` configures connection pooling, proxies or a custom transport. The fragment below keeps the standard transport settings, sets the same response-header timeout as the SDK default, and raises the idle connection limit:

```go
func clientWithTransport(token string) managed.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = 10 * time.Minute
	transport.MaxIdleConnsPerHost = 20
	return managed.NewClient(
		option.WithPAT(token),
		option.WithBaseURL("https://api.qoder.com.cn/api/v1/cloud"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
	)
}
```

`WithMiddleware` intercepts API requests, for example to record timings. Middleware runs on every HTTP attempt:

```go
func requestTiming(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
	started := time.Now()
	response, err := next(req)
	status := 0
	if response != nil {
		status = response.StatusCode
	}
	log.Printf("method=%s path=%s status=%d duration=%s",
		req.Method, req.URL.Path, status, time.Since(started))
	return response, err
}
```

Pass `option.WithMiddleware(requestTiming)` to `NewClient`. API middleware does not apply to the temporary-URL storage download; configure that through the HTTP client or transport instead.

## Accessing raw response data

`WithResponseInto` captures the `*http.Response` alongside the decoded value, which is how response headers such as the request ID become visible:

```go
func inspectResponse(ctx context.Context, client forward.Client, sessionID, traceID string) error {
	var response *http.Response
	_, err := client.Sessions.Get(ctx, sessionID,
		option.WithHeader("X-Trace-ID", traceID),
		option.WithResponseInto(&response),
	)
	if response != nil {
		fmt.Println("Request ID:", response.Header.Get("X-Request-ID"))
	}
	return err
}
```

The captured response's body is not guaranteed to be unread, so prefer the entity's `RawJSON()` when you want the JSON itself. To decode into something else entirely, pass `option.WithResponseBodyInto(&dst)`.

## Making undocumented requests

This library is typed for the documented API surface. Endpoints, parameters and fields outside it are still reachable.

### Undocumented endpoints

`convention.ExecuteNewRequest` issues an arbitrary request against the configured base URL. Passing `client.Options...` reuses the client's authentication, retries and middleware:

```go
func callUndocumentedEndpoint(ctx context.Context, client forward.Client, sessionID string) (map[string]any, error) {
	var response map[string]any
	err := convention.ExecuteNewRequest(ctx, http.MethodPost,
		"sessions/"+sessionID+"/unreleased-action",
		map[string]any{"flag": true},
		&response,
		client.Options...,
	)
	return response, err
}
```

The path is relative to the base URL. No parameter validation or response typing applies here, so treat the result as untyped JSON.

### Undocumented request params

`option.WithJSONSet` adds or overwrites a field in a typed request body, and `option.WithQueryAdd` does the same for the query string:

```go
func createSessionWithUnknownField(ctx context.Context, client forward.Client, identityID, templateID string) (*forward.Session, error) {
	return client.Sessions.New(ctx, forward.SessionNewParams{
		IdentityID: identityID,
		TemplateID: templateID,
	}, option.WithJSONSet("unreleased_flag", true))
}
```

For a field that belongs to the request model rather than to one call, `SetExtraFields` is the better fit — see [Extra fields](#extra-fields).

### Undocumented response properties

Fields the SDK does not model are preserved during decoding and read back through the `JSON` metadata:

```go
func readUnknownField(session *forward.Session) string {
	return session.JSON.ExtraFields["unreleased_field"].Raw()
}
```

`RawJSON()` on the entity returns the whole body it was decoded from.

## Examples

`examples/` holds 18 standalone programs under `examples/forward/` and `examples/managed/`, one `main.go` per scenario, each runnable directly. Each program reads `.env.live` from the working directory — so run them from the repository root, or point `-env` elsewhere. Copy [.env.live.example](.env.live.example) and fill in the URL and PAT for each mode. Unlike the clients, the examples default to the CN hostname, and `-region cn|international` overrides whichever host the configuration names.

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

## Development

```sh
make test
make build
make lint
go test -race ./...
```

The offline tests cover request serialization, API contracts, response decoding, errors, retries, pagination, SSE and cleanup behaviour. Files ending in `_live_test.go` need the `live` build tag and a test configuration; the runnable examples do not.

### API reference

The committed API reference lives at [`docs/api/reference.md`](docs/api/reference.md) and is regenerated from the current source with a pinned `gomarkdoc` release:

```sh
make docs        # regenerate docs/api/reference.md
make docs-check  # regenerate + drift/normalization/core-surface/link gate (used in CI)
```

The gate is what CI enforces; run `make docs` and commit the result whenever you change GoDoc on an exported symbol. Under the hood `make docs` runs `go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@v1.1.0 --repository.url … --repository.default-branch main --repository.path / ./forward ./managed ./convention/...` — the version is pinned in `internal/docs/generate.go`.

## Versioning

The module is pre-1.0. Per semantic versioning, the compatibility guarantee does not apply below `v1.0.0`, so a minor release may change the API; pin a version in `go.mod` and read the release notes before upgrading.

## License

Released under the [MIT License](LICENSE).
