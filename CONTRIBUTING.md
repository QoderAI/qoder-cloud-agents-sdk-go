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
- **Examples** are runnable documentation. Tests must not import helpers from `examples/`; verify examples with `go build ./examples/...`.

Never commit `.env.live`, tokens, credentials, generated logs, or test output. Live tests must register cleanup immediately after creating a resource.

## API and contract changes

When adding or changing an endpoint:

1. Update the relevant `forward/` or `managed/` implementation and public types.
2. Update the applicable fixtures under `forward/testdata/` or `managed/testdata/`.
3. Add focused unit and contract coverage, including failure behavior.
4. Regenerate API documentation with `make docs` and verify it with `make docs-check`.
5. Call out the corresponding Python and TypeScript work in the pull request, or explain why the change is language-specific.

The fixtures are maintained manually. A passing fixture test proves consistency with this repository, not automatically with service routes or the other SDKs.

## Compatibility conventions

The SDK intentionally keeps Qoder-branded `X-Qoder-*` metadata headers and resumable session-event streams. Preserve those extensions unless the change explicitly revises the public contract. Breaking public API changes require a minor-version release while the SDK remains pre-1.0 and must include migration notes.

## Release

Before the first release, create the GitHub `release` Environment with required reviewers and a deployment-branch rule limited to `main`. Add a tag ruleset for `refs/tags/v*` that blocks updates and deletions and allows creation only by the release automation identity used by this workflow. Keep those protections enabled; do not dispatch the workflow until they are configured.

For each release, merge a focused pull request that updates `convention/version.go` and any release notes, then run `make check-version VERSION=<version>` locally. From the workflow page, select the `main` ref and provide the version without `v`, the current full lowercase 40-character `main` commit SHA, and a 1-64 character `batch_id` that starts with a letter or digit and otherwise contains only letters, digits, `.`, `_`, or `-`.

The workflow revalidates `main`, runs the offline lint, build, test, documentation, and example gates, and creates only the annotated module tag after Environment approval. It then waits for the public Go proxy and verifies that the exact tag resolves to the approved commit from empty module and build caches. Release tags are immutable: never move, delete, or overwrite one. Fix a bad release forward with a new version bump and a new workflow run.

## Pull requests

Complete the pull request template, include exact verification commands and results, and identify public API, documentation, integration-test, and cross-SDK effects. Do not combine unrelated refactors with behavior changes.
