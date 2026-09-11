package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestSessionThreadEventAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listSessionThreadEvents", status)
				res, err := client.Sessions.Threads.Events.List(context.Background(), pathSegment, pathSegment, contractParams[forward.SessionThreadEventListParams](t, "listSessionThreadEvents"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("StreamEvents", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "streamSessionThreadEvents", status)
				stream := client.Sessions.Threads.Events.StreamEvents(context.Background(), pathSegment, pathSegment, contractParams[forward.SessionThreadEventStreamParams](t, "streamSessionThreadEvents"))
				defer stream.Close()
				next := stream.Next()
				if status == 200 && !next {
					t.Fatalf("missing event: %v", stream.Err())
				}
				checkError(t, status, stream.Err())
				if status == 200 {
					checkDecoded(t, stream.Current())
				}
			})
		}
	})
}
