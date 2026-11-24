# Changelog

User-facing changes to this SDK are recorded here. History starts at the
existing `0.1.0` release; earlier development prereleases are not listed.

## [Unreleased]

### Added

- Add nullable Forward Session multiagent updates and response snapshots, plus create-only `auto_title`. Title conflicts remain server-validated; Session creation does not accept multiagent overrides.
- Add independent Template A2A configuration and `Templates.GetAgentCard`, with Identity query scope, A2A version headers and unknown-preserving 0.3/1.0 Agent Cards.
- Add typed Forward Environment packages, setup/startup scripts and OSS bucket storage, with write credentials separated from masked response views. Credential omission preserves stored credentials on update.
- Add Managed `always_deny` policies across built-in tools, default configurations and MCP; canonical-name ImageGen, ImageSearch, VideoGen and DeliverArtifacts configurations; and type-only `ManagedBrowserToolset20260714` across Agent and Session/Thread tool unions.

### Fixed

- Replace the per-stream watcher goroutine in `convention/resumablestream` with `context.AfterFunc` so unclosed streams using long-lived contexts do not leak goroutines, and treat permanent TLS/certificate and scheme-mismatch `*url.Error` failures as non-retryable.
- Decode Managed tool configurations using case-insensitive canonical names and snake_case aliases, preserving the authored name when round-tripping request parameters.
- Preserve unknown fields, malformed members, explicit nulls and authored empty objects when decoding typed Forward Environment and Template A2A parameters. Typed edits replace automatically retained raw values; explicit caller overrides retain precedence. Invalid values reach server validation rather than disappearing from a configuration update.
- Clear stale null metadata when reusing a decoded Forward Session multiagent update parameter, so a subsequent roster is sent instead of an unintended clearing request.
- Expose the top-level Managed MCP toolset `permission_policy` for wildcard governance, including `always_deny`. Legacy `default_config.permission_policy` remains unchanged on the wire and does not set the wildcard policy.
- Send no request body from Forward `Dreams.Cancel` and `Dreams.Archive`. Earlier releases sent `{}`, which the service rejects with HTTP 400 `invalid_request_body`.
- Omit an empty Forward `ModelConfigParam.ID`, so an Identity Config update can change `effort`, `speed` or `context_window` while inheriting the stored model ID. Other endpoints still require the ID.
- Decode Managed Agent multiagent roster child versions saved as `"latest"`, and accept `"latest"` or positive integer strings in roster requests. Session and Deployment Agent references keep numeric versions.

### Breaking changes

- Forward `EnvironmentNewParams.Config` and `EnvironmentUpdateParams.Config` now use `EnvironmentConfigParams` instead of `map[string]any`.
- `ManagedAgentsMultiagentAgentUnion.Version` is now `ManagedAgentsMultiagentAgentVersionUnion`, and `AsAgent()` returns `ManagedAgentsMultiagentAgentReference`. Roster entries use `ManagedAgentsMultiagentAgentParams`, whose `Version` is a `ManagedAgentsMultiagentAgentVersionUnionParam`.

### Migration

- Replace `Config: map[string]any{"type": "cloud"}` with `Config: forward.EnvironmentConfigParams{Type: forward.String("cloud")}`. Use `EnvironmentConfigPackagesParam`, `EnvironmentStorageParam` and `EnvironmentStorageBucketParam` for nested configuration. Omit the entire Config on update to preserve it. A supplied Config replaces it, so resend buckets to retain mounts. Set `EnvironmentStorageCredentialParam` only when creating or replacing credentials; omit it on a bucket with the same Name to retain that bucket's credentials. Request OSS access with `option.WithHeader("X-Qoder-Beta", "oss-bucket-2026-09-23")`; Forward permits at most seven buckets.
- Clear a Session multiagent override with `param.NullStruct[forward.MultiagentConfigParam]()`, not a create-time Config field. Omit Template A2A (or send `{}`) to preserve it, use `param.NullStruct[forward.TemplateA2AConfigParam]()` to disable and clear it, and `Skills: []forward.TemplateA2ASkillParam{}` to replace the advertised skills with an empty list. A2A member nulls are rejected by the service.
- Read a saved numeric version from `Version.OfInt` and `"latest"` from `Version.OfString`; check `Version.JSON.OfInt.Valid()` first.
- Replace roster `&managed.ManagedAgentsAgentParams{ID: id, Type: "agent", Version: managed.Int(3)}` with `&managed.ManagedAgentsMultiagentAgentParams{ID: id, Type: "agent", Version: managed.ManagedAgentsMultiagentAgentVersionUnionParam{OfInt: managed.Int(3)}}`; use `OfString: managed.String("latest")` to follow the latest version.

