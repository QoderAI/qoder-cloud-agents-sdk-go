package main

import (
	"context"
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
		Name:        "session",
		Description: "创建 Identity 和 Template，发送一条消息，通过 SSE 接收助手回复。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

// run demonstrates the Identity -> Template -> Session lifecycle and SSE completion.
func run(ctx context.Context, r *live.Run, s *forwardutil.Example) error {
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
	session, err := s.NewSession(ctx, r, forward.SessionNewParams{IdentityID: identity, TemplateID: template})
	if err != nil {
		return err
	}
	marker := live.Marker()
	return s.Turn(ctx, r, session, "请用一句话介绍你能提供什么帮助，并在回复末尾原样附上："+marker, []string{marker}, false, true)
}
