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

## Pull requests

Complete the pull request template, include exact verification commands and results, and identify public API, documentation, integration-test, and cross-SDK effects. Do not combine unrelated refactors with behavior changes.
