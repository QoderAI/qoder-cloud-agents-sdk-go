# Qoder Cloud Agents Go SDK

[Changelog](CHANGELOG.md) · [Releases](https://github.com/QoderAI/qoder-cloud-agents-sdk-go/releases)

The Qoder Cloud Agents Go SDK provides access to the Qoder Cloud Agents API from Go. It offers typed request parameters, context cancellation, streamed events and per-request configuration.

The SDK ships two clients that share authentication, transport, retries, error handling, parameter encoding and pagination:

| Mode | Package | Constructor |
| --- | --- | --- |
| Forward | `forward` | `forward.NewClient` |
| Managed | `managed` | `managed.NewClient` |

- [Installation](#installation)
- [Requirements](#requirements)
- [Usage](#usage)
- [PAT and Service Account authentication](#pat-and-service-account-authentication)
- [Streaming](#streaming)
- [Request fields](#request-fields)
- [Response objects](#response-objects)
- [Error handling](#error-handling)
- [Retries](#retries)
- [Timeouts](#timeouts)
- [Pagination](#pagination)
- [Request options](#request-options)
- [HTTP client customization](#http-client-customization)
- [Accessing raw response data](#accessing-raw-response-data)
- [Making undocumented requests](#making-undocumented-requests)
- [Documentation and examples](#documentation-and-examples)
- [Development](#development)
- [Versioning](#versioning)
- [License](#license)

## Installation

```sh
go get github.com/QoderAI/qoder-cloud-agents-sdk-go@v0.7.0
```

This SDK is pre-1.0, so pin a version explicitly and keep the version your build resolved in `go.mod` and `go.sum`.

## Requirements

- Go **1.23** or later.
- A PAT or Service Account credential for the target environment, with access to the resources you call. Managed SA access also needs a workspace binding.

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
| Token variable | `QODER_ACCESS_TOKEN` (PAT mode falls back to `QODER_PAT`) | `QODER_ACCESS_TOKEN` (PAT mode falls back to `QODER_PAT`) |
| Base URL variable | `QODER_FORWARD_BASE_URL` | `QODER_BASE_URL` |
| Default base URL | `https://api.qoder.com/api/v1/forward/` | `https://api.qoder.com/api/v1/cloud/` |

The defaults point at the international site; pass the CN base URL explicitly to reach CN. A base URL must include the full API root path, and the trailing `/` is optional. If the two modes use different accounts, pass each token to its own client.

The clients do not read `.env` files. Examples and explicit live-test commands have their own config loaders. Examples default to PAT; see the authentication section for their SA configuration.

For rotating credentials, `option.WithCredential(provider)` accepts any implementation of `convention.Credential`. The provider is evaluated per request. An explicit static token or Authorization header takes precedence; an explicit credential now overrides ambient PAT.

The remaining examples are function fragments meant to be dropped into an application; `package` and `import` blocks are omitted. Import the standard library packages each fragment uses, plus the SDK packages:

| Package | Import path |
| --- | --- |
| `forward` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/forward` |
| `managed` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/managed` |
| `convention` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention` |
| `option` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option` |
| `param` | `github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param` |

## PAT and Service Account authentication

PAT and SA are parallel modes. Existing `option.WithPAT(token)` remains supported.
For an existing SAT use `option.WithAuthMode("sa")` and `option.WithAuthToken(token)`.
Tokens are opaque: `WithPAT` remains an alias of `WithAuthToken`; the SDK does not
inspect a token to infer its owner. `QODER_ACCESS_TOKEN` supplies the selected mode;
only PAT mode falls back to `QODER_PAT`. Clients never implicitly read an SA Key.

```go
credential, err := convention.ServiceAccountCredentialFromEnv(
    convention.ServiceAccountCredentialConfig{
        Region: "global",
        Scopes: []string{"qca.access", "forward.access"},
    },
)
if err != nil { return err }
forwardClient := forward.NewClient(option.WithCredential(credential))
managedClient := managed.NewClient(option.WithCredential(credential))
```

Both clients share one cache. Use `Region: "cn"` with matching
`https://api.qoder.com.cn/api/v1/forward/` and `https://api.qoder.com.cn/api/v1/cloud/`
BaseURLs. Managed-only integrations can request `[]string{"qca.access"}`.
Permissions and workspace bindings are configured server-side.

The SA Key is exchanged lazily and sent only to the exchange service. Cache expiry
uses the server's actual lifetime and a monotonic clock. Concurrent callers share
results; failures enter a cooldown. A transient failure can reuse a still-safe SAT;
authentication and invalid-response failures cannot. There is no PAT fallback,
background refresh, or automatic replay of a 401 write. Existing streams remain
open. A cancelled caller stops waiting; the shared exchange has its own timeout
(default 30 seconds) and can finish/cache for other callers.

Do not combine the native SA credential with a static token. Explicit static
Token + legacy custom Credential retains static precedence. Explicit credentials
now override ambient PAT. Header overrides/deletion and custom authentication
middleware remain available. Rotate a Key by constructing a new credential.

Run the separate read-only CN/Global smoke using `make test-live-sa`, with
`QODER_SA_CN_KEY` and `QODER_SA_GLOBAL_KEY` in `.env.live` or the process environment.
Process variables take precedence; `LIVE_ENV_FILE` selects another file.
Both keys are checked before the first API call. For write/read, Session SSE and
cleanup verification, run `make test-live-sa-execution` with
`QODER_SA_LIVE_ALLOW_WRITE=true` and `QODER_SA_LIVE_ALLOW_EXECUTION=true`.
`QODER_SA_MODEL` defaults to `auto`. Offline tests and existing PAT suites remain
independent; `QODER_AUTH_MODE` switches examples, not the PAT test suites.
See the [official authentication contract](https://docs.qoder.com/cloud-agents/api/conventions/authentication).

## Streaming

Streaming responses are consumed incrementally. Close the stream when finished and handle interruption explicitly; closing a local stream does not cancel server-side work.

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

Managed response unions expose `AsAny()` and `As...()` accessors; check the discriminator before calling one. An `agent.message` event, for instance, is read through `event.AsAgentMessage()`.

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

## Documentation and examples

- [API reference](docs/api/reference.md)
- [Runnable examples](examples/)
- [Contribution guide](CONTRIBUTING.md)

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
