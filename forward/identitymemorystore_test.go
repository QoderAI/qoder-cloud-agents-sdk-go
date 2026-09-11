package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestIdentityMemoryStoreAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listMemoryStoreMounts", status)
				res, err := client.Identities.MemoryStores.List(context.Background(), pathSegment, pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Mount", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "mountMemoryStore", status)
				res, err := client.Identities.MemoryStores.Mount(context.Background(), pathSegment, pathSegment, contractParams[forward.IdentityMemoryStoreMountParams](t, "mountMemoryStore"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Detach", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "detachMemoryStore", status)
				res, err := client.Identities.MemoryStores.Detach(context.Background(), pathSegment, pathSegment, pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
}
