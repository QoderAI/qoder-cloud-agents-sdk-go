package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestMemoryStoreMemoryVersionAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listMemoryStoreMemoryVersion", status)
				res, err := client.MemoryStores.MemoryVersions.List(context.Background(), pathSegment, contractParams[forward.MemoryStoreMemoryVersionListParams](t, "listMemoryStoreMemoryVersion"))
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
				client := contractClient(t, "getMemoryStoreMemoryVersion", status)
				res, err := client.MemoryStores.MemoryVersions.Get(context.Background(), pathSegment, pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Redact", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "redactMemoryStoreMemoryVersion", status)
				res, err := client.MemoryStores.MemoryVersions.Redact(context.Background(), pathSegment, pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
}
