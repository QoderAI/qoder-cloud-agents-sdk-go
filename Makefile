PYTHON ?= python3
LIVE_ENV_FILE ?= .env.live
.DEFAULT_GOAL := test

.PHONY: build test test-unit test-contract test-live test-live-check test-live-managed test-live-managed-check test-live-all check-version docs docs-check lint

build:
	go build ./...

# Regenerate the committed API reference from source (public forward / managed
# / convention packages via a pinned gomarkdoc; see internal/docs).
docs:
	go run ./internal/docs/cmd/generate

# Drift + normalization + coarse core-surface + internal-link gate. Red when
# committed docs/api/reference.md lags source, when any source link still
# carries a #Lxx anchor, when the core public surface is missing from the
# output, or when an internal relative link is broken.
docs-check:
	go run ./internal/docs/cmd/check

# gofmt + go vet gate. Fails if any tracked .go file is not gofmt'd or if vet
# reports a diagnostic. Wired into the CI matrix (see .github/workflows/ci.yml).
lint:
	@unformatted=$$(gofmt -l .); \
		if test -n "$$unformatted"; then \
			echo "gofmt reports unformatted files:" >&2; \
			echo "$$unformatted" >&2; \
			exit 1; \
		fi
	go vet ./...

# Run before tagging a release: the reported version is a compile-time constant,
# so tagging without bumping it makes the SDK report a version it is not.
check-version:
	@const=$$(sed -n 's/^const packageVersion = "\(.*\)"$$/\1/p' convention/version.go); \
		tag=$$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//'); \
		if test -z "$$tag"; then echo "no tag reachable; skipped" >&2; exit 0; fi; \
		if test "$$const" != "$$tag"; then \
			echo "packageVersion is $$const but the latest tag is v$$tag" >&2; exit 1; \
		fi; \
		echo "packageVersion matches v$$tag"

test: test-unit test-contract

test-unit:
	go test ./...
	go test -tags live -run 'Cleanup.*Offline$$' ./forward ./managed

test-contract:
	go test -v \
		-run '^(Test.*APIContracts|TestForwardAPIInventory|TestForwardFailureContracts)$$' ./forward
	go test -v ./managed

test-live-check:
	@test -f "$(LIVE_ENV_FILE)" || (echo "missing $(LIVE_ENV_FILE)" >&2; exit 2)
	@set -a; . "$(abspath $(LIVE_ENV_FILE))"; set +a; \
		test -n "$$QODER_FORWARD_PAT" || (echo "QODER_FORWARD_PAT is required in $(LIVE_ENV_FILE)" >&2; exit 2)

test-live-managed-check:
	@test -f "$(LIVE_ENV_FILE)" || (echo "missing $(LIVE_ENV_FILE)" >&2; exit 2)
	@set -a; . "$(abspath $(LIVE_ENV_FILE))"; set +a; \
		test -n "$$QODER_MANAGED_PAT" || (echo "QODER_MANAGED_PAT is required in $(LIVE_ENV_FILE)" >&2; exit 2)

test-live: test-live-check
	@set -a; . "$(abspath $(LIVE_ENV_FILE))"; set +a; \
		go test -tags live -v \
			-run 'Live$$' -skip 'E2ELive$$' ./forward

test-live-managed: test-live-managed-check
	@set -a; . "$(abspath $(LIVE_ENV_FILE))"; set +a; \
		go test -tags live -v -run 'Live$$' -skip 'E2ELive$$' ./managed

test-live-all: test-live test-live-managed

.PHONY: test-e2e-check test-e2e

test-e2e-check: test-live-check test-live-managed-check
	@set -a; . "$(abspath $(LIVE_ENV_FILE))"; set +a; \
		for mode in FORWARD MANAGED; do \
			for key in MODEL LIVE_ALLOW_WRITE LIVE_ALLOW_EXECUTION; do \
				name="QODER_$${mode}_$${key}"; value=$$(printenv "$$name"); \
				if test -z "$$value" || { test "$$key" != MODEL && test "$$value" != true; }; then \
					echo "E2E requires $$name (see .env.live.example)" >&2; exit 2; \
				fi; \
			done; \
		done

# Prerequisites run before the recipe, so no live call can precede a failing unit suite.
test-e2e: test-unit
	@$(MAKE) --no-print-directory test-e2e-check
	@mkdir -p build/test-results
	@set -a; . "$(abspath $(LIVE_ENV_FILE))"; set +a; \
		go test -tags live -count=1 -p=1 -parallel=1 -timeout=45m \
			-json -run 'E2ELive$$' ./forward ./managed > build/test-results/e2e.json; \
		result=$$?; echo "E2E result: build/test-results/e2e.json (exit $$result)"; exit $$result
