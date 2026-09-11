package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/managedutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func main() {
	config, err := live.Load("managed")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	s := &managedutil.Example{Client: managed.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "dream",
		Description: "让 Dream 整理记忆，再读取输出记忆库验证整理结果。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

// run verifies an asynchronous consolidation job by reading its output memory.
func run(ctx context.Context, r *live.Run, s *managedutil.Example) error {
	if s.Model == "" {
		if err := s.SelectModel(ctx, r); err != nil {
			return err
		}
	}
	r.Step("创建记忆库（Memory Store）")
	store, err := s.Client.MemoryStores.New(ctx, managed.MemoryStoreNewParams{Name: live.Name("dream-input")})
	if err != nil {
		return err
	}
	r.Track("input_memory_store", store.ID, func(ctx context.Context) error {
		_, err := s.Client.MemoryStores.Delete(ctx, store.ID, managed.MemoryStoreDeleteParams{})
		return err
	})
	marker := live.Marker()
	r.Step("向记忆库写入示例内容和随机校验值")
	_, err = s.Client.MemoryStores.Memories.New(ctx, store.ID, managed.MemoryStoreMemoryNewParams{Path: "sdk-example/source.md", Content: managed.String("Permanent project verification code: " + marker + ". Preserve this exact code during consolidation.")})
	if err != nil {
		return err
	}
	r.Step("创建 Dream，整理输入记忆库中的内容")
	dream, err := s.Client.Dreams.New(ctx, managed.DreamNewParams{Inputs: []managed.DreamInputUnionParam{{OfMemoryStore: &managed.DreamMemoryStoreInputParam{Type: "memory_store", MemoryStoreID: store.ID}}}, Model: managed.DreamNewParamsModelUnion{OfString: managed.String(s.Model)}, Instructions: managed.String("Consolidate the supplied memory into sdk-example/consolidated.md. Preserve the exact project verification code. Keep the original source.")})
	if err != nil {
		return err
	}
	dreamID := dream.ID
	r.Track("dream", dreamID, func(ctx context.Context) error { return finishDream(ctx, s, r, dreamID, store.ID) })
	stopWaiting := r.Wait("等待 Dream 整理记忆")
	defer stopWaiting()
	previous := ""
	for dream.Status == "pending" || dream.Status == "running" {
		if previous != string(dream.Status) {
			r.Log("info", "Dream 状态："+live.Status(string(dream.Status)))
			previous = string(dream.Status)
		}
		if err = live.Pause(ctx); err != nil {
			return err
		}
		dream, err = s.Client.Dreams.Get(ctx, dreamID, managed.DreamGetParams{})
		if err != nil {
			return err
		}
	}
	stopWaiting()
	r.Step("读取 Dream 输出，检查整理后的记忆")
	if dream.Status != "completed" || len(dream.Outputs) == 0 {
		return fmt.Errorf("dream=%s status=%s outputs=%d", dreamID, dream.Status, len(dream.Outputs))
	}
	for _, out := range dream.Outputs {
		page := s.Client.MemoryStores.Memories.ListAutoPaging(ctx, out.MemoryStoreID, managed.MemoryStoreMemoryListParams{})
		for page.Next() {
			memory := page.Current()
			if memory.Path != "sdk-example/consolidated.md" {
				continue
			}
			got, err := s.Client.MemoryStores.Memories.Get(ctx, memory.ID, managed.MemoryStoreMemoryGetParams{MemoryStoreID: out.MemoryStoreID})
			if err != nil {
				return err
			}
			if strings.Contains(got.Content, marker) {
				r.Log("checked", "输出记忆保留了原始校验值："+got.Path)
				r.Message("assistant", got.Content)
				return nil
			}
		}
		if err := page.Err(); err != nil {
			return err
		}
	}
	return errors.New("Dream did not persist consolidated memory with the original verification code")
}
func finishDream(ctx context.Context, s *managedutil.Example, r *live.Run, id, inputID string) error {
	current, err := s.Client.Dreams.Get(ctx, id, managed.DreamGetParams{})
	if err != nil {
		return err
	}
	if current.Status == "pending" || current.Status == "running" {
		if _, err = s.Client.Dreams.Cancel(ctx, id, managed.DreamCancelParams{}); err != nil {
			return err
		}
		for current.Status == "pending" || current.Status == "running" {
			if err = live.Pause(ctx); err != nil {
				return err
			}
			current, err = s.Client.Dreams.Get(ctx, id, managed.DreamGetParams{})
			if err != nil {
				return err
			}
		}
	}
	var failures []error
	if current.SessionID != "" {
		r.Log("info", "清理 Dream 关联的会话："+current.SessionID)
		if err = s.FinishSession(ctx, current.SessionID); err != nil {
			failures = append(failures, err)
		}
	}
	seen := map[string]bool{inputID: true}
	for _, out := range current.Outputs {
		if out.MemoryStoreID == "" || seen[out.MemoryStoreID] {
			continue
		}
		seen[out.MemoryStoreID] = true
		r.Log("info", "清理 Dream 输出记忆库："+out.MemoryStoreID)
		if _, err = s.Client.MemoryStores.Delete(ctx, out.MemoryStoreID, managed.MemoryStoreDeleteParams{}); err != nil {
			failures = append(failures, err)
		}
	}
	if _, err = s.Client.Dreams.Archive(ctx, id, managed.DreamArchiveParams{}); err != nil {
		failures = append(failures, err)
	}
	return errors.Join(failures...)
}
