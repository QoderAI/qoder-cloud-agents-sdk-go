# Contributing

Thank you for contributing to the Qoder Cloud Agents Go SDK.

## Development workflow

1. Create a focused branch from the latest `main`.
2. Keep each pull request limited to one independently reviewable change.
3. Use Conventional Commits, for example `fix(forward): preserve request id`.
4. Open a pull request and wait for all required checks before merging.

GitHub `main` is the source of truth. Do not develop against or copy changes from the internal CI mirror.

## Setup and checks

Use Go 1.23 or newer.

```bash
go mod download
make lint
make build
make test
make docs-check
go build ./examples/...
```

`make test` runs the offline unit and contract suites. It must not require network access or credentials.

## Test layers

- **Unit tests** cover transport, errors, pagination, streaming, serialization, credentials, and helpers. Run `make test-unit`.
- **Contract tests** verify HTTP methods, paths, headers, request bodies, and response decoding with local transports. Run `make test-contract`.
- **Integration tests** exercise account-backed resource lifecycles. Copy `.env.live.example` to `.env.live`, use a dedicated test account, then run `make test-live-all`.
- **E2E tests** execute real models and tools. They additionally require the documented model and execution gates, then run with `make test-e2e`.
- **Examples** are runnable documentation. Business integration tests must not import helpers from `examples/`; verify examples with `go build ./examples/...`.

Never commit `.env.live`, tokens, credentials, generated logs, or test output. Live tests must register cleanup immediately after creating a resource.

Managed execution tests log the user prompt and resource IDs. If a turn fails,
they also log observed public events before session cleanup: event IDs/types,
Read/Bash inputs, text results and stop reasons. The diagnostic tail is limited
to 128 events / 256 KiB with 4 KiB per text field; omissions are reported.
Reasoning and binary blocks are excluded, and the configured PAT and URLs are
redacted. Retain the `go test -json` output (`build/test-results/e2e.json` for
`make test-e2e`) to investigate after the test deletes the session. These logs
contain test-account content and must not be committed or published.

`make test-live-all` includes Forward hourly Usage, Credential merge patches
and secret redaction, and Managed Session cancellation before and after sending
a turn. Managed deployment tests cover scoped Run listing/retrieval. Active
cancellation and scoped Run execution require the corresponding
`LIVE_ALLOW_WRITE=true` and `LIVE_ALLOW_EXECUTION=true` gates. Use a separate
`LIVE_ENV_FILE` with matching URL and PAT for CN and Global.

Usage queries the last 24 completed whole hours in Asia/Shanghai in both regions;
empty pages verify only the collection. A Session may finish before cancellation
and return HTTP 200 instead of 202; tests record which response occurred. Those
responses do not prove active cancellation occurred. Cleanup failures remain
failures; offline replay exercises assertions and cleanup without account
credentials.

## API and contract changes

When adding or changing an endpoint:

1. Update the relevant `forward/` or `managed/` implementation and public types.
2. Update the applicable fixtures under `forward/testdata/` or `managed/testdata/`.
3. Add focused unit and contract coverage, including failure behavior.
4. Regenerate API documentation with `make docs` and verify it with `make docs-check`.
5. Call out the corresponding Python and TypeScript work in the pull request, or explain why the change is language-specific.

The fixtures are maintained manually. A passing fixture test proves consistency with this repository, not automatically with service routes or the other SDKs.

Realtime protocol regressions use a local WebSocket server, never an account or
model. Run `go test -race ./forward -run Realtime -count=1` when changing the
connection, shared upgrade authentication, history or event models. The
`forward/testdata/realtime-events.json` fixture records the `voice.realtime.v1`
contract and source snapshot. Real CN/Global Realtime checks additionally need a
pre-enabled Template and an allowed Identity; the standalone example can perform
create/ready/ping/close/history, with optional text or raw PCM input. It retains
Conversation IDs because no public deletion API exists. Mock tests do not prove
that either production region has deployed the protocol.

