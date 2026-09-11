package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/managedutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func main() {
	sessionID := flag.String("session-id", "", "continue an existing idle session; its resources are not cleaned up")
	keepSession := flag.Bool("keep-session", false, "retain a newly created session and its resources after success")
	prompt := flag.String("prompt", "", "follow-up message; defaults to asking for the project code")
	config, err := live.Load("managed")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	s := &managedutil.Example{Client: managed.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "conversation",
		Description: "复用同一个 Session 进行多轮对话，读取历史，并支持下次运行继续会话。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s, *sessionID, *keepSession, *prompt)
		},
	}))
}

func run(ctx context.Context, r *live.Run, s *managedutil.Example, sessionID string, keepSession bool, prompt string) error {
	verifyRecall := prompt == ""
	if verifyRecall {
		prompt = "只根据本次会话上文，告诉我刚才约定的项目代号。只回复代号，不要使用工具。"
	}
	createdHere := sessionID == ""
	var expected []string
	if createdHere {
		env, err := s.Environment(ctx, r)
		if err != nil {
			return err
		}
		agent, err := s.Agent(ctx, r, managed.AgentNewParams{})
		if err != nil {
			return err
		}
		created, err := s.NewSession(ctx, r, managed.SessionNewParams{EnvironmentID: env, Agent: managed.SessionNewParamsAgentUnion{OfString: managed.String(agent)}})
		if err != nil {
			return err
		}
		sessionID = created
		code := "project-" + live.Marker()
		if verifyRecall {
			expected = []string{code}
		}
		if err := s.Turn(ctx, r, sessionID, "这次项目代号是 "+code+"。请在本次对话中记住它，不要使用工具或写入记忆库。现在只回复：已记住。", []string{"已记住"}, false, true); err != nil {
			return err
		}
	} else {
		r.Log("info", "继续已有 Session；本程序不会清理它及其关联资源。")
	}
	r.Log("info", "Session ID："+sessionID)
	r.Step("重新创建 SDK 客户端，通过 Session ID 读取服务端会话")
	// A new client has no local chat history; the server retains it under sessionID.
	s.Client = managed.NewClient(s.Config.Options()...)
	current, err := s.Client.Sessions.Get(ctx, sessionID, managed.SessionGetParams{})
	if err != nil {
		return err
	}
	if current.Status != "idle" {
		return fmt.Errorf("session must be idle before continuing; status=%s", current.Status)
	}
	r.Step("分页读取已有的用户消息和助手回复")
	page := s.Client.Sessions.Events.ListAutoPaging(ctx, sessionID, managed.SessionEventListParams{Order: "asc", Limit: managed.Int(100)})
	count := 0
	for page.Next() {
		count++
		if count > 2000 {
			return errors.New("history exceeds this example's 2000-event limit")
		}
		role, text, err := live.HistoryMessage(page.Current().RawJSON())
		if err != nil {
			return err
		}
		if text != "" {
			r.Message("history", role+"："+text)
		}
	}
	if err := page.Err(); err != nil {
		return err
	}
	if err := s.Turn(ctx, r, sessionID, prompt, expected, false, true); err != nil {
		return err
	}
	if createdHere && keepSession {
		r.RetainResources()
		r.Log("info", "下次运行：go run ./examples/managed/conversation -session-id "+sessionID+"（沿用本次 -env 和 -region 配置）")
	}
	return nil
}
