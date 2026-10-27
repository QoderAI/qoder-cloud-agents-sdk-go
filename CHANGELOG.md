# Changelog

User-facing changes to this SDK are recorded here. History starts at the
existing `0.1.0` release; earlier development prereleases are not listed.

## [Unreleased]

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
