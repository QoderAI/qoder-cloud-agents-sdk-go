package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestSessionEventAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listSessionEvents", status)
				res, err := client.Sessions.Events.List(context.Background(), pathSegment, contractParams[forward.SessionEventListParams](t, "listSessionEvents"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Send", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "sendSessionEvents", status)
				res, err := client.Sessions.Events.Send(context.Background(), pathSegment, contractParams[forward.SessionEventSendParams](t, "sendSessionEvents"))
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
				client := contractClient(t, "streamSessionEvents", status)
				stream := client.Sessions.Events.StreamEvents(context.Background(), pathSegment, contractParams[forward.SessionEventStreamParams](t, "streamSessionEvents"))
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