## [0.15.0]

### Added

- Bind a Managed Session to a BYOC MCP tunnel with `tunnel_id` at creation. On Forward, set the tunnel on a Template for its sessions to inherit. Managed Deployments also accept and return the binding. Template and Deployment updates can replace or clear it.

- Add typed Forward Session budgets across create/update responses and events, including omission, explicit null clearing and decimal-string limits.
- Add Forward Session `config.model` overrides and model `speed` configuration, plus optional approval evaluation details, Outcome/status event fields and cumulative top-level Credits on `session.usage`.
- Add `startup_script` to Forward/Managed Environment models and Managed cloud/self-hosted parameters. It runs on sandbox start/resume; the external worker controls self-hosted execution.
- Add Managed Memory `version`, version `redacted`, and Vault creation-response `credentials` fields. Vault creation still does not accept credentials in its request.

### Fixed

- Decode Forward event `evaluated_permission` as a string and align Managed Thread/coordinator nested model effort and context-window fields.
- Add the current Memory version `updated` operation and list-item metadata/version fields; retain the legacy `Modified` constant as deprecated.

### Breaking changes

- Forward `SessionEvent.EvaluatedPermission` and the Managed Thread/coordinator union model `Effort` fields now use strings instead of object types.

### Migration

- Compare `event.EvaluatedPermission` directly to `"allow"`, `"ask"` or `"deny"`; automatic approval details are under `event.Evaluation.EvaluatedPermission`.
- Read nested model `Effort` directly as a string. Legacy Managed model `Speed` and Thread `Usage` remain available but deprecated; current Thread responses omit usage.
- Use `param.NullStruct[forward.SessionBudgetParam]()` to clear a budget and `param.NullStruct[forward.ModelConfigUnionParam]()` to clear an update's model override. Omit the model on create to inherit.
- Treat `session.usage` Credits as cumulative snapshots: replace the previous values rather than adding successive events.

### Changed

- Remove resource-specific sections and inventories from the README to keep it focused on SDK setup and shared client behavior.

## [0.14.0]

### Added

- Add automatic Thread SSE resumption for Forward and Managed, preserving the latest event cursor and request options across retryable disconnects; thread termination stops the subscription while idle events keep it open.
- Add a typed `Last-Event-ID` parameter to Managed Thread streams.

## [0.13.0]

### Added

- Automatically mark direct Forward and Managed Session creation with `__qca_sdk_created`, `__qca_sdk_language`, and the actual `__qca_sdk_version`, preserving caller metadata.
- Add Forward Dreams (create, list, retrieve, cancel, archive), channel pairing list/retrieve/update, and Managed OAuth start with typed requests and responses.
- Add Identity automatic memory-store control, mount access and cursor pagination, sparse managed-capability overrides, and full tools-array replacement.
- Support string Template references and Advisor model entries in Forward Multiagent rosters, Schedule file resources, and optional Run input snapshots.
- Add typed Managed OSS bucket storage with separate write credentials/read views, and Memory Store ID cursor parameters.

### Fixed

- Preserve Identity `tools` arrays and explicit null when loading JSON into request parameters, so empty-array replacement and clearing overrides retain their intended semantics.
- Align Forward Outcome rubric objects and Managed Agent/Dream model configuration with the official API. Dream model defaults to `auto` when omitted; object configuration supports `effort` and `context_window`.

### Breaking changes

- Session creation rejects conflicting values for the three reserved SDK metadata keys before sending; matching string values remain accepted. Missing or null metadata now contains the markers, which count toward API metadata limits.
- Forward event `rubric` is a text/content or file/file_id object instead of a string.
- Managed Agent `effort` is a level string instead of an object; Agent and Dream model configuration no longer declares the unsupported `speed` field.
- Identity mount `List` now takes `IdentityMemoryStoreListParams` and returns `pagination.Page[MemoryStoreMount]`; supply an empty params struct to preserve an unfiltered call.
- Wrap Identity tool maps in `IdentityToolsUnionParam{OfOverrides: ...}`; use `OfArray` for replacement. Read maps through `Tools.OfOverrides`.