`make test-live-realtime` separately verifies CN/Global × PAT/SA with the two
`QODER_REALTIME_LIVE_ALLOW_WRITE=true` and
`QODER_REALTIME_LIVE_ALLOW_EXECUTION=true` gates. Configure
`QODER_FORWARD_{CN,GLOBAL}_PAT`, `QODER_REALTIME_{CN,GLOBAL}_SAT` and the four
`QODER_REALTIME_{CN,GLOBAL}_{PAT,SA}_TEMPLATE_ID` values in `.env.live` or the
process environment. Realtime currently requires Forward-issued bound SATs;
the platform SAT exchanged by `ServiceAccountCredential` does not create the
Forward binding. Without an explicit regional SAT, the suite falls back to the
corresponding `QODER_SA_{CN,GLOBAL}_KEY` to retain this platform compatibility
check, which currently receives 403 from Realtime. Per region, one SA credential
source must exist. All required inputs are checked before any request.
The suite creates unique test Identities, creates/connects Conversations,
checks ready and pong, sends a short text without fixed model wording, decodes
output PCM, confirms graceful closure, and exercises filtered history at page
size 1. It finishes the associated Forward Sessions and deletes only its test
Identities; it retains/logs Conversation IDs and never edits the supplied
Templates. This entry consumes model usage. Optionally set
`QODER_REALTIME_LIVE_PCM_FILE` to a raw 16 kHz mono PCM16 little-endian speech
fixture (at most 30 seconds, including trailing silence). The file is checked
before requests; the additional audio round requires a user transcript and an
assistant transcript/audio response. Without this setting, the suite tests text
input and does not prove microphone ASR. The standalone example also supports
`-pcm`.

## Compatibility conventions

The SDK intentionally keeps Qoder-branded `X-Qoder-*` metadata headers and resumable session-event streams. Preserve those extensions unless the change explicitly revises the public contract. Breaking public API changes require a minor-version release while the SDK remains pre-1.0 and must include migration notes.

## Changelog

Keep user-facing notes in [CHANGELOG.md](CHANGELOG.md). Add a concise bullet under
`## [Unreleased]` in the same PR as a feature, fix, deprecation, or behavior change.
Use English consistently with the public documentation, identify Forward or Managed
when relevant, and describe the effect on SDK users. Pure CI changes and internal
refactors normally need no entry; explain that in the PR checklist.

Use `### Added`, `### Fixed`, `### Changed`, `### Breaking changes`, and
`### Migration` as needed. Omit empty categories. Breaking changes must explain
what callers need to change, preferably with a small migration example. Use inline
links for issues and PRs so the extracted entry also works on the Releases page.

In the release PR, move the completed `Unreleased` notes into exactly one
`## [<version>]` section immediately below it, alongside the version-file updates.
Keep an empty `Unreleased` section for subsequent work. A date is optional; when
included, use `## [<version>] - YYYY-MM-DD`. Use the exact canonical package version
without `v`, including any prerelease suffix. Each SDK keeps its own version and
notes; a shared `batch_id` can associate coordinated releases.

Validate the file and preview a prepared version locally (Python 3.10+):

```bash
python3 .github/scripts/release_notes.py check
python3 .github/scripts/release_notes.py extract --version <version>
python3 -m unittest discover -s .github/scripts -p 'test_*.py'
```

PR CI checks the format and release automation tests. Release preflight requires a
nonempty entry for the requested version, rejects placeholders such as `TODO` or
`TBD`, and displays the extracted notes in the workflow summary before approval.
Only approved-commit notes are used; the workflow never writes back to `main`.

After public package verification succeeds, the workflow creates a GitHub Release
on the existing annotated `v<version>` tag using that entry. Prereleases are marked
as such and are not promoted to Latest. A rerun checks the remote tag's commit and
reuses a published Release only when its title, notes, and prerelease flag match.
Conflicting or draft Releases fail for manual inspection rather than being
silently overwritten. If only this final step fails, rerun the failed job to finish
the Release; the package may already be publicly available.

