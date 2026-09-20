// Command check verifies docs/api/reference.md matches the current source
// tree (drift gate), that every source link is normalized, that the core
// public surface appears in the output, and that internal relative links
// resolve. It is the `make docs-check` entrypoint.
package main

import (
	"fmt"
	"os"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/docs"
)

func main() {
	if err := docs.Check("."); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
