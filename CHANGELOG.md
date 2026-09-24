# Changelog

User-facing changes to this SDK are recorded here. History starts at the
existing `0.1.0` release; earlier development prereleases are not listed.

## [Unreleased]

## [0.2.0] - 2026-09-24

### Changed

- Forward and Managed clients now wait at most 10 minutes for response headers by default, matching Anthropic's Go SDK. This timeout does not limit response-body reads or SSE duration; explicit HTTP clients and request/context deadlines retain their own settings.

## [0.1.0]

### Added

- Forward and Managed: typed Go clients with shared authentication, retries, pagination, and streaming.
- Support context cancellation and per-request configuration on Go 1.23 and newer.