The `0.1.0` entry documents the existing baseline. This automation applies to future
release commits that contain the changelog and scripts; it does not move old tags
or republish historical packages.

## Release

Before the first release, create the GitHub `release` Environment with required reviewers and a deployment-branch rule limited to `main`. Add a tag ruleset for `refs/tags/v*` that blocks updates and deletions and allows creation only by the release automation identity used by this workflow. Keep those protections enabled; do not dispatch the workflow until they are configured.

For each release, merge a focused pull request that updates `convention/version.go` and prepares the matching version entry in `CHANGELOG.md`, then run `make check-version VERSION=<version>` locally. From the workflow page, select the `main` ref and provide the version without `v`, the current full lowercase 40-character `main` commit SHA, and a 1-64 character `batch_id` that starts with a letter or digit and otherwise contains only letters, digits, `.`, `_`, or `-`.

The workflow revalidates `main`, runs the offline lint, build, test, documentation, and example gates, and creates the annotated module tag after Environment approval. It then waits for the public Go proxy and verifies that the exact tag resolves to the approved commit from empty module and build caches, then publishes the changelog entry as a GitHub Release. Release tags are immutable: never move, delete, or overwrite one. Fix a bad release forward with a new version bump and a new workflow run.

Before creating the tag, proxy checks use the `/cached-only` endpoint so they do
not fetch an unpublished version and cache a miss. After publishing, verification
uses the normal public proxy, with 35 minutes of retry delays and a 45-minute job
timeout. This accommodates the [proxy's 30-minute negative cache](https://proxy.golang.org/)
if another client requested the version before its tag existed. A final failure
prints the response headers and body for diagnosis. Once the proxy serves the
expected version, rerun the failed jobs with the original version, SHA and batch;
do not recreate the tag.

## Pull requests

Complete the pull request template, include exact verification commands and results, and identify public API, documentation, integration-test, and cross-SDK effects. Do not combine unrelated refactors with behavior changes.

## Service Account authentication checks

The separate read-only SA suite requires `QODER_SA_CN_KEY` and
`QODER_SA_GLOBAL_KEY` in `.env.live` or the process environment and provisioned SA permissions
(and Managed workspace bindings). Run `make test-live-sa`. It tests both API
surfaces in CN and Global using automatic exchange and static SAT. Required keys
are checked before the first API call; normal PAT suites remain independent.
The read-only command does not prove write, model execution, SSE, or
user/identity isolation coverage. Do not print
SA Keys or SATs in test reports.

For write/read, Session SSE and cleanup verification, use dedicated SA accounts
and run `make test-live-sa-execution` with
`QODER_SA_LIVE_ALLOW_WRITE=true` and `QODER_SA_LIVE_ALLOW_EXECUTION=true`.
`QODER_SA_MODEL` selects the model (default `auto`). This separate entry creates
unique test Environments, Templates/Agents, Forward Identities and Sessions in
both regions, sends one message per Session, checks the assistant reply and idle
event through SSE, and cleans up resources even after execution fails. Cleanup
failures fail the test.
Missing keys or execution flags fail this explicit command before resource creation.
The read-only SA command does not run these execution cases.

Process environment variables override `.env.live`; explicit empty values do not
fall back to the file. Use `LIVE_ENV_FILE` to select another file.
`QODER_AUTH_MODE` selects example authentication only. PAT and SA test entries
are separate; setting it to `sa` does not switch the ordinary PAT suites.

`make test` compiles the SA live suites and runs their configuration regression
test without credentials or network calls.

Business integration scenarios keep independent test helpers. Dedicated offline
regression tests for example configuration and lifecycle may import examples;
those tests do not replace SDK integration coverage.
