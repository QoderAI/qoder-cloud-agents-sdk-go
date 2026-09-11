package forward_test

import (
	"context"
	"testing"

	forward "github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestMemoryStoreMemoryAPIContracts(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "listMemoryStoreMemoryMemories", status)
				res, err := client.MemoryStores.Memories.List(context.Background(), pathSegment, contractParams[forward.MemoryStoreMemoryListParams](t, "listMemoryStoreMemoryMemories"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("New", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "createMemoryStoreMemory", status)
				res, err := client.MemoryStores.Memories.New(context.Background(), pathSegment, contractParams[forward.MemoryStoreMemoryNewParams](t, "createMemoryStoreMemory"))
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
				client := contractClient(t, "getMemoryStoreMemory", status)
				res, err := client.MemoryStores.Memories.Get(context.Background(), pathSegment, pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Update", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "updateMemoryStoreMemory", status)
				res, err := client.MemoryStores.Memories.Update(context.Background(), pathSegment, pathSegment, contractParams[forward.MemoryStoreMemoryUpdateParams](t, "updateMemoryStoreMemory"))
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
	t.Run("Delete", func(t *testing.T) {
		for _, status := range []int{200, 400} {
			t.Run(statusName(status), func(t *testing.T) {
				client := contractClient(t, "deleteMemoryStoreMemory", status)
				res, err := client.MemoryStores.Memories.Delete(context.Background(), pathSegment, pathSegment)
				checkError(t, status, err)
				if status == 200 {
					checkDecoded(t, res)
				}
			})
		}
	})
}
