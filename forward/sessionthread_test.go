package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestSessionThreadAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listSessionThreads", status)
				res, err := client.Sessions.Threads.List(context.Background(), pathSegment, contractParams[forward.SessionThreadListParams](t, "listSessionThreads"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Get", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "getSessionThread", status)
				res, err := client.Sessions.Threads.Get(context.Background(), pathSegment, pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Archive", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "archiveSessionThread", status)
				res, err := client.Sessions.Threads.Archive(context.Background(), pathSegment, pathSegment, contractParams[forward.SessionThreadArchiveParams](t, "archiveSessionThread"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
}