### Migration

- Rename business fields that collide with the three `__qca_sdk_` keys and leave space for automatic markers. See “SDK-created Session metadata” in the README.
- See the README's “Official contract alignment” section for updated model/rubric types, pagination calls and explicit null clearing.

## [0.12.0]

### Added

- Forward resource search for Environments, Skills, Vaults, and Files, with PAT Identity scope in the query.
- Managed resource search for Agents, Sessions, Deployments, Environments, Skills, Vaults, Files, and Memory Stores, including typed metadata and resource filters.
- Automatic Search beta headers and POST body cursor pagination that preserves filters and request options across pages.

## [0.11.0] - 2026-10-27

### Changed

- Clarify that the Forward Decision API is in Beta and its behavior may change in future updates.

### Breaking changes

- Align Beta Forward Decision responses with the updated API: remove `usage`, the `refusal` answer type/field, and internal context fields (`policy_version`, `token_estimator`, `estimated_input_tokens`, `max_input_tokens`, `snapshot_at`, `scan_limited`, `truncated`, `truncation_reasons`).

## [0.10.0]

### Added

- Add Forward Decisions for synchronous evaluation of explicit state or accessible Session history, with choice and score questions, idempotent creation, retrieval, cursor pagination, and typed context truncation/token-estimate metadata.

## [0.9.0] - 2026-10-26

### Added

- Add Forward Webhook endpoint creation, listing, retrieval, partial update, deletion, enable/disable and synthetic test delivery.
- Add raw-body webhook signature verification with HMAC-SHA256, multiple v1 signatures and a 300-second timestamp window; include receiver examples and offline contracts.

## [0.8.0]

### Fixed

- Align existing Managed Session and Deployment resource types with generic HTTP(S) `git_repository`, including write-only credentials, branch/commit checkout, and optional GitHub authorization for public repositories.
- Add Forward Template `max_tool_rounds`, Batch `ignore_idle_window` / `queue_reason`, and Session list `search`; preserve omitted, null, false, and zero values according to each API contract.
- Expose Forward Session outcome evaluations and Thread statistics/Credits, Schedule and Run source/input-version metadata, and trigger time boundaries.
- Support Forward MemoryStore creation with a native `mount` (Identity, Template, and optional access), preserving the create request's idempotency key for retries.
- Expose self-hosted Environment setup scripts, Forward Environment binding metadata, Managed MemoryStore status/count/size, and Dream output file paths.
- Decode flat `agent.artifact_delivered` fields and diagnostic `session.error.error.qoder_error_code` in Managed event history, Session SSE, and Thread SSE. No additional API endpoints are introduced.

### Migration

- Omitting `max_tool_rounds` on update preserves the configured limit; null clears it. `ignore_idle_window` accepts a boolean, with false as the server default.
- Treat an absent Schedule Run input version or Thread usage as unknown, not zero. Use `retry_status` rather than `qoder_error_code` to decide whether execution can retry.

## [0.7.0]

### Fixed

- Align Managed Session budgets with the official decimal-string `max_credit_cost` contract, shared by Deployment budgets. The limit covers cumulative model and sandbox Credits; reaching it rejects the next user message without interrupting the current turn.
- Expose Managed and Forward Session `model_credits`, `sandbox_runtime_credits`, and `total_credits`; decode top-level `session.usage` Credits and per-call `model_usage.credits` / compaction `usage.credits`. Session usage may be absent. Legacy response token/list-cost fields remain available.
- Route Identity ownership through query parameters on existing Forward Environment, Skill/Version, Vault/Credential, and File operations, including multipart uploads and content downloads.
- Allow existing channel pairing creation with either a code or `channel_pairing_id`, and expose pairing scope details. No additional API endpoints are introduced.

### Migration

