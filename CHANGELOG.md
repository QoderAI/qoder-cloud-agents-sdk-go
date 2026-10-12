# Changelog

User-facing changes to this SDK are recorded here. History starts at the
existing `0.1.0` release; earlier development prereleases are not listed.

## [Unreleased]

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
