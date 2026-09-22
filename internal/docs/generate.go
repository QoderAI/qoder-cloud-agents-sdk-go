// Package docs implements the deterministic API reference generator and
// checker for this SDK. It wraps a pinned gomarkdoc release invoked via
// os/exec (not required into go.mod) and post-processes the output to strip
// line-anchor churn from source links so the committed docs/api/reference.md
// is a stable byte-for-byte artifact of the current commit.
package docs

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// gomarkdocVersion is the pinned exact version of the gomarkdoc CLI used by
// this generator. Bump this to intentionally accept doc drift. Never let it
// float.
const gomarkdocVersion = "v1.1.0"

// publicPackages is the set of packages documented by `make docs`. Kept in
// sync with the parent Spec §2 scope: exported forward / managed / convention
// surface only — internal/ (including this tool) and examples/ are excluded.
var publicPackages = []string{
	"./forward",
	"./managed",
	"./convention/...",
}

// repoURL is the canonical repository URL used to compute source links.
// Hard-coded (not read from git remotes) so `make docs` yields identical
// output regardless of clone URL / remote name.
const repoURL = "https://github.com/QoderAI/qoder-cloud-agents-sdk-go"

// outputRelPath is the single-file location of the committed API reference.
// A single file (rather than one per package) keeps the docs/api tree flat
// and simplifies both drift diffing and internal-link checks.
const outputRelPath = "docs/api/reference.md"

// sourceLinkAnchorRe matches `#Lxx` / `#Lxx-Lyy` fragments on QoderAI blob
// URLs. gomarkdoc emits per-line anchors natively; they are churn-sensitive
// (any unrelated line shift falsely reddens the drift gate) and add no value
// for a reader who is landing on `main` at the current commit anyway.
var sourceLinkAnchorRe = regexp.MustCompile(
	`(https://github\.com/QoderAI/qoder-cloud-agents-sdk-go/blob/main/[^)\s>#]+)(#L\d+(?:-L\d+)?)`,
)

// normalizeSourceLinks strips `#Lxx(-Lyy)?` line anchors from QoderAI blob
// URLs pointing at `main`. Non-QoderAI links and QoderAI links on non-`main`
// refs pass through untouched — the latter is intentional so the check stage
// can surface them as unnormalized (i.e. gomarkdoc misconfiguration or a
// SHA-based ref sneaking in).
func normalizeSourceLinks(md string) string {
	return sourceLinkAnchorRe.ReplaceAllString(md, "$1")
}

// stripNonDeterministic is a hook for future non-determinism strippers
// (timestamps, absolute paths, per-user identifiers). The empirical spike on
// gomarkdoc v1.1.0 observed none, so this is currently a pass-through — it
// exists so `Generate` has a single obvious place to add strippers when a
// gomarkdoc bump introduces one.
func stripNonDeterministic(md string) string {
	return md
}

// Generate regenerates docs/api/reference.md from the current source tree.
// repoRoot must be the repo root (module directory). The output directory is
// wiped before write so deleted/renamed exported symbols show up as diff
// noise (no orphan pages) — this is the drift-gate contract.
func Generate(repoRoot string) error {
	apiDir := filepath.Join(repoRoot, "docs", "api")
	if err := os.RemoveAll(apiDir); err != nil {
		return fmt.Errorf("wipe %s: %w", apiDir, err)
	}
	if err := os.MkdirAll(apiDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", apiDir, err)
	}

	raw, err := runGomarkdoc(repoRoot)
	if err != nil {
		return err
	}

	processed := stripNonDeterministic(normalizeSourceLinks(raw))
	out := filepath.Join(repoRoot, outputRelPath)
	if err := os.WriteFile(out, []byte(processed), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", out, err)
	}
	return nil
}

// runGomarkdoc shells out to a pinned gomarkdoc via `go run` and returns the
// raw markdown. `--repository.*` flags force `main`-branch file-level source
// links (see spike observations) — without them gomarkdoc emits no source
// links at all.
func runGomarkdoc(repoRoot string) (string, error) {
	args := []string{
		"run",
		"github.com/princjef/gomarkdoc/cmd/gomarkdoc@" + gomarkdocVersion,
		"--repository.url", repoURL,
		"--repository.default-branch", "main",
		"--repository.path", "/",
	}
	args = append(args, publicPackages...)

	cmd := exec.Command("go", args...)
	cmd.Dir = repoRoot
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("gomarkdoc %s failed: %w\nstderr: %s", gomarkdocVersion, err, stderr.String())
	}
	return stdout.String(), nil
}
