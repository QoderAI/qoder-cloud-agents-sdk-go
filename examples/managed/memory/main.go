package main

import (
	"context"
	"errors"
	"fmt"
	"os"

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
		Name:        "memory",
		Description: "写入项目发布约定并挂载到 Session，验证全新会话的上线安排体现这些记忆。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

func run(ctx context.Context, r *live.Run, s *managedutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	memory := live.NewProjectMemory()
	r.Step("创建记忆库（Memory Store）")
	store, err := s.Client.MemoryStores.New(ctx, managed.MemoryStoreNewParams{Name: live.Name("memory")})
	if err != nil {
		return err
	}
	r.Track("memory_store", store.ID, func(ctx context.Context) error {
		_, err := s.Client.MemoryStores.Delete(ctx, store.ID, managed.MemoryStoreDeleteParams{})
		return err
	})
	r.Step("通过 SDK 写入项目背景和发布约定")
	entry, err := s.Client.MemoryStores.Memories.New(ctx, store.ID, managed.MemoryStoreMemoryNewParams{Path: live.ProjectMemoryPath, Content: managed.String(memory.Content())})
	if err != nil {
		return err
	}
	r.Message("memory", memory.Content())
	r.Step("写入 MEMORY.md 索引，供新会话发现项目记忆")
	index, err := s.Client.MemoryStores.Memories.New(ctx, store.ID, managed.MemoryStoreMemoryNewParams{Path: "MEMORY.md", Content: managed.String(memory.Index())})
	if err != nil {
		return err
	}
	r.Log("info", "索引只包含正文入口，不包含发布时间、联系人和回滚版本")
	r.Step("通过 SDK 重新读取，确认记忆正文和索引已保存")
	saved, err := s.Client.MemoryStores.Memories.Get(ctx, entry.ID, managed.MemoryStoreMemoryGetParams{MemoryStoreID: store.ID})
	if err != nil {
		return err
	}
	if saved.Content != memory.Content() {
		return errors.New("读取的记忆内容与写入内容不一致")
	}
	savedIndex, err := s.Client.MemoryStores.Memories.Get(ctx, index.ID, managed.MemoryStoreMemoryGetParams{MemoryStoreID: store.ID})
	if err != nil {
		return err
	}
	if savedIndex.Content != memory.Index() {
		return errors.New("读取的记忆索引与写入内容不一致")
	}
	r.Log("checked", "记忆正文和索引已保存，内容与写入一致")
	agent, err := s.Agent(ctx, r, managed.AgentNewParams{})
	if err != nil {
		return err
	}
	r.Log("info", "下面创建全新会话；项目约定只保存在记忆库中，不加入系统指令或聊天历史")
	session, err := s.NewSession(ctx, r, managed.SessionNewParams{
		EnvironmentID: env,
		Agent:         managed.SessionNewParamsAgentUnion{OfString: managed.String(agent)},
		Resources: []managed.SessionNewParamsResourceUnion{{OfMemoryStore: &managed.ManagedAgentsMemoryStoreResourceParam{
			Type: "memory_store", MemoryStoreID: store.ID, Access: "read_only",
		}}},
	})
	if err != nil {
		return err
	}
	r.Step("查询新会话的资源，确认记忆库已挂载")
	page := s.Client.Sessions.Resources.ListAutoPaging(ctx, session, managed.SessionResourceListParams{})
	found := false
	for page.Next() {
		resource := page.Current()
		if resource.Type == "memory_store" && resource.MemoryStoreID == store.ID {
			found = true
		}
	}
	if err := page.Err(); err != nil {
		return err
	}
	if !found {
		return errors.New("新会话的资源中未找到指定记忆库")
	}
	r.Log("checked", "会话已挂载记忆库："+store.ID)
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
