package docs

import "testing"

func TestFindUnnormalizedLinks_FlagsQoderAIAnchors(t *testing.T) {
	md := `Fine: https://github.com/QoderAI/qoder-cloud-agents-sdk-go/blob/main/forward/client.go
Bad line anchor: https://github.com/QoderAI/qoder-cloud-agents-sdk-go/blob/main/forward/client.go#L14
Bad range: https://github.com/QoderAI/qoder-cloud-agents-sdk-go/blob/main/managed/client.go#L10-L20
Bad SHA ref: https://github.com/QoderAI/qoder-cloud-agents-sdk-go/blob/abcd1234/forward/client.go
External untouched: https://github.com/other/repo/blob/main/x.go#L1
`
	got := findUnnormalizedLinks(md)
	if len(got) != 3 {
		t.Fatalf("want 3 unnormalized findings, got %d: %v", len(got), got)
	}
}

func TestFindUnnormalizedLinks_ZeroOnCleanInput(t *testing.T) {
	md := `Clean file-level link https://github.com/QoderAI/qoder-cloud-agents-sdk-go/blob/main/forward/client.go here.`
	if got := findUnnormalizedLinks(md); len(got) != 0 {
		t.Fatalf("want zero findings, got %v", got)
	}
}

func TestCoreSurfacePresent_OKWhenAllPresent(t *testing.T) {
	md := "# forward\n\n## type [Client]\n### func [NewClient]\n\n# managed\n\n## type [Client]\n### func [NewClient]\n\n# apierror\n\n## type [Error]"
	if err := coreSurfacePresent(md); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCoreSurfacePresent_FailsWhenClientMissing(t *testing.T) {
	// forward/client.go section deleted (simulates generator misconfiguration
	// that drops the core public surface).
	md := "# forward\n\n## type [Batch]\n\n# managed\n\n## type [Client]\n### func [NewClient]\n\n# apierror\n\n## type [Error]"
	if err := coreSurfacePresent(md); err == nil {
		t.Fatal("expected error when forward Client is missing, got nil")
	}
}

func TestFindBrokenInternalLinks_FlagsMissingFile(t *testing.T) {
	// A relative link to a sibling doc that does not exist must be flagged.
	base := t.TempDir()
	broken := findBrokenInternalLinks(base, "See [nope](./missing.md).")
	if len(broken) != 1 {
		t.Fatalf("want 1 broken link, got %v", broken)
	}
}
