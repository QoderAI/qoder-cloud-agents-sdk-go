package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestSessionResourceAPIContracts(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "addSessionResource", status)
				res, err := client.Sessions.Resources.Add(context.Background(), pathSegment, contractParams[forward.SessionResourceAddParams](t, "addSessionResource"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
}
