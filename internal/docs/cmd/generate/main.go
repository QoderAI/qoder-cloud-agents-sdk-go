// Command generate rewrites docs/api/reference.md from the current source
// tree. It is the `make docs` entrypoint.
package main

import (
	"fmt"
	"os"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/docs"
)

func main() {
	if err := docs.Generate("."); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
