package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/forwardutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	config, err := live.Load("forward")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	s := &forwardutil.Example{Client: forward.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "memory",
		Description: "写入项目发布约定并绑定到 Identity + Template，验证全新会话的上线安排体现这些记忆。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

// run checks persisted Identity binding and its use by a fresh Session.
func run(ctx context.Context, r *live.Run, s *forwardutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	r.Step("创建有展示名的专用 Identity，用于记忆绑定")
	// Forward currently uses Identity.Name when provisioning its default store.
	identity, err := s.Client.Identities.New(ctx, forward.IdentityNewParams{ExternalID: live.Name("identity"), Name: forward.String("记忆示例用户")})
	if err != nil {
		return err
	}
	r.Track("identity", identity.ID, func(ctx context.Context) error {
		// This Identity was created by this run. Clear also removes the default
		// memory store automatically provisioned when its Session was created.
		cleared, err := s.Client.Identities.Clear(ctx, identity.ID, forward.IdentityClearParams{Reason: forward.String("SDK memory example cleanup")})
		if err != nil {
			return err
		}
		if cleared.Status != "completed" {
			return errors.New("专用 Identity 的关联资源清理未完成")
		}
		r.Log("checked", "已清理专用 Identity 的关联资源（包含自动生成的默认记忆库）")
		_, err = s.Client.Identities.Delete(ctx, identity.ID)
		return err
	})
	memory := live.NewProjectMemory()
	r.Step("创建记忆库（Memory Store）")
	store, err := s.Client.MemoryStores.New(ctx, forward.MemoryStoreNewParams{Name: live.Name("memory"), IdempotencyKey: live.Name("memory-key")})
	if err != nil {
		return err
	}
	r.Track("memory_store", store.ID, func(ctx context.Context) error { _, err := s.Client.MemoryStores.Delete(ctx, store.ID); return err })
	r.Step("通过 SDK 写入项目背景和发布约定")
	entry, err := s.Client.MemoryStores.Memories.New(ctx, store.ID, forward.MemoryStoreMemoryNewParams{Path: live.ProjectMemoryPath, Content: memory.Content()})
	if err != nil {
		return err
	}
	r.Message("memory", memory.Content())
	r.Step("写入 MEMORY.md 索引，供新会话发现项目记忆")
	index, err := s.Client.MemoryStores.Memories.New(ctx, store.ID, forward.MemoryStoreMemoryNewParams{Path: "MEMORY.md", Content: memory.Index()})
	if err != nil {
		return err
	}
	r.Log("info", "索引只包含正文入口，不包含发布时间、联系人和回滚版本")
	r.Step("通过 SDK 重新读取，确认记忆正文和索引已保存")
	saved, err := s.Client.MemoryStores.Memories.Get(ctx, store.ID, entry.ID)
	if err != nil {
		return err
	}
	if saved.Content != memory.Content() {
		return errors.New("读取的记忆内容与写入内容不一致")
	}
	savedIndex, err := s.Client.MemoryStores.Memories.Get(ctx, store.ID, index.ID)
	if err != nil {
		return err
	}
	if savedIndex.Content != memory.Index() {
		return errors.New("读取的记忆索引与写入内容不一致")
	}
	r.Log("checked", "记忆正文和索引已保存，内容与写入一致")

	template, err := s.Template(ctx, r, forward.TemplateNewParams{EnvironmentID: env})
	if err != nil {
		return err
	}
	r.Step("将记忆库绑定到 Identity 和 Template")
	_, err = s.Client.Identities.MemoryStores.Mount(ctx, identity.ID, template, forward.IdentityMemoryStoreMountParams{MemoryStoreID: store.ID})
	if err != nil {
		return err
	}
	r.Track("memory_mount", store.ID, func(ctx context.Context) error {
		_, err := s.Client.Identities.MemoryStores.Detach(ctx, identity.ID, template, store.ID)
		return err
	})
	r.Step("查询绑定，确认记忆库已关联")
	mounts, err := s.Client.Identities.MemoryStores.List(ctx, identity.ID, template)
	if err != nil {
		return err
	}
	found := false
	for _, mount := range mounts.Data {
		if mount.MemoryStoreID == store.ID {
			found = true
		}
	}
	if !found {
		return errors.New("memory mount was not persisted")
	}
	r.Log("checked", "已查到绑定的记忆库："+store.ID)
	r.Log("info", "下面创建全新会话；项目约定只保存在记忆库中，不加入系统指令或聊天历史")
	session, err := s.NewSession(ctx, r, forward.SessionNewParams{IdentityID: identity.ID, TemplateID: template})
	if err != nil {
		return err
	}
	r.Step("请助手结合记忆拟定上线安排")
	r.Message("user", memory.Prompt())
	after, err := s.Send(ctx, session, memory.Prompt())
	if err != nil {
		return err
	}
	reply, err := s.WaitReply(ctx, r, session, after, false)
	if err != nil {
		return err
	}
	return memory.Verify(r, reply)
}