- **Breaking budget field change:** Replace `MaxListCost: managed.MonetaryAmountParam{...}` with `MaxCreditCost: "100.00"` in `ManagedAgentsBudgetLimitParam` (keep `Type: "limit"`). Amounts are positive decimal strings, not currency objects. Omitting a budget update preserves the limit; explicit null removes it.
- Consume the three top-level Credits fields from `session.usage` events. Replace cumulative snapshots rather than adding them; compaction Credits are already included in Session totals.
- Set `IdentityID: forward.String("idn_...")` in resource params; for methods without params, pass `forward.WithIdentityID("idn_...")`. Use this only with a PAT. With a SAT, omit `identity_id`; the server infers ownership from the token and rejects explicit values (including an empty string).
- Channel pairing requests must provide exactly one of `code` or `channel_pairing_id`; existing code-based calls remain supported.

## [0.6.0]

### Added

- Forward Drive Beta: typed directory listing, upload/download signed URLs, entry deletion and Drive clearing, with required Identity scope, independent cursor pagination, service ordering and canonical directory query encoding.
- Template `managed_tool_config.enabled_tools` in create/update requests and responses, including omission, empty-list and null semantics. Runnable Drive examples keep storage requests separate from QCA authentication; opt-in CN/Global lifecycle tests cover PAT, admin SAT and Identity SAT.

## [0.5.0]

### Added

- Forward admin and Identity SAT issuance and Forward SAT revocation, using explicit per-call SA Key authentication, typed issuance/revocation metadata, and no automatic retries. Existing platform credentials and PAT/SAT business authentication remain unchanged.
- Forward Realtime Beta: create Conversations with a stable idempotency key, connect over WebSocket using PAT/SAT or shared Credentials, exchange text and PCM16 audio, observe typed events and Work progress, report playback, interrupt, and confirm graceful closure.
- Realtime Conversation history with independent `before` pagination, persisted history IDs and preserved unknown event fields. Connections have separate stage timeouts, bounded receive queues, and explicit protocol failures; they do not follow redirects, reconnect, or replay input automatically.
- Realtime SA access requires a Forward-issued bound SAT and a Realtime-enabled Template; a platform SAT from `ServiceAccountCredential` alone does not establish the required Forward binding.

## [0.4.0] - 2026-10-12

### Added

- Read-only SA Key and SAT authentication examples, plus explicit PAT/SA selection for existing business examples without duplicating scenarios.
- Shared PAT and Service Account authentication for Forward and Managed, including explicit SA mode, direct SAT access, and lazy SA Key exchange with an instance-owned cache, early renewal, concurrent request sharing, and failure cooldown.
- Separate CN/Global SA authentication smoke checks for both API surfaces, covering automatic exchange and direct SAT use, plus opt-in write/read, Session SSE and cleanup verification.

### Changed

- Support `QODER_ACCESS_TOKEN` for the selected mode. SA never falls back to `QODER_PAT`; ordinary clients never implicitly read `QODER_SA_KEY`.
- Resolve authentication in the common layer; explicit credentials now take precedence over ambient PAT. Explicit static token plus a legacy custom credential retains static precedence.

### Migration

- For an existing SAT, select SA mode and supply the neutral token option. For automatic exchange, construct a Service Account credential with explicit region and scopes; do not combine it with a static token. Existing PAT calls require no new mode argument.
- Applications that configured both an ambient PAT and an explicit credential now use the credential. Remove the credential and configure an explicit PAT to retain personal-user authentication.

### Fixed

- Resource examples still require actual tool execution when reading files and environment variables; the later Skill answer validates its secret marker without requiring the model to reread a Skill already in context.

## [0.3.0]

### Added

- Forward Vault Credential updates for secret rotation and metadata merge patches, with automatic retries disabled for write-only updates.
- Forward Usage aggregation by Identity and Template using hourly `start_at` / `end_at` windows in Asia/Shanghai, with fractional `active_seconds` and multi-ID filters. Legacy timestamp parameters are not exposed.
- Managed Session cancellation with the lightweight acknowledgement for both active and idle sessions, plus deployment-scoped Run listing and retrieval.

## [0.2.0] - 2026-09-24

### Changed

- Forward and Managed clients now wait at most 10 minutes for response headers by default, matching Anthropic's Go SDK. This timeout does not limit response-body reads or SSE duration; explicit HTTP clients and request/context deadlines retain their own settings.

## [0.1.0]

### Added

- Forward and Managed: typed Go clients with shared authentication, retries, pagination, and streaming.
- Support context cancellation and per-request configuration on Go 1.23 and newer.
