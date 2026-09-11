package main

import (
	"context"
	"fmt"
	"os"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
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
		Name:        "streaming-deltas",
		Description: "订阅文本增量，更新消息预览，再用最终完整消息替换预览。",
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
	agent, err := s.Agent(ctx, r, managed.AgentNewParams{})
	if err != nil {
		return err
	}
	session, err := s.NewSession(ctx, r, managed.SessionNewParams{EnvironmentID: env, Agent: managed.SessionNewParamsAgentUnion{OfString: managed.String(agent)}})
	if err != nil {
		return err
	}
	// Preview events are not stored in history: subscribe before sending the message.
	r.Step("先订阅 SSE，开启 agent.message 文本增量")
	stream := s.Client.Sessions.Events.StreamEvents(ctx, session, managed.SessionEventStreamParams{EventDeltas: []managed.ManagedAgentsDeltaType{"agent.message"}}, option.WithRequestTimeout(s.Config.Timeout))
	defer stream.Close()
	if err := stream.Err(); err != nil {
		return err
	}
	marker := live.Marker()
	prompt := "请分三句话解释为什么多轮对话要复用 Session ID，最后原样附上：" + marker
	r.Step("发送消息")
	r.Message("user", prompt)
	after, err := s.Send(ctx, session, prompt)
	if err != nil {
		return err
	}
	stopWaiting := r.Wait("接收增量和最终回复")
	defer stopWaiting()
	var preview live.MessagePreview
	result := live.TurnResult{LastID: after}
	deltas := 0
	for stream.Next() {
		raw := stream.Current().RawJSON()
		update, err := preview.Observe(raw)
		if err != nil {
			return err
		}
		if update != nil {
			if update.Final {
				r.Log("info", "消息 "+update.ID+" 收到完整内容，替换该消息的预览。")
			} else {
				deltas++
				// Each snapshot replaces the previous value for this message ID.
				r.Message("preview", update.ID+"："+update.Text)
			}
		}
		if err := r.ObserveTurn(&result, raw); err != nil {
			return err
		}
		if result.Complete {
			break
		}
	}
	if err := stream.Err(); err != nil {
		return err
	}
	stopWaiting()
	if err := result.Verify([]string{marker}, false); err != nil {
		return err
	}
	r.Log("checked", fmt.Sprintf("收到 %d 个文本增量；最终完整回复已校验。", deltas))
	if deltas == 0 {
		r.Log("info", "增量是尽力提供的预览；本次未收到增量，使用最终完整消息。")
	}
	return nil
}
