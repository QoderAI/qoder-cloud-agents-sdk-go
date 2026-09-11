package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/forwardutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	sessionID := flag.String("session-id", "", "continue an existing idle session; its resources are not cleaned up")
	keepSession := flag.Bool("keep-session", false, "retain a newly created session and its resources after success")
	prompt := flag.String("prompt", "", "follow-up message; defaults to asking for the project code")
	config, err := live.Load("forward")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	s := &forwardutil.Example{Client: forward.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "conversation",
		Description: "复用同一个 Session 进行多轮对话，读取历史，并支持下次运行继续会话。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s, *sessionID, *keepSession, *prompt)
		},
	}))
}

func run(ctx context.Context, r *live.Run, s *forwardutil.Example, sessionID string, keepSession bool, prompt string) error {
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
		identity, err := s.Identity(ctx, r)
		if err != nil {
			return err
		}
		template, err := s.Template(ctx, r, forward.TemplateNewParams{EnvironmentID: env})
		if err != nil {
			return err
		}
		created, err := s.NewSession(ctx, r, forward.SessionNewParams{IdentityID: identity, TemplateID: template})
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
	s.Client = forward.NewClient(s.Config.Options()...)
	current, err := s.Client.Sessions.Get(ctx, sessionID)
	if err != nil {
		return err
	}
	if current.Status != "idle" {
		return fmt.Errorf("session must be idle before continuing; status=%s", current.Status)
	}
	r.Step("分页读取已有的用户消息和助手回复")
	page := s.Client.Sessions.Events.ListAutoPaging(ctx, sessionID, forward.SessionEventListParams{Order: forward.String("asc"), Limit: forward.Int(100)})
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
		r.Log("info", "下次运行：go run ./examples/forward/conversation -session-id "+sessionID+"（沿用本次 -env 和 -region 配置）")
	}
	return nil
}
