package docs

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// coreSurfaceMarkers pairs a substring with the minimum number of times it
// must appear in the rendered reference. `## type [Client]` appearing twice
// enforces that BOTH forward.Client and managed.Client survived generation
// (a single-count check would pass when the forward section is dropped).
var coreSurfaceMarkers = []struct {
	needle string
	min    int
}{
	{"# forward", 1},
	{"# managed", 1},
	{"# apierror", 1},
	{"## type [Client]", 2},     // forward.Client AND managed.Client
	{"### func [NewClient]", 2}, // both constructors
	{"## type [Error]", 1},      // apierror.Error
}

// relativeLinkRe matches `[text](target)` markdown links whose target is a
// relative path (no scheme, no leading `#`, no leading `/`).
var relativeLinkRe = regexp.MustCompile(`\[[^\]]*\]\(([^)#][^)]*?)\)`)

// findUnnormalizedLinks returns each QoderAI blob URL that violates the
// normalization contract (SHA ref, or #Lxx anchor on `main`). External links
// and clean `blob/main/...` URLs are ignored.
func findUnnormalizedLinks(md string) []string {
	var findings []string
	rawURLRe := regexp.MustCompile(`https://github\.com/QoderAI/qoder-cloud-agents-sdk-go/blob/([^/]+)/([^)\s>]+)`)
	for _, m := range rawURLRe.FindAllStringSubmatch(md, -1) {
		ref, path := m[1], m[2]
		if ref != "main" {
			findings = append(findings, m[0])
			continue
		}
		if strings.Contains(path, "#L") {
			findings = append(findings, m[0])
		}
	}
	return findings
}

// coreSurfacePresent asserts that every coreSurfaceMarker appears at least
// its `min` times in md. Errors name the missing markers to make CI failures
// self-diagnosing.
func coreSurfacePresent(md string) error {
	var missing []string
	for _, m := range coreSurfaceMarkers {
		if strings.Count(md, m.needle) < m.min {
			missing = append(missing, fmt.Sprintf("%q (need %d)", m.needle, m.min))
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("core public surface missing from generated docs: %v", missing)
	}
	if strings.TrimSpace(md) == "" {
		return fmt.Errorf("generated docs are empty")
	}
	return nil
}

// codeFenceRe matches fenced code blocks (```lang\n...\n```). Content inside
// is stripped before the relative-link scan — Go generics `[T C](args)` and
// similar syntax inside code samples parse as markdown link syntax and
// generate noise.
var codeFenceRe = regexp.MustCompile("(?s)```[^\n]*\n.*?\n```")

// findBrokenInternalLinks returns each relative `[text](target)` link whose
// target does not resolve to an existing file under baseDir. External URLs
// (http/https/mailto) and in-page anchors (`#...`) are ignored. gomarkdoc
// wraps targets in `<...>` — that syntactic wrapper is stripped before
// analysis. Fenced code blocks are excluded from the scan so Go generics
// (`[T Constraint](args)` in code samples) don't false-positive as links.
func findBrokenInternalLinks(baseDir, md string) []string {
	stripped := codeFenceRe.ReplaceAllString(md, "")
	var broken []string
	for _, m := range relativeLinkRe.FindAllStringSubmatch(stripped, -1) {
		target := m[1]
		// gomarkdoc emits `[Text](<target>)`; drop the wrapper if present.
		if strings.HasPrefix(target, "<") && strings.HasSuffix(target, ">") {
			target = target[1 : len(target)-1]
		}
		if strings.HasPrefix(target, "#") {
			continue // in-page anchor
		}
		if strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
			continue
		}
		// Strip fragment / query.
		if i := strings.IndexAny(target, "#?"); i != -1 {
			target = target[:i]
		}
		if target == "" {
			continue
		}
		full := filepath.Join(baseDir, target)
		if _, err := os.Stat(full); err != nil {
			broken = append(broken, target)
		}
	}
	return broken
}

// Check regenerates docs/api/reference.md, verifies zero drift against the
// working tree, and runs the normalization / core-surface / internal-link
// gates in sequence. Any failure is returned as an error — the CLI wrapper
// translates it to a non-zero exit.
func Check(repoRoot string) error {
	if err := Generate(repoRoot); err != nil {
		return fmt.Errorf("regenerate for drift check: %w", err)
	}

	out, err := exec.Command("git", "-C", repoRoot, "diff", "--exit-code", "--", "docs/api").CombinedOutput()
	if err != nil {
		return fmt.Errorf("docs/api drift detected — run `make docs` and commit the result:\n%s", out)
	}

	refPath := filepath.Join(repoRoot, outputRelPath)
	body, err := os.ReadFile(refPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", refPath, err)
	}
	md := string(body)

	if bad := findUnnormalizedLinks(md); len(bad) > 0 {
		return fmt.Errorf("unnormalized source links detected (need `blob/main/<file>` without #Lxx or SHA):\n  %s", joinFirst(bad, 5))
	}
	if err := coreSurfacePresent(md); err != nil {
		return err
	}
	// Relative links resolve from the doc's own directory.
	if broken := findBrokenInternalLinks(filepath.Dir(refPath), md); len(broken) > 0 {
		return fmt.Errorf("broken internal links in %s:\n  %s", outputRelPath, joinFirst(broken, 5))
	}
	return nil
}

// joinFirst formats up to n entries of ss, separated by newlines and 2-space
// indent, appending an ellipsis if truncated.
func joinFirst(ss []string, n int) string {
	var buf bytes.Buffer
	limit := len(ss)
	if limit > n {
		limit = n
	}
	for i := 0; i < limit; i++ {
		if i > 0 {
			buf.WriteString("\n  ")
		}
		buf.WriteString(ss[i])
	}
	if len(ss) > n {
		fmt.Fprintf(&buf, "\n  ... (%d more)", len(ss)-n)
	}
	return buf.String()
}
